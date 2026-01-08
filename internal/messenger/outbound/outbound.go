package outbound

import (
	"github.com/netbill/evebox/box/outbox"
	"github.com/netbill/logium"
)

type Outbound struct {
	log    logium.Logger
	addr   []string
	outbox outbox.Box
}

func New(log logium.Logger, ob outbox.Box, addr ...string) *Outbound {
	return &Outbound{
		log:    log,
		addr:   addr,
		outbox: ob,
	}
}
