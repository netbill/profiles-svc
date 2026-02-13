package messenger

import (
	"context"

	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
	"github.com/segmentio/kafka-go"
)

type Outbox struct {
	log *logium.Entry
	db  *pgdbx.DB

	brokers []string
	config  eventpg.OutboxWorkerConfig
}

func NewOutbox(
	log *logium.Entry,
	db *pgdbx.DB,
	brokers []string,
	config eventpg.OutboxWorkerConfig,
) *Outbox {
	return &Outbox{
		db:      db,
		log:     log.WithComponent("outbox"),
		brokers: brokers,
		config:  config,
	}
}

func (b *Outbox) Run(ctx context.Context) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(b.brokers...),
		RequiredAcks: kafka.RequireAll,
		Compression:  kafka.Snappy,
		Balancer:     &kafka.LeastBytes{},
	}

	defer func() {
		if err := writer.Close(); err != nil {
			b.log.WithError(err).Error("failed to close kafka writer")
		}
	}()

	id := BuildProcessID("outbox")
	worker := eventpg.NewOutboxWorker(id, b.log, b.db, writer, b.config)

	defer func() {
		if err := worker.Stop(context.Background()); err != nil {
			b.log.WithError(err).WithField("worker_id", id).Error("failed to stop outbox worker")
		}
	}()

	worker.Run(ctx)
}
