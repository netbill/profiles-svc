package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (s Service) OpenProfileUpdateSession(
	ctx context.Context,
	accountID uuid.UUID,
) (models.UpdateProfileMedia, models.Profile, error) {
	profile, err := s.GetProfileByAccountID(ctx, accountID)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, err
	}

	uploadSessionID := uuid.New()
	uploadURL, getURL, err := s.bucket.GetPreloadLinkForUpdateProfileAvatar(
		ctx,
		accountID,
		uploadSessionID,
	)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, fmt.Errorf(
			"failed to get preload link for avatarKey upload url: %w", err,
		)
	}

	uploadToken, err := s.token.NewUploadProfileMediaToken(accountID, uploadSessionID)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, fmt.Errorf(
			"failed to generate upload token for avatarKey upload url: %w", err,
		)
	}

	return models.UpdateProfileMedia{
		UploadSessionID: uploadSessionID,
		UploadURL:       uploadURL,
		GetURL:          getURL,
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

func (s Service) UpdateProfile(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateParams,
) (profile models.Profile, err error) {
	profile, err = s.GetProfileByAccountID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}

	params.Media.avatarKey = profile.Avatar
	switch params.Media.DeleteAvatar {
	case true:
		err = s.bucket.DeleteProfileAvatar(
			ctx,
			accountID,
		)
		if err != nil {
			return models.Profile{}, err
		}

		params.Media.avatarKey = nil
	case false:
		avatar, err := s.bucket.AcceptUpdateProfileAvatar(
			ctx,
			accountID,
			params.Media.UploadSessionID,
		)
		switch {
		case errors.Is(err, errx.ErrorNoContentUploaded):
			// No new avatar uploaded, keep the existing one
		case err != nil:
			return models.Profile{}, err
		default:
			params.Media.avatarKey = &avatar
		}
	}

	err = s.bucket.CleanProfileMediaSession(
		ctx,
		accountID,
		params.Media.UploadSessionID,
	)
	if err != nil {
		return models.Profile{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.UpdateProfile(ctx, accountID, params)
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

func (s Service) DeleteUploadProfileAvatarInSession(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) error {
	err := s.bucket.CancelUpdateProfileAvatar(ctx, accountID, sessionID)
	if err != nil {
		return err
	}

	return nil
}
