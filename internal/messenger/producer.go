package messenger

import (
	"time"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/profiles-svc/pkg/log"
)

type ProducerConfig struct {
	Producer string   `json:"producer"`
	Brokers  []string `json:"brokers"`

	ProfilesV1 ProduceKafkaConfig `json:"profiles_v1"`
}

type ProduceKafkaConfig struct {
	RequiredAcks string        `json:"required_acks"`
	Compression  string        `json:"compression"`
	Balancer     string        `json:"balancer"`
	BatchSize    int           `json:"batch_size"`
	BatchTimeout time.Duration `json:"batch_timeout"`
}

func NewProducer(log *log.Logger, cfg ProducerConfig) *eventbox.Producer {
	producer := eventbox.NewProducer(log, cfg.Brokers...)

	err := producer.AddWriter(evtypes.ProfilesTopicV1, eventbox.WriterTopicConfig{
		RequiredAcks: cfg.ProfilesV1.RequiredAcks,
		Compression:  cfg.ProfilesV1.Compression,
		Balancer:     cfg.ProfilesV1.Balancer,
		BatchSize:    cfg.ProfilesV1.BatchSize,
		BatchTimeout: cfg.ProfilesV1.BatchTimeout,
	})
	if err != nil {
		panic(err)
	}

	return producer
}
