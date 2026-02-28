package controller

import (
	"context"

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

	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)

	UpdateOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error)

	CreateUploadMediaLinks(
		ctx context.Context,
		actor models.AccountActor,
	) (models.Profile, models.UploadProfileMediaLinks, error)
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

type Controller struct {
	modules Modules
}

func New(modules Modules) *Controller {
	return &Controller{
		modules: modules,
	}
}
