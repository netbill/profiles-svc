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

const accountsTable = "accounts"

const accountsColumns = "id, username, role, version, source_created_at, replica_created_at, source_updated_at, replica_updated_at"

func scanAccount(row sq.RowScanner) (r repository.AccountRow, err error) {
	err = row.Scan(
		&r.ID,
		&r.Username,
		&r.Role,
		&r.Version,
		&r.SourceCreatedAt,
		&r.ReplicaCreatedAt,
		&r.SourceUpdatedAt,
		&r.ReplicaUpdatedAt,
	)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return repository.AccountRow{}, nil
	case err != nil:
		return repository.AccountRow{}, fmt.Errorf("scanning account: %w", err)
	}
	return r, nil
}

type accounts struct {
	db       *pgdbx.DB
	inserter sq.InsertBuilder
	selector sq.SelectBuilder
	counter  sq.SelectBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
}

func NewAccountsQ(db *pgdbx.DB) repository.AccountsQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &accounts{
		db:       db,
		inserter: builder.Insert(accountsTable),
		selector: builder.Select(accountsColumns).From(accountsTable),
		counter:  builder.Select("COUNT(*) AS count").From(accountsTable),
		updater:  builder.Update(accountsTable),
		deleter:  builder.Delete(accountsTable),
	}
}

func (q *accounts) New() repository.AccountsQ {
	return NewAccountsQ(q.db)
}

func (q *accounts) Insert(ctx context.Context, input repository.AccountRow) (repository.AccountRow, error) {
	id := pgtype.UUID{Bytes: [16]byte(input.ID), Valid: true}
	username := pgtype.Text{String: input.Username, Valid: true}
	role := pgtype.Text{String: input.Role, Valid: true}
	version := pgtype.Int4{Int32: input.Version, Valid: true}

	sca := pgtype.Timestamptz{Time: input.SourceCreatedAt, Valid: !input.SourceCreatedAt.IsZero()}
	sua := pgtype.Timestamptz{Time: input.SourceUpdatedAt, Valid: !input.SourceUpdatedAt.IsZero()}

	query, args, err := q.inserter.SetMap(map[string]interface{}{
		"id":                id,
		"username":          username,
		"role":              role,
		"version":           version,
		"source_created_at": sca,
		"source_updated_at": sua,
	}).Suffix("RETURNING " + accountsColumns).ToSql()
	if err != nil {
		return repository.AccountRow{}, fmt.Errorf("building insert query for %s: %w", accountsTable, err)
	}

	return scanAccount(q.db.QueryRow(ctx, query, args...))
}

func (q *accounts) Get(ctx context.Context) (repository.AccountRow, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return repository.AccountRow{}, fmt.Errorf("building get query for %s: %w", accountsTable, err)
	}

	row := q.db.QueryRow(ctx, query, args...)
	acc, err := scanAccount(row)

	if err != nil {
		return repository.AccountRow{}, err
	}

	return acc, nil
}

func (q *accounts) Select(ctx context.Context) ([]repository.AccountRow, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", accountsTable, err)
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]repository.AccountRow, 0)
	for rows.Next() {
		r, err := scanAccount(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q *accounts) UpdateMany(ctx context.Context) (int64, error) {
	now := time.Now().UTC()

	q.updater = q.updater.Set("replica_updated_at", pgtype.Timestamptz{Time: now, Valid: true})

	query, args, err := q.updater.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building update query for %s: %w", accountsTable, err)
	}

	tag, err := q.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("executing update query for %s: %w", accountsTable, err)
	}

	return tag.RowsAffected(), nil
}

func (q *accounts) UpdateOne(ctx context.Context) (repository.AccountRow, error) {
	now := time.Now().UTC()

	q.updater = q.updater.Set("replica_updated_at", pgtype.Timestamptz{Time: now, Valid: true})

	query, args, err := q.updater.Suffix("RETURNING " + accountsColumns).ToSql()
	if err != nil {
		return repository.AccountRow{}, fmt.Errorf("building update query for %s: %w", accountsTable, err)
	}

	return scanAccount(q.db.QueryRow(ctx, query, args...))
}

func (q *accounts) UpdateRole(role string) repository.AccountsQ {
	q.updater = q.updater.Set("role", pgtype.Text{String: role, Valid: true})
	return q
}

func (q *accounts) UpdateUsername(username string) repository.AccountsQ {
	q.updater = q.updater.Set("username", pgtype.Text{String: username, Valid: true})
	return q
}

func (q *accounts) UpdateVersion(version int32) repository.AccountsQ {
	q.updater = q.updater.Set("version", pgtype.Int4{Int32: version, Valid: true})
	return q
}

func (q *accounts) UpdateSourceUpdatedAt(source time.Time) repository.AccountsQ {
	q.updater = q.updater.Set("source_updated_at", pgtype.Timestamp{Time: source, Valid: true})
	return q
}

func (q *accounts) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", accountsTable, err)
	}
	_, err = q.db.Exec(ctx, query, args...)
	return err
}

func (q *accounts) Exists(ctx context.Context) (bool, error) {
	subSQL, subArgs, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return false, err
	}

	sql := "SELECT EXISTS (" + subSQL + ")"

	var exists bool
	err = q.db.QueryRow(ctx, sql, subArgs...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("sql=%s args=%v: %w", sql, subArgs, err)
	}

	return exists, nil
}

func (q *accounts) FilterID(id uuid.UUID) repository.AccountsQ {
	pid := pgtype.UUID{Bytes: [16]byte(id), Valid: true}

	q.selector = q.selector.Where(sq.Eq{"id": pid})
	q.counter = q.counter.Where(sq.Eq{"id": pid})
	q.updater = q.updater.Where(sq.Eq{"id": pid})
	q.deleter = q.deleter.Where(sq.Eq{"id": pid})
	return q
}

func (q *accounts) FilterUsername(username string) repository.AccountsQ {
	val := pgtype.Text{String: username, Valid: true}

	q.selector = q.selector.Where(sq.Eq{"username": val})
	q.counter = q.counter.Where(sq.Eq{"username": val})
	q.updater = q.updater.Where(sq.Eq{"username": val})
	q.deleter = q.deleter.Where(sq.Eq{"username": val})
	return q
}

func (q *accounts) FilterVersion(version int32) repository.AccountsQ {
	val := pgtype.Int4{Int32: version, Valid: true}

	q.selector = q.selector.Where(sq.Eq{"version": val})
	q.counter = q.counter.Where(sq.Eq{"version": val})
	q.updater = q.updater.Where(sq.Eq{"version": val})
	q.deleter = q.deleter.Where(sq.Eq{"version": val})
	return q
}

func (q *accounts) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", accountsTable, err)
	}

	var count int64
	err = q.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count < 0 {
		return 0, fmt.Errorf("invalid count for %s: %d", accountsTable, count)
	}

	return uint(count), nil
}

func (q *accounts) Page(limit, offset uint) repository.AccountsQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}
