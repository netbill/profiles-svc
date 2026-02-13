package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) OpenUpdateSession(
	ctx context.Context,
	actor models.AccountActor,
) (models.UpdateProfileMedia, models.Profile, error) {
	profile, err := m.GetByAccountID(ctx, actor)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, err
	}

	uploadSessionID := uuid.New()
	links, err := m.bucket.GetPreloadLinkForProfileMedia(
		ctx,
		actor,
		uploadSessionID,
	)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, err
	}

	uploadToken, err := m.token.GenerateUploadProfileMediaToken(actor, uploadSessionID)
	if err != nil {
		return models.UpdateProfileMedia{}, models.Profile{}, err
	}

	return models.UpdateProfileMedia{
		Links:           links,
		UploadSessionID: uploadSessionID,
		UploadToken:     uploadToken,
	}, profile, nil
}

type UpdateMediaParams struct {
	DeleteAvatar bool
	avatarKey    *string
}

type UpdateParams struct {
	Pseudonym   *string
	Description *string
	Media       UpdateMediaParams
}

func (p UpdateParams) GetUpdatedAvatar() *string {
	if p.Media.DeleteAvatar {
		return nil
	}
	return p.Media.avatarKey
}

func (m *Module) ConfirmUpdateSession(
	ctx context.Context,
	actor models.AccountActor,
	scope models.UploadScope,
	params UpdateParams,
) (profile models.Profile, err error) {
	profile, err = m.GetByAccountID(ctx, actor)
	if err != nil {
		return models.Profile{}, err
	}

	if params.Media.DeleteAvatar {
		if err = m.bucket.DeleteProfileAvatar(ctx, actor); err != nil {
			return models.Profile{}, err
		}
		params.Media.avatarKey = nil
	} else {
		key, err := m.bucket.UpdateProfileAvatar(ctx, actor, scope)
		switch {
		case errors.Is(err, errx.ErrorNoContentUploaded):
			params.Media.avatarKey = profile.Avatar
		case err != nil:
			return models.Profile{}, err
		default:
			params.Media.avatarKey = &key
		}
	}

	err = m.bucket.CleanProfileMediaSession(ctx, actor, scope)
	if err != nil {
		return models.Profile{}, err
	}

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
	scope models.UploadScope,
) error {
	err := m.bucket.CancelUpdateProfileAvatar(ctx, actor, scope)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) CancelUpdateSession(
	ctx context.Context,
	actor models.AccountActor,
	scope models.UploadScope,
) error {
	err := m.bucket.CleanProfileMediaSession(ctx, actor, scope)
	if err != nil {
		return err
	}

	return nil
}
