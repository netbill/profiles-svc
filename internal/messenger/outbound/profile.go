package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox/headers"
	"github.com/netbill/profiles-svc/internal/core/models"
	"github.com/netbill/profiles-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

func (o *Outbound) WriteProfileCreated(
	ctx context.Context,
	profile models.Profile,
) error {
	payload, err := json.Marshal(evtypes.ProfileCreatedPayload{
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
			Topic: evtypes.ProfilesTopicV1,
			Key:   []byte(profile.AccountID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.ProfileCreatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile created: %w", err)
	}

	return nil
}

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
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile updated, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteProfileDeleted(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	payload, err := json.Marshal(evtypes.ProfileDeletedPayload{
		AccountID: accountID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal profile deleted payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.ProfilesTopicV1,
			Key:   []byte(accountID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.ProfileDeletedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile deleted, cause: %w", err)
	}

	return nil
}
