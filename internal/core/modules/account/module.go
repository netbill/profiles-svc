package account

import (
	"context"
	"time"

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

type CreateAccountParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Version  int32     `json:"version"`

	CreatedAt time.Time `json:"created_at"`
}

func (m *Module) Create(
	ctx context.Context,
	params CreateAccountParams,
) error {
	exist, err := m.repo.ExistsAccountByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		account, err := m.repo.CreateAccount(ctx, params)
		if err != nil {
			return err
		}

		profile, err := m.repo.CreateProfile(ctx, account.ID, account.Username)
		if err != nil {
			return err
		}

		err = m.messenger.WriteProfileCreated(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	})
}

type UpdateUsernameParams struct {
	Username  string
	Version   int32
	UpdatedAt time.Time
}

func (m *Module) UpdateUsername(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateUsernameParams,
) error {
	account, err := m.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account.Version >= params.Version {
		return nil
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		_, err = m.repo.UpdateAccountUsername(ctx, accountID, params)
		if err != nil {
			return err
		}

		profile, err := m.repo.UpdateProfileUsername(ctx, accountID, params.Username)
		if err != nil {
			return err
		}

		err = m.messenger.WriteProfileUpdated(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	})
}

func (m *Module) Delete(ctx context.Context, accountID uuid.UUID) error {
	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err := m.repo.DeleteProfile(ctx, accountID)
		if err != nil {
			return err
		}

		err = m.repo.DeleteAccount(ctx, accountID)
		if err != nil {
			return err
		}

		err = m.messenger.WriteProfileDeleted(ctx, accountID)
		if err != nil {
			return err
		}

		return nil
	})
}
