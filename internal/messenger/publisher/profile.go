package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/profiles-svc/internal/models"
)

func (p *Publisher) WriteProfileCreated(
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

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.ProfileCreatedEvent,
		Version:  1,
		Topic:    evtypes.ProfilesTopicV1,
		Key:      profile.AccountID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile created: %w", err)
	}

	return nil
}

func (p *Publisher) WriteProfileUpdated(
	ctx context.Context,
	profile models.Profile,
) error {
	payload, err := json.Marshal(evtypes.ProfileUpdatedPayload{
		AccountID:   profile.AccountID,
		Username:    profile.Username,
		Pseudonym:   profile.Pseudonym,
		Description: profile.Description,
		AvatarKey:   profile.AvatarKey,
		Version:     profile.Version,
		UpdatedAt:   profile.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal profile updated payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.ProfileUpdatedEvent,
		Version:  1,
		Topic:    evtypes.ProfilesTopicV1,
		Key:      profile.AccountID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile updated, cause: %w", err)
	}

	return nil
}

func (p *Publisher) WriteProfileDeleted(
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

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.ProfileDeletedEvent,
		Version:  1,
		Topic:    evtypes.ProfilesTopicV1,
		Key:      accountID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create outbox event for profile deleted, cause: %w", err)
	}

	return nil
}
