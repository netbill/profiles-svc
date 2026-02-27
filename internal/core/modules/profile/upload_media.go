package profile

import (
	"context"
	"fmt"

	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) CreateUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
) (models.Profile, models.UploadProfileMediaLinks, error) {
	profile, err := m.repo.GetProfileByAccountID(ctx, actor)
	if err != nil {
		return models.Profile{}, models.UploadProfileMediaLinks{}, err
	}

	links, err := m.bucket.CreateProfileAvatarUploadMediaLinks(ctx, actor)
	if err != nil {
		return models.Profile{}, models.UploadProfileMediaLinks{}, err
	}

	return profile, models.UploadProfileMediaLinks{
		Avatar: links,
	}, nil
}

func (m *Module) updateProfileAvatar(
	ctx context.Context,
	profile models.Profile,
	params UpdateParams,
) (newKey *string, err error) {
	if params.AvatarKey != nil {
		if err = m.bucket.ValidateUploadProfileAvatar(ctx, profile.AccountID, *params.AvatarKey); err != nil {
			return nil, fmt.Errorf("failed to validate profile avatar: %w", err)
		}

		avatarKey, err := m.bucket.UpdateProfileAvatar(ctx, profile.AccountID, *params.AvatarKey)
		if err != nil {
			return nil, fmt.Errorf("failed to update profile avatar: %w", err)
		}

		if err = m.bucket.DeleteUploadProfileAvatar(ctx, profile.AccountID, *params.AvatarKey); err != nil {
			return nil, fmt.Errorf("failed to delete temp profile avatar: %w", err)
		}

		newKey = &avatarKey
	}

	if profile.AvatarKey != nil {
		if err = m.bucket.DeleteProfileAvatar(ctx, profile.AccountID, *profile.AvatarKey); err != nil {
			return nil, fmt.Errorf("failed to delete profile avatar: %w", err)
		}
	}

	return newKey, nil
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
