package account

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
)

type CreateAccountParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`

	CreatedAt time.Time `json:"created_at"`
}

func (m *Module) Create(
	ctx context.Context,
	params CreateAccountParams,
) error {
	buried, err := m.repo.AccountIsBuried(ctx, params.ID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorAccountDeleted.Raise(
			fmt.Errorf("account with id %s is already deleted", params.ID),
		)
	}

	exist, err := m.repo.ExistsAccountByID(ctx, params.ID)
	if err != nil {
		return err
	}
	if exist {
		return errx.ErrorAccountAlreadyExists.Raise(
			fmt.Errorf("account with id %s already exists", params.ID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		_, err = m.repo.CreateAccount(ctx, params)
		if err != nil {
			return err
		}

		profile, err := m.repo.CreateProfile(ctx, params.ID, params.Username)
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
