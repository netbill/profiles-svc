package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/kafkakit/box"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal/domain/models"
	"github.com/netbill/profiles-svc/internal/messanger/contracts"
)

type InboxWorker struct {
	log    logium.Logger
	inbox  inbox
	domain domain
}

type inbox interface {
	GetInboxEventByID(
		ctx context.Context,
		id uuid.UUID,
	) (box.InboxEvent, error)

	GetPendingInboxEvents(
		ctx context.Context,
		limit int32,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsProcessed(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsFailed(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsPending(
		ctx context.Context,
		ids []uuid.UUID,
		delay time.Duration,
	) ([]box.InboxEvent, error)
}

type domain interface {
	CreateProfile(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error)
	UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
}

func NewInboxWorker(
	log logium.Logger,
	inbox inbox,
	domain domain,
) InboxWorker {
	return InboxWorker{
		log:    log,
		inbox:  inbox,
		domain: domain,
	}
}

const eventInboxRetryDelay = 1 * time.Minute

func (w InboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		events, err := w.inbox.GetPendingInboxEvents(ctx, 10)
		if err != nil {
			w.log.Errorf("failed to get pending inbox events, cause: %v", err)
			continue
		}
		if len(events) == 0 {
			continue
		}

		var processed []uuid.UUID
		var delayed []uuid.UUID

		for _, ev := range events {
			w.log.Infof("processing inbox event: %s, type %s", ev.ID, ev.Type)

			key, err := uuid.Parse(ev.Key)
			if err != nil {
				w.log.Errorf("bad inbox event key, id: %s, key: %s, error: %v", ev.ID, ev.Key, err)
				processed = append(processed, ev.ID)
				continue
			}

			switch ev.Type {
			case contracts.AccountCreatedEvent:
				var p contracts.AccountCreatedPayload
				if err = json.Unmarshal(ev.Payload, &p); err != nil {
					w.log.Errorf("bad payload for %s, id: %s, error: %v", ev.Type, ev.ID, err)
					processed = append(processed, ev.ID)
					continue
				}

				if _, err = w.domain.CreateProfile(ctx, key, p.Account.Username); err != nil {
					w.log.Errorf("failed to create profile, id: %s, error: %v", ev.ID, err)
					delayed = append(delayed, ev.ID)
					continue
				}
				processed = append(processed, ev.ID)

			case contracts.AccountUsernameChangeEvent:
				var p contracts.AccountUsernameChangePayload
				if err = json.Unmarshal(ev.Payload, &p); err != nil {
					w.log.Errorf("bad payload for %s, id: %s, error: %v", ev.Type, ev.ID, err)
					processed = append(processed, ev.ID)
					continue
				}

				if _, err = w.domain.UpdateProfileUsername(ctx, key, p.Account.Username); err != nil {
					w.log.Errorf("failed to update profile username, id: %s, error: %v", ev.ID, err)
					delayed = append(delayed, ev.ID)
					continue
				}
				processed = append(processed, ev.ID)

			default:
				w.log.Warnf("unknown inbox event type: %s, id: %s", ev.Type, ev.ID)
				processed = append(processed, ev.ID)
			}
		}

		if len(processed) > 0 {
			_, err = w.inbox.MarkInboxEventsAsProcessed(ctx, processed)
			if err != nil {
				w.log.Errorf("failed to mark inbox events as processed, ids: %v, error: %v", processed, err)
			}
		}

		if len(delayed) > 0 {
			_, err = w.inbox.MarkInboxEventsAsPending(ctx, delayed, eventInboxRetryDelay)
			if err != nil {
				w.log.Errorf("failed to delay inbox events, ids: %v, error: %v", delayed, err)
			}
		}
	}
}
