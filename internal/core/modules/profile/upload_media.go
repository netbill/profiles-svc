package profile

import (
	"context"

	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) CreateProfileUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
) (models.Profile, models.UploadProfileMediaLinks, error) {
	profile, err := m.repo.GetProfileByAccountID(ctx, actor)
	if err != nil {
		return models.Profile{}, models.UploadProfileMediaLinks{}, err
	}

	links, err := m.bucket.CreateProfileUploadMediaLinks(ctx, actor)
	if err != nil {
		return models.Profile{}, models.UploadProfileMediaLinks{}, err
	}

	return profile, links, nil
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
