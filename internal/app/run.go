package app

import (
	"context"
	"fmt"
	"sync"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/netbill/awsx"
	"github.com/netbill/eventbox"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/bucket"
	"github.com/netbill/profiles-svc/internal/core/modules/account"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/messenger"
	"github.com/netbill/profiles-svc/internal/messenger/handler"
	"github.com/netbill/profiles-svc/internal/messenger/publisher"
	"github.com/netbill/profiles-svc/internal/repository"
	"github.com/netbill/profiles-svc/internal/repository/pg"
	"github.com/netbill/profiles-svc/internal/rest"
	"github.com/netbill/profiles-svc/internal/rest/controller"
	"github.com/netbill/profiles-svc/internal/rest/middlewares"
	"github.com/netbill/profiles-svc/internal/tokenmanager"
)

func (a *App) Run(ctx context.Context) error {
	var wg = &sync.WaitGroup{}

	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	pool, err := a.config.PoolDB(ctx)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	db := pgdbx.NewDB(pool)

	a.log.Info("starting application")

	repo := &repository.Repository{
		ProfilesQ:     pg.NewProfilesQ(db),
		AccountsQ:     pg.NewAccountsQ(db),
		Transactioner: pg.NewTransaction(db),
	}

	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithRegion(a.config.S3.Aws.Region),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				a.config.S3.Aws.AccessKeyID,
				a.config.S3.Aws.SecretAccessKey,
				a.config.S3.Aws.SessionToken,
			),
		),
	)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}

	s3 := bucket.NewStorage(awsx.New(a.config.S3.Aws.BucketName, cfg), bucket.Config{
		LinkTTL:       a.config.S3.Media.Link.TTL,
		ProfileAvatar: a.config.S3.Media.Resources.Profile.Avatar,
	})

	outbox := eventpg.NewOutbox(db)
	inbox := eventpg.NewInbox(db)

	producer := messenger.NewProducer(a.log, messenger.ProducerConfig{
		Producer: a.config.Kafka.Identity,
		Brokers:  a.config.Kafka.Brokers,
		ProfilesV1: messenger.ProduceKafkaConfig{
			RequiredAcks: a.config.Kafka.Produce.Topics.ProfilesV1.RequiredAcks,
			Compression:  a.config.Kafka.Produce.Topics.ProfilesV1.Compression,
			Balancer:     a.config.Kafka.Produce.Topics.ProfilesV1.Balancer,
			BatchSize:    a.config.Kafka.Produce.Topics.ProfilesV1.BatchSize,
			BatchTimeout: a.config.Kafka.Produce.Topics.ProfilesV1.BatchTimeout,
		},
	})
	defer producer.Close()

	outbound := publisher.New(a.config.Kafka.Identity, outbox, producer)

	tokenManager := tokenmanager.New(tokenmanager.Config{
		Issuer:   a.config.Auth.Tokens.Issuer,
		AccessSK: a.config.Auth.Tokens.AccountAccess.SecretKey,
	})

	profileModule := profile.New(repo, outbound, s3)
	accountModule := account.New(repo, outbound)

	ctrl := controller.New(controller.Modules{
		Profile: profileModule,
	})
	mdll := middlewares.New(tokenManager)
	router := rest.New(mdll, ctrl)

	run(func() {
		router.Run(ctx, a.log, rest.Config{
			Port:              a.config.Rest.Port,
			ReadTimeout:       a.config.Rest.Timeouts.Read,
			ReadHeaderTimeout: a.config.Rest.Timeouts.ReadHeader,
			WriteTimeout:      a.config.Rest.Timeouts.Write,
			IdleTimeout:       a.config.Rest.Timeouts.Idle,
		})
	})

	outboxWorker := messenger.NewOutboxWorker(a.log, outbox, producer, eventbox.OutboxWorkerConfig{
		Routines:       a.config.Kafka.Outbox.Routines,
		Slots:          a.config.Kafka.Outbox.Slots,
		BatchSize:      a.config.Kafka.Outbox.BatchSize,
		Sleep:          a.config.Kafka.Outbox.Sleep,
		MinNextAttempt: a.config.Kafka.Outbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Outbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Outbox.MaxAttempts,
	})
	defer outboxWorker.Clean()

	run(func() {
		outboxWorker.Run(ctx)
	})

	inbound := handler.New(accountModule)

	inboxWorker := messenger.NewInboxWorker(a.log, inbox, eventbox.InboxWorkerConfig{
		Routines:       a.config.Kafka.Inbox.Routines,
		Slots:          a.config.Kafka.Inbox.Slots,
		BatchSize:      a.config.Kafka.Inbox.BatchSize,
		Sleep:          a.config.Kafka.Inbox.Sleep,
		MinNextAttempt: a.config.Kafka.Inbox.MinNextAttempt,
		MaxNextAttempt: a.config.Kafka.Inbox.MaxNextAttempt,
		MaxAttempts:    a.config.Kafka.Inbox.MaxAttempts,
	}, inbound)
	defer inboxWorker.Clean()

	run(func() {
		inboxWorker.Run(ctx)
	})

	consumer := messenger.NewConsumer(a.log, inbox, messenger.ConsumerConfig{
		GroupID:    a.config.Kafka.Identity,
		Brokers:    a.config.Kafka.Brokers,
		MinBackoff: a.config.Kafka.Consume.Backoff.Min,
		MaxBackoff: a.config.Kafka.Consume.Backoff.Max,
		AccountsV1: messenger.ConsumeKafkaConfig{
			Instances:     a.config.Kafka.Consume.Topics.AccountsV1.Instances,
			MinBytes:      a.config.Kafka.Consume.Topics.AccountsV1.MinBytes,
			MaxBytes:      a.config.Kafka.Consume.Topics.AccountsV1.MaxBytes,
			MaxWait:       a.config.Kafka.Consume.Topics.AccountsV1.MaxWait,
			QueueCapacity: a.config.Kafka.Consume.Topics.AccountsV1.QueueCapacity,
		},
	})
	defer consumer.Close()

	run(func() {
		consumer.Run(ctx)
	})

	wg.Wait()
	return nil
}
