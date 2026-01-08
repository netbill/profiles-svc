package outbound

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (p Outbound) WriteProfileUpdated(
	ctx context.Context,
	profile models.Profile,
) error {
	payload, err := json.Marshal(contracts.ProfileUpdatedPayload{
		Profile: profile,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.ProfilesTopicV1,
			Key:   []byte(profile.AccountID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: header.EventType, Value: []byte(contracts.ProfileUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.ProfilesSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
