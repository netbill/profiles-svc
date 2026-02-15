package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

type Module struct {
	repo      repo
	messenger messenger
}

func New(db repo, messenger messenger) *Module {
	return &Module{
		repo:      db,
		messenger: messenger,
	}
}

type repo interface {
	CreateAccount(
		ctx context.Context,
		params CreateAccountParams,
	) (models.Account, error)
	ExistsAccountByID(ctx context.Context, accountID uuid.UUID) (bool, error)
	GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	UpdateAccountUsername(
		ctx context.Context,
		accountID uuid.UUID,
		params UpdateUsernameParams,
	) (models.Account, error)
	DeleteAccount(ctx context.Context, accountID uuid.UUID) error

	CreateProfile(
		ctx context.Context,
		accountID uuid.UUID,
		username string,
	) (models.Profile, error)
	GetAccountByID(ctx context.Context, accountID uuid.UUID) (models.Account, error)
	ExistsProfileByID(ctx context.Context, accountID uuid.UUID) (bool, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	DeleteProfile(ctx context.Context, accountID uuid.UUID) error

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteProfileCreated(ctx context.Context, profile models.Profile) error
	WriteProfileUpdated(ctx context.Context, profile models.Profile) error
	WriteProfileDeleted(ctx context.Context, accountID uuid.UUID) error
}
