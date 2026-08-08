package app

import (
	"context"
	"fmt"
	"sync"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/netbill/awsx"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/api/rest"
	"github.com/netbill/profiles-svc/internal/api/rest/controller"
	"github.com/netbill/profiles-svc/internal/api/rest/middlewares"
	"github.com/netbill/profiles-svc/internal/media"
	"github.com/netbill/profiles-svc/internal/messenger"
	"github.com/netbill/profiles-svc/internal/messenger/handler"
	"github.com/netbill/profiles-svc/internal/modules/account"
	"github.com/netbill/profiles-svc/internal/modules/profile"
	"github.com/netbill/profiles-svc/internal/repo/pg"
	"github.com/netbill/profiles-svc/pkg/tokenmanager"
	"github.com/netbill/profiles-svc/pkg/username"
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

	s3 := media.NewStorage(awsx.New(a.config.S3.Aws.BucketName, cfg), media.Config{
		LinkTTL:       a.config.S3.Media.Link.TTL,
		ProfileAvatar: a.config.S3.Media.Resources.Profile.Avatar,
	})

	tokenManager := tokenmanager.New(tokenmanager.Config{
		Issuer:   a.config.Auth.Tokens.Issuer,
		AccessSK: a.config.Auth.Tokens.AccountAccess.SecretKey,
	})

	profileRepo := pg.NewProfileRepo(db)
	accountRepo := pg.NewAccountRepo(db)
	outbox := pg.NewOutboxRepo(db, a.config.Kafka.Identity)
	usernameValidator := username.NewValidator()

	profileModule := profile.NewProfileModule(profile.ServiceDeps{
		Repo:      profileRepo,
		Bucket:    s3,
		Messenger: outbox,
		Tx:        db,
		Username:  usernameValidator,
	})

	accountModule := account.NewAccountModule(account.ServiceDeps{
		AccountRepo: accountRepo,
		ProfileRepo: profileRepo,
		Transaction: db,
		Messenger:   outbox,
		Username:    usernameValidator,
	})

	ctrl := controller.New(profileModule)
	mdll := middlewares.New(tokenManager)
	router := rest.NewServer(rest.ServerDeps{
		Middlewares: mdll,
		Profile:     ctrl,
		Log:         a.log,
		Resolver:    media.NewResolver(a.config.S3.Aws.BaseURL),
	})

	run(func() {
		router.Run(ctx, rest.Config{
			Port:              a.config.Rest.Port,
			ReadTimeout:       a.config.Rest.Timeouts.Read,
			ReadHeaderTimeout: a.config.Rest.Timeouts.ReadHeader,
			WriteTimeout:      a.config.Rest.Timeouts.Write,
			IdleTimeout:       a.config.Rest.Timeouts.Idle,
		})
	})

	inbound := handler.New(a.log, handler.Modules{
		Account: accountModule,
	})

	consumer := messenger.NewConsumer(inbound, messenger.ConsumerConfig{
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
