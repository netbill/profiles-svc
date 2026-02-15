package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/modules/account"
)

type Inbound struct {
	modules
}

type modules struct {
	account accountModule
}

func New(profileModule accountModule) *Inbound {
	return &Inbound{
		modules: modules{
			account: profileModule,
		},
	}
}

type accountModule interface {
	Create(ctx context.Context, params account.CreateAccountParams) error
	UpdateUsername(
		ctx context.Context,
		accountID uuid.UUID,
		module account.UpdateUsernameParams,
	) error
	Delete(ctx context.Context, accountID uuid.UUID) error
}
