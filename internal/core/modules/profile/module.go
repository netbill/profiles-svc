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
	bucket    bucket
}

func New(repo repo, messenger messenger, bucket bucket) *Module {
	return &Module{
		repo:      repo,
		messenger: messenger,
		bucket:    bucket,
	}
}

type repo interface {
	GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateProfile(
		ctx context.Context,
		accountID uuid.UUID,
		params UpdateParams,
	) (models.Profile, error)

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

type bucket interface {
	CreateProfileUploadMediaLinks(
		ctx context.Context,
		accountID uuid.UUID,
	) (models.UploadProfileMediaLinks, error)

	ValidateProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
		tempKey string,
	) error

	DeleteUploadProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
		tempKey string,
	) error

	DeleteProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
		finalKey string,
	) error

	UpdateProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
		oldFinalKey *string,
		tempKey *string,
	) (*string, error)
}
