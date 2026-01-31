package controller

import (
	"context"
	"net/http"

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

type responser interface {
	Render(w http.ResponseWriter, status int, res ...interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type Controller struct {
	log *logium.Logger

	domain    domain
	responser responser
}

func New(log *logium.Logger, profile domain, responser responser) Controller {
	return Controller{
		domain:    profile,
		log:       log,
		responser: responser,
	}
}
