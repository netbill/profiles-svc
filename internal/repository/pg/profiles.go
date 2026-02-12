package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/repository"
)

const profilesTable = "profiles"
const ProfilesColumns = "account_id, username, official, pseudonym, description, avatar, created_at, updated_at"

func scanProfile(row sq.RowScanner) (p repository.ProfileRow, err error) {
	pseudonym := pgtype.Text{}
	description := pgtype.Text{}
	avatarURL := pgtype.Text{}

	err = row.Scan(
		&p.AccountID,
		&p.Username,
		&p.Official,
		&pseudonym,
		&description,
		&avatarURL,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.ProfileRow{}, nil
	case err != nil:
		return repository.ProfileRow{}, fmt.Errorf("scanning profile: %w", err)
	}

	if pseudonym.Valid {
		p.Pseudonym = &pseudonym.String
	}
	if description.Valid {
		p.Description = &description.String
	}
	if avatarURL.Valid {
		p.Avatar = &avatarURL.String
	}

	return p, nil
}

type profiles struct {
	db       *pgdbx.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewProfilesQ(db *pgdbx.DB) repository.ProfilesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &profiles{
		db:       db,
		selector: builder.Select("profiles.*").From(profilesTable),
		inserter: builder.Insert(profilesTable),
		updater:  builder.Update(profilesTable),
		deleter:  builder.Delete(profilesTable),
		counter:  builder.Select("COUNT(*) AS count").From(profilesTable),
	}
}

func (q *profiles) New() repository.ProfilesQ {
	return NewProfilesQ(q.db)
}

func (q *profiles) Insert(ctx context.Context, input repository.ProfileRow) (repository.ProfileRow, error) {
	query, args, err := q.inserter.SetMap(map[string]interface{}{
		"account_id":  input.AccountID,
		"username":    input.Username,
		"official":    input.Official,
		"pseudonym":   input.Pseudonym,
		"description": input.Description,
	}).Suffix("RETURNING " + ProfilesColumns).ToSql()
	if err != nil {
		return repository.ProfileRow{}, fmt.Errorf("building insert query for %s: %w", profilesTable, err)
	}

	return scanProfile(q.db.QueryRow(ctx, query, args...))
}

func (q *profiles) UpdateMany(ctx context.Context) (int64, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", profilesTable, err)
	}

	tag, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (q *profiles) UpdateOne(ctx context.Context) (repository.ProfileRow, error) {
	q.updater = q.updater.Set("updated_at", time.Now().UTC())

	query, args, err := q.updater.Suffix("RETURNING " + ProfilesColumns).ToSql()
	if err != nil {
		return repository.ProfileRow{}, fmt.Errorf("building update query for %s: %w", profilesTable, err)
	}

	return scanProfile(q.db.QueryRow(ctx, query, args...))
}

func (q *profiles) UpdateUsername(username string) repository.ProfilesQ {
	q.updater = q.updater.Set("username", username)
	return q
}

func (q *profiles) UpdateOfficial(official bool) repository.ProfilesQ {
	q.updater = q.updater.Set("official", official)
	return q
}

func (q *profiles) UpdatePseudonym(v *string) repository.ProfilesQ {
	q.updater = q.updater.Set("pseudonym", v)
	return q
}

func (q *profiles) UpdateDescription(v *string) repository.ProfilesQ {
	q.updater = q.updater.Set("description", v)
	return q
}

func (q *profiles) UpdateAvatar(v *string) repository.ProfilesQ {
	q.updater = q.updater.Set("avatar", v)
	return q
}

func (q *profiles) Get(ctx context.Context) (repository.ProfileRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.ProfileRow{}, fmt.Errorf("building get query for %s: %w", profilesTable, err)
	}

	return scanProfile(q.db.QueryRow(ctx, query, args...))
}

func (q *profiles) Select(ctx context.Context) ([]repository.ProfileRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", profilesTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]repository.ProfileRow, 0)
	for rows.Next() {
		p, err := scanProfile(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning profile: %w", err)
		}
		out = append(out, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *profiles) FilterAccountID(accountID ...uuid.UUID) repository.ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"account_id": accountID})
	q.counter = q.counter.Where(sq.Eq{"account_id": accountID})
	q.deleter = q.deleter.Where(sq.Eq{"account_id": accountID})
	q.updater = q.updater.Where(sq.Eq{"account_id": accountID})
	return q
}

func (q *profiles) FilterUsername(username string) repository.ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"username": username})
	q.counter = q.counter.Where(sq.Eq{"username": username})
	q.deleter = q.deleter.Where(sq.Eq{"username": username})
	q.updater = q.updater.Where(sq.Eq{"username": username})
	return q
}

func (q *profiles) FilterOfficial(official bool) repository.ProfilesQ {
	q.selector = q.selector.Where(sq.Eq{"official": official})
	q.counter = q.counter.Where(sq.Eq{"official": official})
	q.deleter = q.deleter.Where(sq.Eq{"official": official})
	q.updater = q.updater.Where(sq.Eq{"official": official})
	return q
}

func (q *profiles) FilterBestMatch(term string) repository.ProfilesQ {
	like := "%" + term + "%"
	prefix := term + "%"

	q.selector = q.selector.Where(sq.Or{
		sq.ILike{"username": like},
		sq.ILike{"pseudonym": like},
	})
	q.counter = q.counter.Where(sq.Or{
		sq.ILike{"username": like},
		sq.ILike{"pseudonym": like},
	})

	q.selector = q.selector.OrderByClause(sq.Expr(
		`CASE
			WHEN lower(username) = lower(?) THEN 0
			WHEN lower(pseudonym) = lower(?) THEN 1
			WHEN lower(username) LIKE lower(?) THEN 2
			WHEN lower(pseudonym) LIKE lower(?) THEN 3
			WHEN lower(username) LIKE lower(?) THEN 4
			WHEN lower(pseudonym) LIKE lower(?) THEN 5
			ELSE 6
		END`,
		term, term,
		prefix, prefix,
		like, like,
	))

	q.selector = q.selector.OrderBy("username ASC")
	return q
}

func (q *profiles) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", profilesTable, err)
	}

	var count uint

	err = q.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (q *profiles) Page(limit, offset uint) repository.ProfilesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q *profiles) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", profilesTable, err)
	}

	_, err = q.db.Exec(ctx, query, args...)
	return err
}
