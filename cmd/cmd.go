package cmd

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/awsx"
	eventpg "github.com/netbill/eventbox/pg"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/cmd/config"
	"github.com/netbill/profiles-svc/internal/bucket"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/messenger"
	"github.com/netbill/profiles-svc/internal/messenger/contracts"
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

func StartServices(ctx context.Context, cfg config.Config, log *logium.Logger, wg *sync.WaitGroup) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	db := pgdbx.NewDB(pool)

	awsCfg := aws.Config{
		Region: cfg.S3.AWS.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.S3.AWS.AccessKeyID,
			cfg.S3.AWS.SecretAccessKey,
			"",
		),
	}

	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)

	awsS3 := awsx.New(
		cfg.S3.AWS.BucketName,
		s3Client,
		presignClient,
	)

	s3Bucket := bucket.New(awsS3, bucket.Config{
		Profile: bucket.ProfileConfig{
			TokenTTL:  cfg.S3.Upload.Token.TTL.Profile,
			MaxSize:   cfg.S3.Upload.Profile.Avatar.ContentLengthMax,
			MaxWidth:  cfg.S3.Upload.Profile.Avatar.MaxWidth,
			MaxHeight: cfg.S3.Upload.Profile.Avatar.MaxHeight,
			Formats:   cfg.S3.Upload.Profile.Avatar.AllowedFormats,
		},
	})

	profilesSqlQ := pg.NewProfilesQ(db)
	transactionSqlQ := pg.NewTransaction(db)
	repo := repository.New(transactionSqlQ, profilesSqlQ)

	kafkaProducer := messenger.NewProducer(db, cfg.Kafka.Brokers...)
	kafkaOutbound := outbound.New(kafkaProducer)

	tokenManager := tokenmanager.New(
		cfg.S3.Upload.Token.SecretKey,
		cfg.S3.Upload.Token.TTL.Profile,
	)

	profileModule := profile.New(repo, kafkaOutbound, tokenManager, s3Bucket)

	responser := restkit.NewResponser()
	ctrl := controller.New(log, responser, controller.Modules{
		Profile: profileModule,
	})
	mdll := middlewares.New(log, responser, middlewares.Config{
		AccountAccessSK: cfg.Auth.Account.Token.Access.SecretKey,
		UploadFilesSK:   cfg.S3.Upload.Token.SecretKey,
	})
	router := rest.New(log, mdll, ctrl)

	run(func() {
		router.Run(ctx, rest.Config{
			Port:              cfg.Rest.Port,
			TimeoutRead:       cfg.Rest.Timeouts.Read,
			TimeoutReadHeader: cfg.Rest.Timeouts.ReadHeader,
			TimeoutWrite:      cfg.Rest.Timeouts.Write,
			TimeoutIdle:       cfg.Rest.Timeouts.Idle,
		})
	})

	kafkaConsumer := messenger.NewConsumerArchitect(log, db, cfg.Kafka.Brokers, map[string]int{
		contracts.AccountsTopicV1: cfg.Kafka.Readers.AccountsV1,
	})

	run(func() { kafkaConsumer.Start(ctx) })

	kafkaInboxArh := messenger.NewInbox(log, db, inbound.New(profileModule), eventpg.InboxWorkerConfig{
		Routines:       cfg.Kafka.Inbox.Routines,
		MinSleep:       cfg.Kafka.Inbox.MinSleep,
		MaxSleep:       cfg.Kafka.Inbox.MaxSleep,
		MinBatch:       cfg.Kafka.Inbox.MinBatch,
		MaxBatch:       cfg.Kafka.Inbox.MaxBatch,
		MinNextAttempt: cfg.Kafka.Inbox.MinNextAttempt,
		MaxNextAttempt: cfg.Kafka.Inbox.MaxNextAttempt,
		MaxAttempts:    cfg.Kafka.Inbox.MaxAttempts,
	})

	run(func() { kafkaInboxArh.Start(ctx) })

	kafkaOutboxArch := messenger.NewOutbox(log, db, cfg.Kafka.Brokers, eventpg.OutboxWorkerConfig{
		Routines:       cfg.Kafka.Outbox.Routines,
		MinSleep:       cfg.Kafka.Outbox.MinSleep,
		MaxSleep:       cfg.Kafka.Outbox.MaxSleep,
		MinBatch:       cfg.Kafka.Outbox.MinBatch,
		MaxBatch:       cfg.Kafka.Outbox.MaxBatch,
		MinNextAttempt: cfg.Kafka.Outbox.MinNextAttempt,
		MaxNextAttempt: cfg.Kafka.Outbox.MaxNextAttempt,
		MaxAttempts:    cfg.Kafka.Outbox.MaxAttempts,
	})

	run(func() { kafkaOutboxArch.Start(ctx) })
}
