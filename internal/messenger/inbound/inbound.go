package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal/core/models"
)

type Inbound struct {
	log    logium.Logger
	domain domain
}

func New(log logium.Logger, domain domain) Inbound {
	return Inbound{
		log:    log,
		domain: domain,
	}
}

type domain interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	DeleteProfile(ctx context.Context, accountID uuid.UUID) error
}
