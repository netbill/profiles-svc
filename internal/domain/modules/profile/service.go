package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/pagi"
	"github.com/netbill/profiles-svc/internal/domain/models"
)

type Service struct {
	repo      repo
	messanger messanger
}

func New(db repo, messanger messanger) Service {
	return Service{
		repo:      db,
		messanger: messanger,
	}
}

type repo interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)

	GetProfileByAccountID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateProfile(ctx context.Context, userID uuid.UUID, params UpdateParams) (models.Profile, error)

	UpdateProfileUsername(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)
	UpdateProfileOfficial(ctx context.Context, userID uuid.UUID, official bool) (models.Profile, error)

	DeleteProfile(ctx context.Context, userID uuid.UUID) error

	FilterProfiles(
		ctx context.Context,
		params FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Profile], error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messanger interface {
	WriteProfileUpdated(ctx context.Context, profile models.Profile) error
}
