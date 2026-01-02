package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/netbill/kafkakit/box"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/messenger/consumer"
	"github.com/netbill/profiles-svc/internal/messenger/consumer/callback"
	"github.com/netbill/profiles-svc/internal/messenger/producer"
	"github.com/netbill/profiles-svc/internal/repository"
	"github.com/netbill/restkit/mdlv"

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
	mdll := mdlv.New(cfg.JWT.User.AccessToken.SecretKey, rest.AccountDataCtxKey)
	router := rest.New(log, mdll, ctrl)

	kafkaConsumer := consumer.New(log, cfg.Kafka.Brokers, box.New(pg), callbacker.New(log, profileSvc))

	run(func() { router.Run(ctx, cfg) })

	run(func() { kafkaConsumer.Run(ctx) })

	run(func() { kafkaProducer.Run(ctx) })
}
