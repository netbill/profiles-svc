package messenger

import (
	"time"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/profiles-svc/pkg/log"
)

type ConsumerConfig struct {
	GroupID string   `json:"group_id"`
	Brokers []string `json:"brokers"`

	MinBackoff time.Duration `json:"min_backoff"`
	MaxBackoff time.Duration `json:"max_backoff"`

	AccountsV1 ConsumeKafkaConfig `json:"accounts_v1"`
}

type ConsumeKafkaConfig struct {
	Instances     int           `json:"instances"`
	MinBytes      int           `json:"min_bytes"`
	MaxBytes      int           `json:"max_bytes"`
	MaxWait       time.Duration `json:"max_wait"`
	QueueCapacity int           `json:"queue_capacity"`
}

func NewConsumer(
	logger *log.Logger,
	inbox eventbox.Inbox,
	config ConsumerConfig,
) *eventbox.Consumer {
	consumer := eventbox.NewConsumer(logger, inbox, eventbox.ConsumerConfig{
		MinBackoff: config.MinBackoff,
		MaxBackoff: config.MaxBackoff,
	})

	consumer.AddReader(eventbox.ReaderConfig{
		Brokers:       config.Brokers,
		GroupID:       config.GroupID,
		Topic:         evtypes.AccountsTopicV1,
		Instances:     config.AccountsV1.Instances,
		MaxWait:       config.AccountsV1.MaxWait,
		MinBytes:      config.AccountsV1.MinBytes,
		MaxBytes:      config.AccountsV1.MaxBytes,
		QueueCapacity: config.AccountsV1.QueueCapacity,
	})

	return consumer
}
