package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
)

func (m *Module) Create(ctx context.Context, accountID uuid.UUID, username string) error {
	profile, err := m.repo.GetProfileByAccountID(ctx, accountID)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		// continue to create profile
	case err != nil:
		return err
	default:
		return nil
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = m.repo.InsertProfile(ctx, accountID, username)
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
