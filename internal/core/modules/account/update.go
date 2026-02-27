package account

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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
	buried, err := m.repo.AccountIsBuried(ctx, accountID)
	if err != nil {
		return err
	}
	if buried {
		return nil
	}

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

		return m.messenger.WriteProfileUpdated(ctx, profile)
	})
}
