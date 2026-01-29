package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/restkit/pagi"
)

type domain interface {
	FilterProfile(
		ctx context.Context,
		params profile.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Profile], error)

	GetProfileByAccountID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateProfileOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)

	UpdateProfile(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (models.Profile, error)
	OpenProfileUpdateSession(
		ctx context.Context,
		accountID uuid.UUID,
	) (models.UpdateProfileMedia, models.Profile, error)
	DeleteUploadProfileAvatarInSession(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) error
}

type Service struct {
	domain domain
	log    *logium.Logger
}

func New(log *logium.Logger, profile domain) Service {
	return Service{
		domain: profile,
		log:    log,
	}
}
