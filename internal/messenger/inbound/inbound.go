package inbound

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/logium"
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
	DeleteProfile(ctx context.Context, accountID uuid.UUID) error
}
