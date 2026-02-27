package account

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) Delete(ctx context.Context, accountID uuid.UUID) error {
	buried, err := m.repo.AccountIsBuried(ctx, accountID)
	if err != nil {
		return err
	}
	if buried {
		return nil
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err := m.repo.BuryAccount(ctx, accountID); err != nil {
			return err
		}

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
