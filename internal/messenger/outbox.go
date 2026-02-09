package messenger

import (
	"context"
	"sync"

	"github.com/netbill/logium"
	pgm "github.com/netbill/msnger/pg"
	"github.com/segmentio/kafka-go"
)

type OutboxArchitectConfig struct {
	KafkaAddr  []string
	Processors []string

	outbox pgm.OutboxProcessorConfig
}

func StartOutboxArchitect(
	ctx context.Context,
	wg *sync.WaitGroup,
	log *logium.Logger,
	cfg OutboxArchitectConfig,
) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		writer := &kafka.Writer{
			Addr:         kafka.TCP(cfg.KafkaAddr...),
			RequiredAcks: kafka.RequireAll,
			Compression:  kafka.Snappy,
			Balancer:     &kafka.LeastBytes{},
		}
		defer func() {
			if err := writer.Close(); err != nil {
				log.Error("failed to close kafka writer", "error", err)
			}
		}()

		processor := pgm.NewOutboxProcessor(log, cfg.outbox, nil, writer)

		wg.Add(len(cfg.Processors))

		for _, p := range cfg.Processors {
			go func(processID string) {
				defer wg.Done()
				defer processor.StopProcess(context.Background(), processID)

				processor.StartProcess(ctx, processID)
			}(p)
		}

		<-ctx.Done()
	}()
}
