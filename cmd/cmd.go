package cmd

import (
	"context"
	"database/sql"
	"sync"

	"github.com/netbill/evebox/box/inbox"
	"github.com/netbill/evebox/box/outbox"
	"github.com/netbill/logium"
	"github.com/netbill/profiles-svc/internal"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/messenger"
	"github.com/netbill/profiles-svc/internal/messenger/inbound"
	"github.com/netbill/profiles-svc/internal/messenger/outbound"
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

	outBox := outbox.New(pg)
	inBox := inbox.New(pg)

	kafkaOutbound := outbound.New(log, outBox)

	profileSvc := profile.New(repo, kafkaOutbound)

	kafkaInbound := inbound.New(log, profileSvc)

	ctrl := controller.New(log, profileSvc)
	mdll := mdlv.New(cfg.JWT.User.AccessToken.SecretKey, rest.AccountDataCtxKey, log)
	router := rest.New(log, mdll, ctrl)

	kafkaConsumer := messenger.NewConsumer(log, inBox, kafkaInbound, cfg.Kafka.Brokers...)

	kafkaProducer := messenger.NewProducer(log, outBox, cfg.Kafka.Brokers...)

	run(func() { router.Run(ctx, cfg) })

	log.Infof("starting kafka brokers %s", cfg.Kafka.Brokers)

	run(func() { kafkaConsumer.Run(ctx) })

	run(func() { kafkaProducer.Run(ctx) })
}
