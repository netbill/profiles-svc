package callback

import (
	"context"
	"fmt"

	"github.com/netbill/kafkakit/box"
	"github.com/segmentio/kafka-go"
)

func (s Service) UpdateUsername(ctx context.Context, event kafka.Message) error {
	_, err := s.inbox.CreateInboxEvent(ctx, box.InboxStatusPending, event)
	if err != nil {
		s.log.Errorf("failed to processed account username change for account %s", string(event.Key))
		return fmt.Errorf("failed to processing account username change event for account %s: %w", string(event.Key), err)
	}

	return nil
}
