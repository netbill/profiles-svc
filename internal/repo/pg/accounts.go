package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
	"github.com/netbill/profiles-svc/internal/modules/account"
)

const (
	accountsTable = "accounts"
	accountsCols  = "id, username, role, version, source_created_at, source_updated_at, deleted_at"
)

// AccountRepo is the local read-mirror of accounts-svc's accounts, kept in
// sync over Kafka (see internal/modules/account). source_created_at/
// source_updated_at come from the upstream event; replica_updated_at is
// this mirror's own bookkeeping column and never leaves this package.
//
// Every mutating query filters "AND deleted_at IS NULL" in its own WHERE
// clause rather than relying on a separate pre-check — a pre-check-then-act
// pair is a TOCTOU race against a concurrent Delete; folding the condition
// into the same statement makes the check-and-mutate atomic.
type AccountRepo struct {
	db *pgdbx.DB
}

func NewAccountRepo(db *pgdbx.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func scanAccount(row pgx.Row) (a models.Account, err error) {
	err = row.Scan(
		&a.ID,
		&a.Username,
		&a.Role,
		&a.Version,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Account{}, errx.ErrorAccountNotExists.Raise(err)
	case err != nil:
		return models.Account{}, fmt.Errorf("scan account: %w", err)
	}
	return a, nil
}

func (r *AccountRepo) Create(ctx context.Context, params account.CreateAccountParams) (models.Account, error) {
	const query = `
		INSERT INTO ` + accountsTable + ` (id, username, role, source_created_at, source_updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING ` + accountsCols

	acc, err := scanAccount(r.db.QueryRow(ctx, query, params.ID, params.Username, params.Role, params.CreatedAt))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Same conflict either way, whether the row is still active or
			// was soft-deleted — id is never reused, so any existing row
			// means this is a duplicate/replayed AccountCreated.
			return models.Account{}, errx.ErrorAccountAlreadyExists.Raise(err)
		}
		return models.Account{}, fmt.Errorf("insert account, cause: %w", err)
	}

	return acc, nil
}

// UpdateUsername applies params.Version unconditionally as the new stored
// version (it comes from the upstream event, not a client-read version),
// but only if it's actually newer than what's stored and the row is still
// active — both folded into the WHERE so a stale/replayed event or a
// concurrent Delete atomically yields zero rows / ErrorAccountNotExists
// instead of racing.
func (r *AccountRepo) UpdateUsername(
	ctx context.Context,
	accountID uuid.UUID,
	params account.UpdateUsernameParams,
) (models.Account, error) {
	const query = `
		UPDATE ` + accountsTable + `
		SET username = $1, version = $2, source_updated_at = $3, replica_updated_at = now()
		WHERE id = $4 AND version < $2 AND deleted_at IS NULL
		RETURNING ` + accountsCols

	acc, err := scanAccount(r.db.QueryRow(ctx, query, params.Username, params.Version, params.UpdatedAt, accountID))
	if err != nil {
		return models.Account{}, fmt.Errorf(
			"failed to update account username for account %s, cause: %w", accountID, err,
		)
	}

	return acc, nil
}

// Delete anonymizes username (freeing the UNIQUE slot for reuse) and sets
// deleted_at. It never removes the row: deletion state has to survive so a
// later replayed AccountCreated for the same id can be recognized as a
// conflict (see Create). Idempotent — a second call matches zero rows.
func (r *AccountRepo) Delete(ctx context.Context, accountID uuid.UUID) (models.Account, error) {
	// username is VARCHAR(32): "deleted_user" (12) + 20 hex chars fits exactly.
	const query = `
		UPDATE ` + accountsTable + `
		SET
			username   = 'deleted_user' || substr(replace(gen_random_uuid()::text, '-', ''), 1, 20),
			deleted_at = now(),
			replica_updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING ` + accountsCols

	acc, err := scanAccount(r.db.QueryRow(ctx, query, accountID))
	if err != nil {
		return models.Account{}, fmt.Errorf("delete account %s, cause: %w", accountID, err)
	}

	return acc, nil
}
