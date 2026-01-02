package consumer

import (
	"context"

	"github.com/netbill/kafkakit/subscriber"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal/messanger/contracts"
	"github.com/segmentio/kafka-go"
)

type Service struct {
	log       logium.Logger
	addr      []string
	callbacks callbacks
}

type callbacks interface {
	CreateAccount(ctx context.Context, event kafka.Message) error
	UpdateUsername(ctx context.Context, event kafka.Message) error
}

func New(log logium.Logger, addr []string, callbacks callbacks) *Service {
	return &Service{
		addr:      addr,
		log:       log,
		callbacks: callbacks,
	}
}

func (s Service) Run(ctx context.Context) {
	sub := subscriber.New(s.addr, contracts.AccountsTopicV1, contracts.GroupProfilesSvc)

	s.log.Info("starting events consumer", "addr", s.addr)

	go func() {
		err := sub.Consume(ctx, func(m kafka.Message) (subscriber.HandlerFunc, bool) {
			et, ok := subscriber.Header(m, "event_type")
			if !ok {
				return nil, false
			}

			switch et {
			case contracts.AccountCreatedEvent:
				return s.callbacks.CreateAccount, true
			case contracts.AccountUsernameChangeEvent:
				return s.callbacks.UpdateUsername, true
			default:
				return nil, false
			}
		})
		if err != nil {
			s.log.Warnf("accounts consumer stopped: %v", err)
		}
	}()
}
