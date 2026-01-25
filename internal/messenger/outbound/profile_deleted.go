package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/profiles-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (o Outbound) WriteProfileDeleted(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.ProfileDeletedPayload{
		AccountID: accountID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal profile deleted payload, cause: %w", err)
	}

	event, err := o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.ProfilesTopicV1,
			Key:   []byte(accountID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: header.EventType, Value: []byte(contracts.ProfileDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.ProfilesSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile deleted, cause: %w", err)
	}

	o.log.Debugf("profile deleted event queued, account_id: %s, event_id: %s", accountID.String(), event.ID)

	return nil
}
