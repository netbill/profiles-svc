package account

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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
