package messenger

import (
	"github.com/netbill/logium"
	"github.com/netbill/msnger"
	pgm "github.com/netbill/msnger/pg"
	"github.com/segmentio/kafka-go"
)

func Producer(log *logium.Logger, addr ...string) msnger.Producer {
	prod := pgm.NewProducer(
		log,
		&kafka.Writer{
			Addr:         kafka.TCP(addr...),
			RequiredAcks: kafka.RequireAll,
			Compression:  kafka.Snappy,
			Balancer:     &kafka.LeastBytes{},
		},
		nil,
	)

	return prod
}
