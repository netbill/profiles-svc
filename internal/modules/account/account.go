package account

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
)

type accountMessenger interface {
	WriteProfileCreated(ctx context.Context, profile models.Profile) error
	WriteProfileUpdated(ctx context.Context, profile models.Profile) error
	WriteProfileDeleted(ctx context.Context, accountID uuid.UUID) error
}

type Service struct {
	account   accountRepo
	profile   profileRepo
	tombstone tombstoneRepo
	tx        transaction

	messenger accountMessenger
}

type ServiceDeps struct {
	AccountRepo accountRepo
	ProfileRepo profileRepo
	Tombstone   tombstoneRepo
	Transaction transaction

	Messenger accountMessenger
}

func NewAccountModule(deps ServiceDeps) *Service {
	return &Service{
		account:   deps.AccountRepo,
		profile:   deps.ProfileRepo,
		tombstone: deps.Tombstone,
		tx:        deps.Transaction,

		messenger: deps.Messenger,
	}
}

type CreateAccountParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`

	CreatedAt time.Time `json:"created_at"`
}

func (m *Service) Create(
	ctx context.Context,
	params CreateAccountParams,
) error {
	buried, err := m.tombstone.AccountIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorAccountDeleted.Raise(
			fmt.Errorf("account with id %s is already deleted", params.ID),
		)
	}

	exist, err := m.account.ExistsByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if exist {
		return errx.ErrorAccountAlreadyExists.Raise(
			fmt.Errorf("account with id %s already exists", params.ID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		_, err = m.account.Create(ctx, params)
		if err != nil {
			return err
		}

		profile, err := m.profile.Create(ctx, params.ID, params.Username)
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

func (m *Service) UpdateUsername(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateUsernameParams,
) error {
	buried, err := m.tombstone.AccountIsBuried(ctx, accountID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorAccountDeleted.Raise(
			fmt.Errorf("account with id %s is already deleted", accountID),
		)
	}

	account, err := m.account.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account.Version >= params.Version {
		return nil
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		_, err = m.account.UpdateUsername(ctx, accountID, params)
		if err != nil {
			return err
		}

		profile, err := m.profile.UpdateUsername(ctx, accountID, params.Username)
		if err != nil {
			return err
		}

		return m.messenger.WriteProfileUpdated(ctx, profile)
	})
}

func (m *Service) Delete(ctx context.Context, accountID uuid.UUID) error {
	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := m.tombstone.BuryAccount(ctx, accountID); err != nil {
			return err
		}

		err := m.profile.Delete(ctx, accountID)
		if err != nil {
			return err
		}

		err = m.account.Delete(ctx, accountID)
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
