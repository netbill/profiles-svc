package profile

import (
	"context"

	"github.com/netbill/profiles-svc/internal/core/models"
)

type UpdateParams struct {
	AvatarKey   *string
	Pseudonym   *string
	Description *string
}

func (m *Module) Update(
	ctx context.Context,
	actor models.AccountActor,
	params UpdateParams,
) (profile models.Profile, err error) {
	profile, err = m.repo.GetProfileByAccountID(ctx, actor)
	if err != nil {
		return models.Profile{}, err
	}

	avatarKey, err := m.bucket.UpdateProfileAvatar(ctx, actor, profile.Avatar, params.AvatarKey)
	if err != nil {
		return models.Profile{}, err
	}

	params.AvatarKey = avatarKey

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = m.repo.UpdateProfile(ctx, actor, params)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteProfileUpdated(ctx, profile); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
