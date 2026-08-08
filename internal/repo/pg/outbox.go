package pg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/evtypes"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/models"
)

const (
	outboxTable    = "outbox_events"
	outboxCols     = "event_id, topic, key, type, version, producer, payload"
	payloadVersion = 1
)

type OutboxRepo struct {
	db       *pgdbx.DB
	producer string
}

func NewOutboxRepo(db *pgdbx.DB, producer string) *OutboxRepo {
	return &OutboxRepo{db: db, producer: producer}
}

func (r *OutboxRepo) write(ctx context.Context, topic, key, eventType string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	const query = `
		INSERT INTO ` + outboxTable + ` (` + outboxCols + `)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if _, err = r.db.Exec(ctx, query, uuid.New(), topic, key, eventType, payloadVersion, r.producer, data); err != nil {
		return fmt.Errorf("insert outbox event: %w", err)
	}

	return nil
}

func (r *OutboxRepo) WriteProfileCreated(
	ctx context.Context,
	profile models.Profile,
) error {
	return r.write(
		ctx,
		evtypes.ProfilesTopicV1,
		profile.AccountID.String(),
		evtypes.ProfileCreatedEvent,
		evtypes.ProfileCreatedPayload{
			Profile: toEvProfile(profile),
		},
	)
}

func (r *OutboxRepo) WriteProfileUpdated(
	ctx context.Context,
	profile models.Profile,
) error {
	return r.write(
		ctx,
		evtypes.ProfilesTopicV1,
		profile.AccountID.String(),
		evtypes.ProfileUpdatedEvent,
		evtypes.ProfileUpdatedPayload{
			Profile: toEvProfile(profile),
		},
	)
}

func (r *OutboxRepo) WriteProfileDeleted(
	ctx context.Context,
	profile models.Profile,
) error {
	return r.write(
		ctx,
		evtypes.ProfilesTopicV1,
		profile.AccountID.String(),
		evtypes.ProfileDeletedEvent,
		evtypes.ProfileDeletedPayload{
			Profile: toEvProfile(profile),
		},
	)
}

func toEvProfile(a models.Profile) evtypes.Profile {
	return evtypes.Profile{
		AccountID:   a.AccountID,
		Username:    a.Username,
		Pseudonym:   a.Pseudonym,
		Description: a.Description,
		AvatarKey:   a.AvatarKey,
		Version:     a.Version,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		DeletedAt:   a.DeletedAt,
	}
}
