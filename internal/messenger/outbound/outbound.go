package outbound

import (
	"github.com/netbill/evebox/box/outbox"
	"github.com/netbill/logium"
)

type Outbound struct {
	log    logium.Logger
	outbox outbox.Box
}

func New(log logium.Logger, ob outbox.Box) *Outbound {
	return &Outbound{
		log:    log,
		outbox: ob,
	}
}
