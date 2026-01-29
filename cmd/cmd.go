package cmd

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/awsx"
	"github.com/netbill/logium"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/bucket"
	"github.com/netbill/profiles-svc/internal/core/modules/profile"
	"github.com/netbill/profiles-svc/internal/messenger"
	"github.com/netbill/profiles-svc/internal/messenger/inbound"
	"github.com/netbill/profiles-svc/internal/messenger/outbound"
	"github.com/netbill/profiles-svc/internal/repository"
	"github.com/netbill/profiles-svc/internal/repository/pg"
	"github.com/netbill/profiles-svc/internal/rest/middlewares"
	"github.com/netbill/profiles-svc/internal/tokenmanager"

	"github.com/netbill/profiles-svc/internal/rest"
	"github.com/netbill/profiles-svc/internal/rest/controller"
)

func StartServices(ctx context.Context, cfg Config, log *logium.Logger, wg *sync.WaitGroup) {
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

	awsxSvc := awsx.New(
		cfg.S3.AWS.BucketName,
		s3Client,
		presignClient,
	)

	s3Bucket := bucket.New(awsxSvc)

	profilesSqlQ := pg.NewProfilesQ(db)
	transactionSqlQ := pg.NewTransaction(db)
	repo := repository.New(transactionSqlQ, profilesSqlQ)

	kafkaOutbound := outbound.New(log, db)

	tokenManager := tokenmanager.New(cfg.Service.Name, cfg.S3.Upload.Token.SecretKey)

	profileSvc := profile.New(repo, kafkaOutbound, tokenManager, s3Bucket)

	ctrl := controller.New(log, profileSvc)
	mdll := middlewares.New(log, middlewares.Config{
		AccountAccessSK: cfg.Auth.Account.Token.Access.SecretKey,
		UploadFilesSK:   cfg.S3.Upload.Token.SecretKey,
	})
	router := rest.New(log, mdll, ctrl)

	msgx := messenger.New(log, db, cfg.Kafka.Brokers...)

	run(func() {
		router.Run(ctx, rest.Config{
			Port:              cfg.Rest.Port,
			TimeoutRead:       cfg.Rest.Timeouts.Read,
			TimeoutReadHeader: cfg.Rest.Timeouts.ReadHeader,
			TimeoutWrite:      cfg.Rest.Timeouts.Write,
			TimeoutIdle:       cfg.Rest.Timeouts.Idle,
		})
	})

	log.Infof("starting kafka brokers %s", cfg.Kafka.Brokers)

	run(func() { msgx.RunProducer(ctx) })

	run(func() { msgx.RunConsumer(ctx, inbound.New(log, profileSvc)) })
}
