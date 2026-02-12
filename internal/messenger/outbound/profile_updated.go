package outbound

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/eventbox/headers"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

func (o *Outbound) WriteProfileUpdated(
	ctx context.Context,
	profile models.Profile,
) error {
	payload, err := json.Marshal(evtypes.ProfileUpdatedPayload{
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

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.ProfilesTopicV1,
			Key:   []byte(profile.AccountID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.ProfileUpdatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(evtypes.ProfilesSvcGroup)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile updated, cause: %w", err)
	}

	return nil
}
