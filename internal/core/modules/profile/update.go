package profile

import (
	"context"
	"errors"

	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) GetAvatarUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
) (models.UploadMediaLink, models.Profile, error) {
	profile, err := m.repo.GetProfileByAccountID(ctx, actor)
	if err != nil {
		return models.UploadMediaLink{}, models.Profile{}, err
	}

	links, err := m.bucket.GetPreloadLinkForProfileAvatar(ctx, actor)
	if err != nil {
		return models.UploadMediaLink{}, models.Profile{}, err
	}

	return links, profile, nil
}

type UpdateParams struct {
	AvatarKey   string
	Pseudonym   string
	Description string
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

	key, err := m.bucket.UpdateProfileAvatar(ctx, actor, params.AvatarKey)
	if err != nil {
		if errors.Is(err, errx.ErrorNoContentUploaded) {
			key = ""
		} else {
			return models.Profile{}, err
		}
	}

	params.AvatarKey = key

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = m.repo.UpdateProfile(ctx, actor, params)
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

func (m *Module) DeleteUploadAvatar(
	ctx context.Context,
	actor models.AccountActor,
	key string,
) error {
	err := m.bucket.DeleteUploadProfileAvatar(ctx, actor, key)
	if err != nil {
		return err
	}

	return nil
}
