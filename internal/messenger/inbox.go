package messenger

import (
	"context"
	"fmt"
	"os"

	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

type handlers interface {
	AccountCreated(
		ctx context.Context,
		message kafka.Message,
	) error
	AccountDeleted(
		ctx context.Context,
		message kafka.Message,
	) error
	AccountUsernameUpdated(
		ctx context.Context,
		message kafka.Message,
	) error
}

type Inbox struct {
	log *logium.Entry
	db  *pgdbx.DB

	handlers handlers
	config   eventpg.InboxWorkerConfig
}

func NewInbox(
	log *logium.Entry,
	db *pgdbx.DB,
	handlers handlers,
	config eventpg.InboxWorkerConfig,
) *Inbox {
	return &Inbox{
		log:      log.WithComponent("inbox"),
		db:       db,
		handlers: handlers,
		config:   config,
	}
}

func (b *Inbox) Start(ctx context.Context) {
	id := BuildProcessID("inbox")
	worker := eventpg.NewInboxWorker(id, b.log, b.db, b.config)

	defer func() {
		if err := worker.Stop(context.Background()); err != nil {
			b.log.WithError(err).Errorf("stop inbox worker %s failed", id)
		}
	}()

	worker.Route(evtypes.AccountCreatedEvent, b.handlers.AccountCreated)
	worker.Route(evtypes.AccountDeletedEvent, b.handlers.AccountDeleted)
	worker.Route(evtypes.AccountUsernameUpdatedEvent, b.handlers.AccountUsernameUpdated)

	worker.Run(ctx)
}

func BuildProcessID(service string) string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s-%s-%d", service, hostname, os.Getpid())
}
