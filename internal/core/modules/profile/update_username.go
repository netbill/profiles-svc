package profile

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) error {
	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err := m.repo.UpdateProfileUsername(ctx, accountID, username)
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
