package outbound

import (
	"github.com/netbill/eventbox"
	"github.com/netbill/profiles-svc/internal/messenger"
)

type Outbound struct {
	groupID string
	outbox  eventbox.Producer
}

func New(producer eventbox.Producer) *Outbound {
	return &Outbound{
		groupID: messenger.ProfilesSvcGroup,
		outbox:  producer,
	}
}
