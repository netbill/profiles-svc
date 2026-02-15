package account

import (
	"context"

	"github.com/google/uuid"
)

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
