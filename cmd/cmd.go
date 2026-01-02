package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/netbill/kafkakit/box"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal"
	"github.com/netbill/profiles-svc/internal/domain/modules/profile"
	"github.com/netbill/profiles-svc/internal/messanger/consumer"
	"github.com/netbill/profiles-svc/internal/messanger/consumer/callback"
	"github.com/netbill/profiles-svc/internal/messanger/producer"
	"github.com/netbill/profiles-svc/internal/repository"
	"github.com/netbill/profiles-svc/internal/rest/middlewares"

	"github.com/netbill/profiles-svc/internal/rest"
	"github.com/netbill/profiles-svc/internal/rest/controller"
)

func StartServices(ctx context.Context, cfg internal.Config, log logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	repo := repository.New(pg)
	kafkaBox := box.New(pg)

	kafkaProducer := producer.New(log, cfg.Kafka.Brokers, kafkaBox)

	profileSvc := profile.New(repo, kafkaProducer)

	ctrl := controller.New(log, profileSvc)
	mdlv := middlewares.New(log)

	kafkaConsumer := consumer.New(log, cfg.Kafka.Brokers, callback.NewService(log, kafkaBox))
	kafkaInboxWorker := consumer.NewInboxWorker(log, kafkaBox, profileSvc)

	run(func() { kafkaConsumer.Run(ctx) })

	run(func() { kafkaInboxWorker.Run(ctx) })

	run(func() { rest.Run(ctx, cfg, log, mdlv, ctrl) })
}
