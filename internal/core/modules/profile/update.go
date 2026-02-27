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

	upd := false

	if !ptrStrEq(params.AvatarKey, profile.AvatarKey) {
		avatarKey, err := m.updateProfileAvatar(ctx, profile, params)
		if err != nil {
			return models.Profile{}, err
		}
		params.AvatarKey = avatarKey
		upd = true
	}

	if !ptrStrEq(params.Pseudonym, profile.Pseudonym) {
		upd = true
	}

	if !ptrStrEq(params.Description, profile.Description) {
		upd = true
	}

	if !upd {
		return profile, nil
	}

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

func ptrStrEq(a, b *string) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}
