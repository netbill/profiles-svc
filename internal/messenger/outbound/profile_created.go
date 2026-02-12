package outbound

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/eventbox/headers"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (o *Outbound) WriteProfileCreated(
	ctx context.Context,
	profile models.Profile,
) error {
	payload, err := json.Marshal(contracts.ProfileCreatedPayload{
		AccountID: profile.AccountID,
		Username:  profile.Username,
		CreatedAt: profile.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal profile created payload: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: contracts.ProfilesTopicV1,
			Key:   []byte(profile.AccountID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(contracts.ProfileCreatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(contracts.ProfilesSvcGroup)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile created: %w", err)
	}

	return nil
}
