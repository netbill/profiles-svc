package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/rest/contexter"
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

	GetByAccountID(ctx context.Context, userID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error)

	ConfirmUpdateSession(
		ctx context.Context,
		accountID uuid.UUID,
		params profile.UpdateParams,
	) (models.Profile, error)
	OpenUpdateSession(
		ctx context.Context,
		accountID uuid.UUID,
	) (models.UpdateProfileMedia, models.Profile, error)
	DeleteUploadAvatar(
		ctx context.Context,
		accountID, sessionID uuid.UUID,
	) error
	CancelUpdateSession(
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

	modules   Modules
	responser responser
}

func New(log *logium.Logger, responser responser, modules Modules) *Controller {
	return &Controller{
		log:       log,
		modules:   modules,
		responser: responser,
	}
}

func (c *Controller) Log(r *http.Request) *logium.Entry {
	log := c.log.WithRequest(r)

	initiator, err := contexter.AccountData(r.Context())
	if err == nil {
		log = log.WithAccount(initiator)
	}

	upload, err := contexter.UploadContentData(r.Context())
	if err == nil {
		log = log.WithUploadSession(upload)
	}

	return log
}
