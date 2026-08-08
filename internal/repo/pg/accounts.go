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
)

const (
	accountsTable = "accounts"
	accountsCols  = "id, role, version, source_created_at, source_updated_at, deleted_at"
)

type AccountRepo struct {
	db *pgdbx.DB
}

func NewAccountRepo(db *pgdbx.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func scanAccount(row pgx.Row) (a models.Account, err error) {
	err = row.Scan(
		&a.ID,
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

func (r *AccountRepo) Create(ctx context.Context, account models.Account) (models.Account, error) {
	const query = `
		INSERT INTO ` + accountsTable + ` (id, role, source_created_at, source_updated_at)
		VALUES ($1, $2, $3, $3)
		RETURNING ` + accountsCols

	acc, err := scanAccount(r.db.QueryRow(ctx, query, account.ID, account.Role, account.CreatedAt))
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

// Delete soft-deletes the account by setting deleted_at. It never removes
// the row: deletion state has to survive so a later replayed AccountCreated
// for the same id can be recognized as a conflict (see Create). Idempotent —
// a second call matches zero rows.
func (r *AccountRepo) Delete(ctx context.Context, accountID uuid.UUID) (models.Account, error) {
	const query = `
		UPDATE ` + accountsTable + `
		SET deleted_at = now(), replica_updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING ` + accountsCols

	acc, err := scanAccount(r.db.QueryRow(ctx, query, accountID))
	if err != nil {
		return models.Account{}, fmt.Errorf("delete account %s, cause: %w", accountID, err)
	}

	return acc, nil
}
