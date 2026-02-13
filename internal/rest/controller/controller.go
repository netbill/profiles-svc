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

	GetByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error)

	OpenUpdateSession(
		ctx context.Context,
		account models.AccountActor,
	) (models.UpdateProfileMedia, models.Profile, error)
	ConfirmUpdateSession(ctx context.Context,
		account models.AccountActor,
		session models.UploadScope,
		params profile.UpdateParams,
	) (models.Profile, error)
	DeleteUploadAvatar(
		ctx context.Context,
		account models.AccountActor,
		session models.UploadScope,
	) error
	CancelUpdateSession(
		ctx context.Context,
		account models.AccountActor,
		session models.UploadScope,
	) error
}

type responser interface {
	Render(w http.ResponseWriter, status int, res ...interface{})
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
