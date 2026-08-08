package pg

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
	"github.com/netbill/profiles-svc/internal/modules/profile"
	"github.com/netbill/restkit/pagi"
)

const (
	profilesTable = "profiles"
	profilesCols  = "account_id, username, pseudonym, description, avatar_key, version, created_at, updated_at, deleted_at"
)

type ProfileRepo struct {
	db *pgdbx.DB
}

func NewProfileRepo(db *pgdbx.DB) *ProfileRepo {
	return &ProfileRepo{db: db}
}

func scanProfile(row pgx.Row) (p models.Profile, err error) {
	err = row.Scan(
		&p.AccountID,
		&p.Username,
		&p.Pseudonym,
		&p.Description,
		&p.AvatarKey,
		&p.Version,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotExists.Raise(err)
	case err != nil:
		return models.Profile{}, fmt.Errorf("scan profile: %w", err)
	}
	return p, nil
}

// nullIfEmpty maps an update param's cleared-value convention ("" clears the
// field) onto SQL NULL; pgx sends a nil *string as NULL and a non-nil one as
// its value.
func nullIfEmpty(v *string) *string {
	if v == nil || *v == "" {
		return nil
	}
	return v
}

func (r *ProfileRepo) Create(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error) {
	const query = `
		INSERT INTO ` + profilesTable + ` (account_id, username)
		VALUES ($1, $2)
		RETURNING ` + profilesCols

	res, err := scanProfile(r.db.QueryRow(ctx, query, accountID, username))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName != "profiles_pkey" {
			// account_id conflict means a genuine duplicate/replayed
			// AccountCreated (caller's problem, not ours); a username
			// conflict is the caller's cue to regenerate and retry.
			return models.Profile{}, errx.ErrorProfileUsernameTaken.Raise(err)
		}
		return models.Profile{}, fmt.Errorf(
			"failed to insert profile for account id %s, cause: %w", accountID, err,
		)
	}

	return res, nil
}

func (r *ProfileRepo) GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error) {
	const query = `SELECT ` + profilesCols + ` FROM ` + profilesTable + ` WHERE account_id = $1 AND deleted_at IS NULL`

	res, err := scanProfile(r.db.QueryRow(ctx, query, accountID))
	if err != nil {
		return models.Profile{}, fmt.Errorf(
			"failed to get profile by account id %s, cause: %w", accountID, err,
		)
	}

	return res, nil
}

func (r *ProfileRepo) GetByUsername(ctx context.Context, username string) (models.Profile, error) {
	const query = `SELECT ` + profilesCols + ` FROM ` + profilesTable + ` WHERE username = $1 AND deleted_at IS NULL`

	res, err := scanProfile(r.db.QueryRow(ctx, query, username))
	if err != nil {
		return models.Profile{}, fmt.Errorf(
			"failed to get profile by username %s, cause: %w", username, err,
		)
	}

	return res, nil
}

func (r *ProfileRepo) ExistByUsername(ctx context.Context, username string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM ` + profilesTable + ` WHERE username = $1 AND deleted_at IS NULL)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, username).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check profile existence by username %s, cause: %w", username, err)
	}

	return exists, nil
}

func (r *ProfileRepo) Update(
	ctx context.Context,
	accountID uuid.UUID,
	input profile.UpdateParams,
) (models.Profile, error) {
	sets := []string{"updated_at = now()", "version = version + 1"}
	args := make([]interface{}, 0, 4)

	if input.Pseudonym != nil {
		args = append(args, nullIfEmpty(input.Pseudonym))
		sets = append(sets, fmt.Sprintf("pseudonym = $%d", len(args)))
	}
	if input.Description != nil {
		args = append(args, nullIfEmpty(input.Description))
		sets = append(sets, fmt.Sprintf("description = $%d", len(args)))
	}
	if input.AvatarKey != nil {
		args = append(args, nullIfEmpty(input.AvatarKey))
		sets = append(sets, fmt.Sprintf("avatar_key = $%d", len(args)))
	}

	args = append(args, accountID)
	query := `UPDATE ` + profilesTable + ` SET ` + strings.Join(sets, ", ") +
		fmt.Sprintf(" WHERE account_id = $%d AND deleted_at IS NULL RETURNING ", len(args)) + profilesCols

	res, err := scanProfile(r.db.QueryRow(ctx, query, args...))
	if err != nil {
		return models.Profile{}, fmt.Errorf(
			"failed to update profile by account id %s, cause: %w", accountID, err,
		)
	}

	return res, nil
}

func (r *ProfileRepo) UpdateUsername(
	ctx context.Context,
	accountID uuid.UUID,
	username string,
) (models.Profile, error) {
	const query = `
		UPDATE ` + profilesTable + `
		SET username = $1, updated_at = now(), version = version + 1
		WHERE account_id = $2 AND deleted_at IS NULL
		RETURNING ` + profilesCols

	res, err := scanProfile(r.db.QueryRow(ctx, query, username, accountID))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.Profile{}, errx.ErrorProfileUsernameTaken.Raise(err)
		}
		return models.Profile{}, fmt.Errorf(
			"failed to update profile username by account id %s, cause: %w", accountID, err,
		)
	}

	return res, nil
}

// Filter powers the public profile search (GET /profiles). Ranking mirrors
// a simple relevance order: exact match first, then prefix match, then
// substring match, username taking priority over pseudonym within each tier.
func (r *ProfileRepo) Filter(
	ctx context.Context,
	params profile.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Profile], error) {
	if limit == 0 {
		limit = 10
	}

	where := " WHERE deleted_at IS NULL"
	var args []interface{}

	if params.Text != nil && *params.Text != "" {
		term := *params.Text
		like := "%" + term + "%"
		prefix := term + "%"

		where += ` AND (username ILIKE $1 OR pseudonym ILIKE $1)`
		args = []interface{}{like, term, prefix}
	}

	const countQuery = `SELECT COUNT(*) FROM ` + profilesTable
	var total uint
	if err := r.db.QueryRow(ctx, countQuery+where, args...).Scan(&total); err != nil {
		return pagi.Page[[]models.Profile]{}, fmt.Errorf("failed to count profiles: %w", err)
	}

	order := " ORDER BY username ASC"
	if len(args) > 0 {
		order = ` ORDER BY
			CASE
				WHEN lower(username) = lower($2) THEN 0
				WHEN lower(pseudonym) = lower($2) THEN 1
				WHEN lower(username) LIKE lower($3) THEN 2
				WHEN lower(pseudonym) LIKE lower($3) THEN 3
				WHEN lower(username) LIKE lower($1) THEN 4
				WHEN lower(pseudonym) LIKE lower($1) THEN 5
				ELSE 6
			END,
			username ASC`
	}

	listQuery := `SELECT ` + profilesCols + ` FROM ` + profilesTable + where + order +
		fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	rows, err := r.db.Query(ctx, listQuery, append(append([]interface{}{}, args...), limit, offset)...)
	if err != nil {
		return pagi.Page[[]models.Profile]{}, fmt.Errorf("failed to filter profiles: %w", err)
	}
	defer rows.Close()

	collection := make([]models.Profile, 0, limit)
	for rows.Next() {
		p, err := scanProfile(rows)
		if err != nil {
			return pagi.Page[[]models.Profile]{}, fmt.Errorf("failed to filter profiles: %w", err)
		}
		collection = append(collection, p)
	}
	if err = rows.Err(); err != nil {
		return pagi.Page[[]models.Profile]{}, fmt.Errorf("failed to filter profiles: %w", err)
	}

	return pagi.Page[[]models.Profile]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: total,
	}, nil
}

// Delete soft-deletes the profile: the row is kept (so a later replayed
// AccountCreated for the same account is recognized as a conflict, see
// AccountRepo.Create), and username is anonymized to free it up for reuse
// despite the UNIQUE constraint — same trick as AccountRepo.Delete. A
// second call against an already-deleted row matches zero rows and
// surfaces as errx.ErrorProfileNotExists, which the caller treats as
// "already deleted, nothing to do" for idempotency against event replay.
func (r *ProfileRepo) Delete(ctx context.Context, accountID uuid.UUID) (models.Profile, error) {
	// username is VARCHAR(32): "deleted_user" (12) + 20 hex chars fits exactly.
	const query = `
		UPDATE ` + profilesTable + `
		SET
			username   = 'deleted_user' || substr(replace(gen_random_uuid()::text, '-', ''), 1, 20),
			deleted_at = now(),
			updated_at = now(),
			version    = version + 1
		WHERE account_id = $1 AND deleted_at IS NULL
		RETURNING ` + profilesCols

	res, err := scanProfile(r.db.QueryRow(ctx, query, accountID))
	if err != nil {
		return models.Profile{}, fmt.Errorf(
			"failed to delete profile for account id %s, cause: %w", accountID, err,
		)
	}

	return res, nil
}
