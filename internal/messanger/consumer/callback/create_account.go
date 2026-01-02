package callback

import (
	"context"
	"fmt"

	"github.com/netbill/kafkakit/box"
	"github.com/segmentio/kafka-go"
)

func (s Service) CreateAccount(ctx context.Context, event kafka.Message) error {
	_, err := s.inbox.CreateInboxEvent(ctx, box.InboxStatusPending, event)
	if err != nil {
		s.log.Errorf("failed to upsert inbox event for account %s: %v", string(event.Key), err)
		return fmt.Errorf("failed to processing create account event for account %s: %w", string(event.Key), err)
	}

	return nil
}
