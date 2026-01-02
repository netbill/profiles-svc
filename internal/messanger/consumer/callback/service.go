package callback

import (
	"context"

	"github.com/netbill/kafkakit/box"
	"github.com/netbill/logium"
	"github.com/segmentio/kafka-go"
)

type Inbox interface {
	CreateInboxEvent(
		ctx context.Context,
		status string,
		message kafka.Message,
	) (box.InboxEvent, error)
}

type Service struct {
	inbox Inbox
	log   logium.Logger
}

func NewService(log logium.Logger, inbox Inbox) *Service {
	return &Service{
		inbox: inbox,
		log:   log,
	}
}
