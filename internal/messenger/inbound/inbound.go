package inbound

import (
	"context"

	"github.com/google/uuid"
)

type Inbound struct {
	modules
}

type modules struct {
	profile profileModule
}

func New(profileModule profileModule) *Inbound {
	return &Inbound{
		modules: modules{
			profile: profileModule,
		},
	}
}

type profileModule interface {
	Create(ctx context.Context, userID uuid.UUID, username string) error
	UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) error
	Delete(ctx context.Context, accountID uuid.UUID) error
}
