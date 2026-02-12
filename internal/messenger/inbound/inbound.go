package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

type Inbound struct {
	domain domain
}

func New(domain domain) *Inbound {
	return &Inbound{
		domain: domain,
	}
}

type domain interface {
	Create(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)
	UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	Delete(ctx context.Context, accountID uuid.UUID) error
}
