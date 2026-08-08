package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
	"github.com/netbill/profiles-svc/internal/modules/account"
)

const (
	accountsTable = "accounts"
	accountsCols  = "id, username, role, version, source_created_at, source_updated_at"
)

// AccountRepo is the local read-mirror of accounts-svc's accounts, kept in
// sync over Kafka (see internal/modules/account). source_created_at/
// source_updated_at come from the upstream event; replica_updated_at is
// this mirror's own bookkeeping column and never leaves this package.
// Rows are never deleted here — deletion state lives on profiles.deleted_at
// (see ProfileRepo), so a deleted account's mirror row just goes stale.
type AccountRepo struct {
	db *pgdbx.DB
}

func NewAccountRepo(db *pgdbx.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func scanAccount(row pgx.Row) (a models.Account, err error) {
	err = row.Scan(&a.ID, &a.Username, &a.Role, &a.Version, &a.CreatedAt, &a.UpdatedAt)
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
		return models.Account{}, fmt.Errorf("insert account, cause: %w", err)
	}

	return acc, nil
}

func (r *AccountRepo) GetByID(ctx context.Context, accountID uuid.UUID) (models.Account, error) {
	const query = `SELECT ` + accountsCols + ` FROM ` + accountsTable + ` WHERE id = $1`

	return scanAccount(r.db.QueryRow(ctx, query, accountID))
}

func (r *AccountRepo) ExistsByID(ctx context.Context, accountID uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM ` + accountsTable + ` WHERE id = $1)`

	var exists bool
	if err := r.db.QueryRow(ctx, query, accountID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check account existence by id %s, cause: %w", accountID, err)
	}

	return exists, nil
}

func (r *AccountRepo) UpdateUsername(
	ctx context.Context,
	accountID uuid.UUID,
	params account.UpdateUsernameParams,
) (models.Account, error) {
	const query = `
		UPDATE ` + accountsTable + `
		SET username = $1, version = $2, source_updated_at = $3, replica_updated_at = now()
		WHERE id = $4
		RETURNING ` + accountsCols

	acc, err := scanAccount(r.db.QueryRow(ctx, query, params.Username, params.Version, params.UpdatedAt, accountID))
	if err != nil {
		return models.Account{}, fmt.Errorf(
			"failed to update account username for account %s, cause: %w", accountID, err,
		)
	}

	return acc, nil
}
