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

func (p Outbound) WriteProfileCreated(
	ctx context.Context,
	profile models.Profile,
) error {
	payload, err := json.Marshal(contracts.AccountProfileCreatedPayload{
		Data: contracts.AccountProfileCreatedPayloadData{
			AccountID:   profile.AccountID,
			Username:    profile.Username,
			Official:    profile.Official,
			Pseudonym:   profile.Pseudonym,
			Description: profile.Description,
			Avatar:      profile.Avatar,
		},
		Timestamp: profile.CreatedAt,
	})
	if err != nil {
		return err
	}

	eventID := uuid.New().String()

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.AccountsTopicV1,
			Key:   []byte(profile.AccountID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(eventID)}, // Outbox will fill this
				{Key: header.EventType, Value: []byte(contracts.ProfileCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.ProfilesSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	p.log.Debugf("profile created event queued, account_id: %s, event_id: %s", profile.AccountID, eventID)

	return err
}
