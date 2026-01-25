package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

type UpdateParams struct {
	Username    *string
	Pseudonym   *string
	Description *string
	Official    *bool
}

func (s Service) UpdateProfile(ctx context.Context, accountID uuid.UUID, input UpdateParams) (profile models.Profile, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.UpdateProfile(ctx, accountID, input)
		if err != nil {
			return err
		}

		err = s.messanger.WriteProfileUpdated(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (s Service) UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (profile models.Profile, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.UpdateProfileUsername(ctx, accountID, username)
		if err != nil {
			return err
		}

		err = s.messanger.WriteProfileUpdated(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (s Service) UpdateProfileOfficial(ctx context.Context, accountID uuid.UUID, official bool) (profile models.Profile, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.UpdateProfileOfficial(ctx, accountID, official)
		if err != nil {
			return err
		}

		err = s.messanger.WriteProfileUpdated(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
