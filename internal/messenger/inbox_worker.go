package messenger

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/profiles-svc/pkg/log"
)

type handlers interface {
	AccountCreated(
		ctx context.Context,
		event eventbox.InboxEvent,
	) error
	AccountDeleted(
		ctx context.Context,
		event eventbox.InboxEvent,
	) error
	AccountUsernameUpdated(
		ctx context.Context,
		event eventbox.InboxEvent,
	) error
}

func NewInboxWorker(
	logger *log.Logger,
	inbox eventbox.Inbox,
	cfg eventbox.InboxWorkerConfig,
	handlers handlers,
) *eventbox.InboxWorker {
	id := uuid.New().String()

	worker := eventbox.NewInboxWorker(id, logger, inbox, cfg)

	worker.Route(evtypes.AccountCreatedEvent, handlers.AccountCreated)
	worker.Route(evtypes.AccountDeletedEvent, handlers.AccountDeleted)
	worker.Route(evtypes.AccountUsernameUpdatedEvent, handlers.AccountUsernameUpdated)

	return worker
}
