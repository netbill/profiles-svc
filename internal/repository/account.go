package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/core/modules/account"
)

type AccountRow struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Role     string    `db:"role"`
	Version  int32     `db:"version"`

	SourceCreatedAt  time.Time `db:"source_created_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
	SourceUpdatedAt  time.Time `db:"source_updated_at"`
	ReplicaUpdatedAt time.Time `db:"replica_updated_at"`
}

func (a AccountRow) IsNil() bool {
	return a.ID == uuid.Nil
}

func (a AccountRow) ToModel() models.Account {
	return models.Account{
		ID:        a.ID,
		Username:  a.Username,
		Role:      a.Role,
		Version:   a.Version,
		CreatedAt: a.SourceCreatedAt,
		UpdatedAt: a.SourceUpdatedAt,
	}
}

type AccountsQ interface {
	New() AccountsQ

	Insert(ctx context.Context, input AccountRow) (AccountRow, error)

	Get(ctx context.Context) (AccountRow, error)
	Select(ctx context.Context) ([]AccountRow, error)
	Exists(ctx context.Context) (bool, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (AccountRow, error)

	UpdateUsername(username string) AccountsQ
	UpdateRole(role string) AccountsQ
	UpdateVersion(version int32) AccountsQ
	UpdateSourceUpdatedAt(source time.Time) AccountsQ

	Delete(ctx context.Context) error

	FilterID(accountID uuid.UUID) AccountsQ
	FilterUsername(username string) AccountsQ
	FilterVersion(version int32) AccountsQ
}

func (r *Repository) CreateAccount(ctx context.Context, params account.CreateAccountParams) (models.Account, error) {
	acc, err := r.accountsQ.New().Insert(ctx, AccountRow{
		ID:              params.ID,
		Username:        params.Username,
		Role:            params.Role,
		Version:         params.Version,
		SourceCreatedAt: params.CreatedAt,
		SourceUpdatedAt: params.CreatedAt,
	})
	if err != nil {
		return models.Account{}, fmt.Errorf("failed to insert account, cause: %w", err)
	}

	return acc.ToModel(), nil
}

func (r *Repository) GetAccountByID(ctx context.Context, accountID uuid.UUID) (models.Account, error) {
	row, err := r.accountsQ.New().FilterID(accountID).Get(ctx)
	switch {
	case err != nil:
		return models.Account{}, fmt.Errorf("failed to get account, cause: %w", err)
	case row.IsNil():
		return models.Account{}, errx.ErrorAccountNotFound.Raise(
			fmt.Errorf("account with id %s not found", accountID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) ExistsAccountByID(ctx context.Context, accountID uuid.UUID) (bool, error) {
	exist, err := r.accountsQ.New().FilterID(accountID).Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check account existence by id %s, cause: %w", accountID, err)
	}

	return exist, nil
}

func (r *Repository) GetAccountByUsername(ctx context.Context, username string) (models.Account, error) {
	row, err := r.accountsQ.New().FilterUsername(username).Get(ctx)
	switch {
	case err != nil:
		return models.Account{}, fmt.Errorf("failed to get account by username, cause: %w", err)
	case row.IsNil():
		return models.Account{}, errx.ErrorAccountNotFound.Raise(
			fmt.Errorf("account with username %s not found", username),
		)
	}

	return row.ToModel(), nil
}
func (r *Repository) ExistsAccountByUsername(ctx context.Context, username string) (bool, error) {
	exist, err := r.accountsQ.New().FilterUsername(username).Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check account existence by username %s, cause: %w", username, err)
	}

	return exist, nil
}

func (r *Repository) UpdateAccountUsername(
	ctx context.Context,
	accountID uuid.UUID,
	params account.UpdateUsernameParams,
) (models.Account, error) {
	row, err := r.accountsQ.New().
		FilterID(accountID).
		UpdateUsername(params.Username).
		UpdateVersion(params.Version).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateOne(ctx)

	if err != nil {
		return models.Account{}, fmt.Errorf(
			"failed to update account username for account %s, cause: %w", accountID, err,
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) DeleteAccount(ctx context.Context, accountID uuid.UUID) error {
	err := r.accountsQ.New().FilterID(accountID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete account %s, cause: %w", accountID, err)
	}

	return nil
}
