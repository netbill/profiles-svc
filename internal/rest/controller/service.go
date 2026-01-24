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

	GetProfileByID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateProfile(ctx context.Context, accountID uuid.UUID, input profile.UpdateParams) (models.Profile, error)
	UpdateProfileOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)

	GetPreloadLinkForUpdateAvatar(
		ctx context.Context,
		accountID uuid.UUID,
	) (models.UpdateProfileAvatar, error)
	CancelUpdateAvatar(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) (models.Profile, error)
	AcceptUpdateAvatar(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) (models.Profile, error)
	DeleteProfileAvatar(
		ctx context.Context,
		accountID uuid.UUID,
	) (models.Profile, error)
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
