package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/modules/account"
	"github.com/netbill/profiles-svc/pkg/log"
)

type Handler struct {
	log     *log.Logger
	modules Modules
}

type Modules struct {
	Account accountModule
}

func New(log *log.Logger, modules Modules) *Handler {
	return &Handler{
		log:     log,
		modules: modules,
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
