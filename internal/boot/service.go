package boot

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/bucket"
	"github.com/netbill/profiles-svc/internal/core/modules/account"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/messenger"
	"github.com/netbill/profiles-svc/internal/messenger/inbound"
	"github.com/netbill/profiles-svc/internal/messenger/outbound"
	"github.com/netbill/profiles-svc/internal/repository"
	"github.com/netbill/profiles-svc/internal/repository/pg"
	"github.com/netbill/profiles-svc/internal/rest"
	"github.com/netbill/profiles-svc/internal/rest/controller"
	"github.com/netbill/profiles-svc/internal/rest/middlewares"
	"github.com/netbill/profiles-svc/internal/tokenmanager"
	"github.com/netbill/restkit"
)

func RunService(ctx context.Context, log *logium.Entry, wg *sync.WaitGroup, cfg *Config) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			defer wg.Done()
		}()
	}

	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	db := pgdbx.NewDB(pool)

	s3 := newAws(cfg.S3.AWS)

	s3Bucket := bucket.New(s3, cfg.S3.Media)

	profilesSqlQ := pg.NewProfilesQ(db)
	transactionSqlQ := pg.NewTransaction(db)
	accountsSqlQ := pg.NewAccountsQ(db)
	repo := repository.New(transactionSqlQ, accountsSqlQ, profilesSqlQ)

	msg := messenger.NewManager(log, db, cfg.Kafka)

	kafkaProducer := msg.NewProducer()
	kafkaOutbound := outbound.New(kafkaProducer)

	tokenManager := tokenmanager.New(ServiceName, cfg.Auth.Tokens)

	profileModule := profile.New(repo, kafkaOutbound, tokenManager, s3Bucket)
	accountModule := account.New(repo, kafkaOutbound)

	responser := restkit.NewResponser()
	ctrl := controller.New(controller.Modules{
		Profile: profileModule,
	}, responser)

	mdll := middlewares.New(responser, tokenManager)

	router := rest.New(mdll, ctrl)

	run(func() {
		router.Run(ctx, log, cfg.Rest)
	})

	run(func() { msg.RunInbox(ctx, inbound.New(accountModule)) })

	run(func() { msg.RunConsumer(ctx) })

	run(func() { msg.RunOutbox(ctx) })
}
