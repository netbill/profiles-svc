package profile

import (
	"context"

	"github.com/google/uuid"

	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Module struct {
	repo      repo
	messenger messenger
	token     token
	bucket    bucket
}

func New(repo repo, messenger messenger, token token, bucket bucket) *Module {
	return &Module{
		repo:      repo,
		messenger: messenger,
		token:     token,
		bucket:    bucket,
	}
}

type repo interface {
	GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateProfile(ctx context.Context, accountID uuid.UUID, params UpdateParams) (models.Profile, error)
	UpdateProfileAvatar(ctx context.Context, accountID uuid.UUID, avatarURL string) (models.Profile, error)
	DeleteProfileAvatar(ctx context.Context, accountID uuid.UUID) (models.Profile, error)

	UpdateProfileOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error)

	DeleteProfile(ctx context.Context, accountID uuid.UUID) error

	FilterProfiles(
		ctx context.Context,
		params FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Profile], error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteProfileUpdated(ctx context.Context, profile models.Profile) error
}

type token interface {
	GenerateUploadProfileMediaToken(
		OwnerAccountID uuid.UUID,
		UploadSessionID uuid.UUID,
	) (string, error)
}

type bucket interface {
	GetPreloadLinkForProfileMedia(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) (links models.UpdateProfileMediaLinks, err error)

	CancelUpdateProfileAvatar(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) error

	DeleteProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
	) error

	UpdateProfileAvatar(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) (string, error)

	CleanProfileMediaSession(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) error
}
