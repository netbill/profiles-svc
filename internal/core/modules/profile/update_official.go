package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) UpdateOfficial(ctx context.Context, accountID uuid.UUID, official bool) (profile models.Profile, err error) {
	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = m.repo.UpdateProfileOfficial(ctx, accountID, official)
		if err != nil {
			return err
		}

		err = m.messenger.WriteProfileUpdated(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
