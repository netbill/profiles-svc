package messenger

import (
	"context"
	"sync"
	"time"

	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

type ConsumerArchitect struct {
	log          *logium.Entry
	db           *pgdbx.DB
	brokers      []string
	topicReaders map[string]int
}

func NewConsumerArchitect(
	log *logium.Entry,
	db *pgdbx.DB,
	brokers []string,
	topicReaders map[string]int,
) *ConsumerArchitect {
	return &ConsumerArchitect{
		log:          log.WithField("component", "kafka-consumer"),
		db:           db,
		brokers:      brokers,
		topicReaders: topicReaders,
	}
}

func (a *ConsumerArchitect) Start(ctx context.Context) {
	var wg sync.WaitGroup

	accountReadersNum, ok := a.topicReaders[evtypes.AccountsTopicV1]
	if !ok || accountReadersNum <= 0 {
		a.log.Fatalf("number of readers for topic %s must be greater than 0", evtypes.AccountsTopicV1)
	}

	accountConsumer := eventpg.NewConsumer(a.log, a.db, eventpg.ConsumerConfig{
		MinBackoff: 200 * time.Millisecond,
		MaxBackoff: 5 * time.Second,
	})

	a.log.Infoln("starting kafka consumers process")

	wg.Add(accountReadersNum)

	for i := 0; i < accountReadersNum; i++ {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  a.brokers,
			GroupID:  evtypes.ProfilesSvcGroup,
			Topic:    evtypes.AccountsTopicV1,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		})

		go func(r *kafka.Reader) {
			defer wg.Done()
			accountConsumer.Read(ctx, r) // Read сам закроет reader
		}(reader)
	}

	wg.Wait()
}
