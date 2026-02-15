package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/restkit/pagi"
)

type Modules struct {
	Profile profileModule
}

type profileModule interface {
	GetList(
		ctx context.Context,
		params profile.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Profile], error)

	GetMy(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error)

	GetAvatarUploadMediaLinks(
		ctx context.Context,
		actor models.AccountActor,
	) (models.UploadMediaLink, models.Profile, error)
	Update(
		ctx context.Context,
		actor models.AccountActor,
		params profile.UpdateParams,
	) (profile models.Profile, err error)
	DeleteUploadAvatar(
		ctx context.Context,
		actor models.AccountActor,
		key string,
	) error
}

type responser interface {
	Status(w http.ResponseWriter, status int)
	Render(w http.ResponseWriter, status int, res interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type Controller struct {
	modules   Modules
	responser responser
}

func New(modules Modules, responser responser) *Controller {
	return &Controller{
		modules:   modules,
		responser: responser,
	}
}
