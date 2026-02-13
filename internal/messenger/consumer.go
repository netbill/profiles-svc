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

type Consumer struct {
	log          *logium.Entry
	db           *pgdbx.DB
	brokers      []string
	topicReaders map[string]int
}

func NewConsumer(
	log *logium.Entry,
	db *pgdbx.DB,
	brokers []string,
	topicReaders map[string]int,
) *Consumer {
	return &Consumer{
		log:          log.WithComponent("kafka-consumer"),
		db:           db,
		brokers:      brokers,
		topicReaders: topicReaders,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	var wg sync.WaitGroup

	accountReadersNum, ok := c.topicReaders[evtypes.AccountsTopicV1]
	if !ok || accountReadersNum <= 0 {
		c.log.Fatalf("number of readers for topic %s must be greater than 0", evtypes.AccountsTopicV1)
	}

	accountConsumer := eventpg.NewConsumer(c.log, c.db, eventpg.ConsumerConfig{
		MinBackoff: 200 * time.Millisecond,
		MaxBackoff: 5 * time.Second,
	})

	c.log.Infof("starting %d readers for topic %s", accountReadersNum, evtypes.AccountsTopicV1)

	wg.Add(accountReadersNum)

	for i := 0; i < accountReadersNum; i++ {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: c.brokers,
			GroupID: evtypes.ProfilesSvcGroup,
			Topic:   evtypes.AccountsTopicV1,
		})
		go func(r *kafka.Reader) {
			defer r.Close()
			defer wg.Done()

			accountConsumer.Read(ctx, r) // Read сам закроет reader
		}(reader)
	}

	wg.Wait()
}
