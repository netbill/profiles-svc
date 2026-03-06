package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/models"
)

type media interface {
	CreateProfileAvatarUploadMediaLinks(
		ctx context.Context,
		accountID uuid.UUID,
	) (models.UploadMediaLink, error)

	DeleteUploadProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
		key string,
	) error

	DeleteProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
		key string,
	) error

	UpdateProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
		key string,
	) (string, error)
}

func (s *Service) CreateUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
) (models.Profile, models.UploadProfileMediaLinks, error) {
	profile, err := s.repo.GetByID(ctx, actor)
	if err != nil {
		return models.Profile{}, models.UploadProfileMediaLinks{}, err
	}

	links, err := s.bucket.CreateProfileAvatarUploadMediaLinks(ctx, actor)
	if err != nil {
		return models.Profile{}, models.UploadProfileMediaLinks{}, err
	}

	return profile, models.UploadProfileMediaLinks{
		Avatar: links,
	}, nil
}

type DeleteUploadMediaParams struct {
	Avatar *string
}

func (s *Service) DeleteUploadMedia(
	ctx context.Context,
	actor models.AccountActor,
	params DeleteUploadMediaParams,
) error {
	if params.Avatar != nil {
		if err := s.bucket.DeleteUploadProfileAvatar(ctx, actor, *params.Avatar); err != nil {
			return fmt.Errorf("failed to delete upload profile avatar: %w", err)
		}
	}

	return nil
}
