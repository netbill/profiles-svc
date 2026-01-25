package outbound

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (o Outbound) WriteProfileUpdated(
	ctx context.Context,
	profile models.Profile,
) error {
	payload, err := json.Marshal(contracts.ProfileUpdatedPayload{
		AccountID:   profile.AccountID,
		Username:    profile.Username,
		Official:    profile.Official,
		Pseudonym:   profile.Pseudonym,
		Description: profile.Description,
		UpdatedAt:   profile.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal profile updated payload, cause: %w", err)
	}

	event, err := o.outbox.CreateOutboxEvent(
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
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile updated, cause: %w", err)
	}

	o.log.Debugf("profile updated event queued, account_id: %s, event_id: %s", profile.AccountID, event.ID)

	return nil
}
