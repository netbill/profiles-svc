package callbacker

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal/core/models"
)

type Callbacker struct {
	log    logium.Logger
	domain core
}

func New(log logium.Logger, core core) Callbacker {
	return Callbacker{
		log:    log,
		domain: core,
	}
}

type core interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	DeleteProfile(
		ctx context.Context,
		accountID uuid.UUID,
	) error
}
