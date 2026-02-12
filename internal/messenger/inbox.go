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
	log      *logium.Entry
	db       *pgdbx.DB
	handlers handlers
	config   eventpg.InboxWorkerConfig
}

func NewInbox(
	log *logium.Logger,
	db *pgdbx.DB,
	handlers handlers,
	config eventpg.InboxWorkerConfig,
) *Inbox {
	return &Inbox{
		log:      log.WithField("component", "inbox"),
		db:       db,
		handlers: handlers,
		config:   config,
	}
}

func (a *Inbox) Start(ctx context.Context) {
	a.log.Infoln("starting inbox worker")

	id := BuildProcessID("profiles-svc", "inbox", 0)
	worker := eventpg.NewInboxWorker(a.log, a.db, id, a.config)

	defer func() {
		if err := worker.Stop(context.Background()); err != nil {
			a.log.WithError(err).Errorf("stop inbox worker %s failed", id)
		}
	}()

	worker.Route(evtypes.AccountCreatedEvent, a.handlers.AccountCreated)
	worker.Route(evtypes.AccountDeletedEvent, a.handlers.AccountDeleted)
	worker.Route(evtypes.AccountUsernameUpdatedEvent, a.handlers.AccountUsernameUpdated)

	worker.Run(ctx)
}

func BuildProcessID(service string, role string, index int) string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return fmt.Sprintf("%s-%s-%d-%s", service, role, index, hostname)
}
