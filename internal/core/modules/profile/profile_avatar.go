package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (s Service) GetPreloadLinkForUpdateAvatar(
	ctx context.Context,
	accountID uuid.UUID,
) (models.UpdateProfileAvatar, error) {
	sessionID := uuid.New()
	uploadURL, getURL, err := s.bucket.GetPreloadLinkForUpdateProfileAvatar(
		ctx,
		accountID,
		sessionID,
	)
	if err != nil {
		return models.UpdateProfileAvatar{}, fmt.Errorf("failed to get preload link for avatar upload url: %w", err)
	}

	uploadToken, err := s.token.NewUploadProfileAvatarToken(sessionID)
	if err != nil {
		return models.UpdateProfileAvatar{}, fmt.Errorf("failed to generate upload token for avatar upload url: %w", err)

	}

	return models.UpdateProfileAvatar{
		UploadURL:   uploadURL,
		GetURL:      getURL,
		UploadToken: uploadToken,
	}, nil
}

func (s Service) AcceptUpdateAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) (models.Profile, error) {
	link, err := s.bucket.AcceptUpdateProfileAvatar(ctx, accountID, sessionID)
	if err != nil {
		return models.Profile{}, err
	}

	var profile models.Profile
	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		profile, err = s.repo.UpdateProfileAvatar(ctx, accountID, link)
		if err != nil {
			return err
		}

		err = s.messanger.WriteProfileUpdated(txCtx, profile)
		if err != nil {
			return fmt.Errorf("failed to send profile updated message: %w", err)
		}

		return nil
	})

	return profile, err
}

func (s Service) CancelUpdateAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) (models.Profile, error) {
	err := s.bucket.CancelUpdateProfileAvatar(ctx, accountID, sessionID)
	if err != nil {
		return models.Profile{}, err
	}

	return s.GetProfileByID(ctx, accountID)
}

func (s Service) DeleteProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	var profile models.Profile
	err := s.repo.Transaction(ctx, func(txCtx context.Context) error {
		var err error
		profile, err = s.repo.DeleteProfileAvatar(ctx, accountID)
		if err != nil {
			return err
		}

		err = s.messanger.WriteProfileUpdated(txCtx, profile)
		if err != nil {
			return err
		}

		return nil
	})

	return profile, err
}
