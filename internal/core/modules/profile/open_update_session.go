package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) OpenUpdateSession(
	ctx context.Context,
	accountID uuid.UUID,
) (models.UpdateProfileMedia, models.Profile, error) {
	profile, err := m.GetByAccountID(ctx, accountID)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, err
	}

	uploadSessionID := uuid.New()
	links, err := m.bucket.GetPreloadLinkForProfileMedia(
		ctx,
		accountID,
		uploadSessionID,
	)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, err
	}

	uploadToken, err := m.token.GenerateUploadProfileMediaToken(accountID, uploadSessionID)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, err
	}

	return models.UpdateProfileMedia{
		Links:           links,
		UploadSessionID: uploadSessionID,
		UploadToken:     uploadToken,
	}, profile, nil
}

type UpdateParams struct {
	Pseudonym   *string
	Description *string

	Media UpdateMediaParams
}

type UpdateMediaParams struct {
	UploadSessionID uuid.UUID

	DeleteAvatar bool
	avatarKey    *string
}

func (p UpdateParams) GetUpdatedAvatar() *string {
	if p.Media.DeleteAvatar {
		return nil
	}

	return p.Media.avatarKey
}
