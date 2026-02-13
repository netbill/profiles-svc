package messenger

import (
	"context"
	"sync"
	"time"

	"github.com/netbill/eventbox"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	log      *logium.Entry
	consumer eventbox.Consumer

	groupID string
	brokers []string
	topics  map[string]ConsumerTopicConfig
}

type ConsumerTopicConfig struct {
	NumReaders     int
	QueueCapacity  int
	MaxBytes       int
	MinBytes       int
	MaxWait        time.Duration
	CommitInterval time.Duration
}

func NewConsumer(
	log *logium.Entry,
	db *pgdbx.DB,
	brokers ...string,
) *Consumer {
	return &Consumer{
		log:      log.WithComponent("kafka-consumer"),
		consumer: eventpg.NewConsumer(log, db, eventpg.ConsumerConfig{}),
		groupID:  evtypes.ProfilesSvcGroup,
		brokers:  brokers,
		topics:   make(map[string]ConsumerTopicConfig),
	}
}

func (g *Consumer) AddTopic(topic string, config ConsumerTopicConfig) {
	if config.NumReaders <= 0 {
		g.log.Fatalf("number of readers for topic %s must be greater than 0", topic)
	}

	g.topics[topic] = config
}

func (g *Consumer) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for topic, config := range g.topics {
		g.log.Infof("starting %d readers for topic %s", config.NumReaders, topic)

		for i := 0; i < config.NumReaders; i++ {
			reader := kafka.NewReader(kafka.ReaderConfig{
				Brokers:        g.brokers,
				Topic:          topic,
				GroupID:        g.groupID,
				QueueCapacity:  config.QueueCapacity,
				MaxBytes:       config.MaxBytes,
				MinBytes:       config.MinBytes,
				MaxWait:        config.MaxWait,
				CommitInterval: config.CommitInterval,
			})

			wg.Add(1)
			go func(r *kafka.Reader) {
				defer r.Close()
				defer wg.Done()

				g.consumer.Read(ctx, r)
			}(reader)
		}
	}

	wg.Wait()
}
