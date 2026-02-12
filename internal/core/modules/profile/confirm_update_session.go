package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) ConfirmUpdateSession(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateParams,
) (profile models.Profile, err error) {
	profile, err = m.GetByAccountID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}

	params.Media.avatarKey = profile.Avatar
	switch params.Media.DeleteAvatar {
	case true:
		if err = m.bucket.DeleteProfileAvatar(
			ctx,
			accountID,
		); err != nil {
			return models.Profile{}, err
		}

		params.Media.avatarKey = nil
	case false:
		avatar, err := m.bucket.AcceptUpdateProfileMedia(
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

	err = m.bucket.CleanProfileMediaSession(
		ctx,
		accountID,
		params.Media.UploadSessionID,
	)
	if err != nil {
		return models.Profile{}, err
	}

	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = m.repo.UpdateProfile(ctx, accountID, params)
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
