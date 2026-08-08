package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/awsx"
	"github.com/netbill/profiles-svc/pkg/log"
)

type ServiceCfg struct {
	Name string
}

type LogConfig struct {
	Level  string
	Format string
}

type SQLConfig struct {
	URL string
}

type DatabaseConfig struct {
	SQL SQLConfig
}

type RestTimeouts struct {
	Read       time.Duration
	ReadHeader time.Duration
	Write      time.Duration
	Idle       time.Duration
}

type RestConfig struct {
	Port     int
	Timeouts RestTimeouts
}

type TokenConfig struct {
	SecretKey string
}

type AuthTokensConfig struct {
	Issuer        string
	AccountAccess TokenConfig
}

type AuthConfig struct {
	Tokens AuthTokensConfig
}

type S3AwsConfig struct {
	Region          string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	BaseURL         string
}

type S3MediaProfileConfig struct {
	Avatar awsx.ImageValidator
}

type S3MediaResourcesConfig struct {
	Profile S3MediaProfileConfig
}

type S3MediaLinkConfig struct {
	TTL time.Duration
}

type S3MediaConfig struct {
	Link      S3MediaLinkConfig
	Resources S3MediaResourcesConfig
}

type S3Config struct {
	Aws   S3AwsConfig
	Media S3MediaConfig
}

type KafkaBackoffConfig struct {
	Min time.Duration
	Max time.Duration
}

type KafkaConsumeAccountsV1Config struct {
	Instances     int
	MinBytes      int
	MaxBytes      int
	MaxWait       time.Duration
	QueueCapacity int
}

type KafkaConsumeTopicsConfig struct {
	AccountsV1 KafkaConsumeAccountsV1Config
}

type KafkaConsumeConfig struct {
	Backoff KafkaBackoffConfig
	Topics  KafkaConsumeTopicsConfig
}

type KafkaProduceProfilesV1Config struct {
	RequiredAcks string
	Compression  string
	Balancer     string
	BatchSize    int
	BatchTimeout time.Duration
}

type KafkaProduceTopicsConfig struct {
	ProfilesV1 KafkaProduceProfilesV1Config
}

type KafkaProduceConfig struct {
	Topics KafkaProduceTopicsConfig
}

type KafkaBoxConfig struct {
	Routines       int
	Slots          int
	BatchSize      int
	Sleep          time.Duration
	MinNextAttempt time.Duration
	MaxNextAttempt time.Duration
	MaxAttempts    int32
}

type KafkaConfig struct {
	Brokers  []string
	Identity string

	Consume KafkaConsumeConfig
	Produce KafkaProduceConfig

	Inbox  KafkaBoxConfig
	Outbox KafkaBoxConfig
}

type Config struct {
	Service  ServiceCfg
	Log      LogConfig
	Database DatabaseConfig
	Rest     RestConfig
	Auth     AuthConfig
	S3       S3Config
	Kafka    KafkaConfig
}

// LoadConfig builds the config entirely from environment variables — no
// config file involved. Vars backing an actual secret or a connection the
// service can't run without (DB, auth secret, S3 credentials, Kafka
// brokers) are required and panic if missing; everything else falls back
// to the same defaults the old config.example.yaml used to hardcode.
func LoadConfig() *Config {
	return &Config{
		Service: ServiceCfg{
			Name: envOr("SERVICE_NAME", "profiles-svc"),
		},
		Log: LogConfig{
			Level:  envOr("LOG_LEVEL", "debug"),
			Format: envOr("LOG_FORMAT", "text"),
		},
		Rest: RestConfig{
			Port: envIntOr("REST_PORT", 8002),
			Timeouts: RestTimeouts{
				Read:       envDurationOr("REST_READ_TIMEOUT", 15*time.Second),
				ReadHeader: envDurationOr("REST_READ_HEADER_TIMEOUT", 15*time.Second),
				Write:      envDurationOr("REST_WRITE_TIMEOUT", 15*time.Second),
				Idle:       envDurationOr("REST_IDLE_TIMEOUT", 60*time.Second),
			},
		},
		Auth: AuthConfig{
			Tokens: AuthTokensConfig{
				Issuer: envOr("AUTH_TOKENS_ISSUER", "profiles-svc"),
				AccountAccess: TokenConfig{
					SecretKey: mustEnv("AUTH_TOKENS_ACCOUNT_ACCESS_SECRET_KEY"),
				},
			},
		},
		Database: DatabaseConfig{
			SQL: SQLConfig{
				URL: mustEnv("DATABASE_SQL_URL"),
			},
		},
		S3: S3Config{
			Aws: S3AwsConfig{
				Region:          mustEnv("S3_AWS_REGION"),
				BucketName:      mustEnv("S3_AWS_BUCKET_NAME"),
				AccessKeyID:     mustEnv("S3_AWS_ACCESS_KEY_ID"),
				SecretAccessKey: mustEnv("S3_AWS_SECRET_ACCESS_KEY"),
				SessionToken:    envOr("S3_AWS_SESSION_TOKEN", ""),
				BaseURL:         mustEnv("S3_AWS_BASE_URL"),
			},
			Media: S3MediaConfig{
				Link: S3MediaLinkConfig{
					TTL: envDurationOr("S3_MEDIA_LINK_TTL", 24*time.Hour),
				},
				Resources: S3MediaResourcesConfig{
					Profile: S3MediaProfileConfig{
						Avatar: awsx.ImageValidator{
							AllowedFormats: envListOr(
								"S3_MEDIA_PROFILE_AVATAR_ALLOWED_FORMATS", []string{"jpeg", "jpg", "png", "gif"},
							),
							MaxWidth:       envIntOr("S3_MEDIA_PROFILE_AVATAR_MAX_WIDTH", 512),
							MinWidth:       envIntOr("S3_MEDIA_PROFILE_AVATAR_MIN_WIDTH", 512),
							MaxHeight:      envIntOr("S3_MEDIA_PROFILE_AVATAR_MAX_HEIGHT", 512),
							MinHeight:      envIntOr("S3_MEDIA_PROFILE_AVATAR_MIN_HEIGHT", 512),
							ContentSizeMax: envInt64Or("S3_MEDIA_PROFILE_AVATAR_CONTENT_SIZE_MAX", 5242880),
						},
					},
				},
			},
		},
		Kafka: KafkaConfig{
			Brokers:  mustEnvList("KAFKA_BROKERS"),
			Identity: envOr("KAFKA_IDENTITY", "profiles-svc"),
			Consume: KafkaConsumeConfig{
				Backoff: KafkaBackoffConfig{
					Min: envDurationOr("KAFKA_CONSUME_BACKOFF_MIN", 100*time.Millisecond),
					Max: envDurationOr("KAFKA_CONSUME_BACKOFF_MAX", time.Second),
				},
				Topics: KafkaConsumeTopicsConfig{
					AccountsV1: KafkaConsumeAccountsV1Config{
						Instances:     envIntOr("KAFKA_CONSUME_ACCOUNTS_V1_INSTANCES", 4),
						MinBytes:      envIntOr("KAFKA_CONSUME_ACCOUNTS_V1_MIN_BYTES", 1),
						MaxBytes:      envIntOr("KAFKA_CONSUME_ACCOUNTS_V1_MAX_BYTES", 10e6),
						MaxWait:       envDurationOr("KAFKA_CONSUME_ACCOUNTS_V1_MAX_WAIT", time.Second),
						QueueCapacity: envIntOr("KAFKA_CONSUME_ACCOUNTS_V1_QUEUE_CAPACITY", 10),
					},
				},
			},
			Produce: KafkaProduceConfig{
				Topics: KafkaProduceTopicsConfig{
					ProfilesV1: KafkaProduceProfilesV1Config{
						RequiredAcks: envOr("KAFKA_PRODUCE_PROFILES_V1_REQUIRED_ACKS", "all"),
						Compression:  envOr("KAFKA_PRODUCE_PROFILES_V1_COMPRESSION", "snappy"),
						Balancer:     envOr("KAFKA_PRODUCE_PROFILES_V1_BALANCER", "hash"),
						BatchSize:    envIntOr("KAFKA_PRODUCE_PROFILES_V1_BATCH_SIZE", 100),
						BatchTimeout: envDurationOr("KAFKA_PRODUCE_PROFILES_V1_BATCH_TIMEOUT", 10*time.Millisecond),
					},
				},
			},
			Inbox: KafkaBoxConfig{
				Routines:       envIntOr("KAFKA_INBOX_ROUTINES", 10),
				Slots:          envIntOr("KAFKA_INBOX_SLOTS", 40),
				BatchSize:      envIntOr("KAFKA_INBOX_BATCH_SIZE", 100),
				Sleep:          envDurationOr("KAFKA_INBOX_SLEEP", 100*time.Millisecond),
				MinNextAttempt: envDurationOr("KAFKA_INBOX_MIN_NEXT_ATTEMPT", time.Minute),
				MaxNextAttempt: envDurationOr("KAFKA_INBOX_MAX_NEXT_ATTEMPT", 10*time.Minute),
				MaxAttempts:    envInt32Or("KAFKA_INBOX_MAX_ATTEMPTS", 0),
			},
			Outbox: KafkaBoxConfig{
				Routines:       envIntOr("KAFKA_OUTBOX_ROUTINES", 10),
				Slots:          envIntOr("KAFKA_OUTBOX_SLOTS", 40),
				BatchSize:      envIntOr("KAFKA_OUTBOX_BATCH_SIZE", 100),
				Sleep:          envDurationOr("KAFKA_OUTBOX_SLEEP", 100*time.Millisecond),
				MinNextAttempt: envDurationOr("KAFKA_OUTBOX_MIN_NEXT_ATTEMPT", time.Minute),
				MaxNextAttempt: envDurationOr("KAFKA_OUTBOX_MAX_NEXT_ATTEMPT", 10*time.Minute),
				MaxAttempts:    envInt32Or("KAFKA_OUTBOX_MAX_ATTEMPTS", 0),
			},
		},
	}
}

func mustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		panic(fmt.Errorf("%s env var is not set", key))
	}
	return v
}

func mustEnvList(key string) []string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		panic(fmt.Errorf("%s env var is not set", key))
	}
	return splitList(v)
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func envListOr(key string, def []string) []string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}
	return splitList(v)
}

func splitList(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func envIntOr(key string, def int) int {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Errorf("invalid int value for %s: %w", key, err))
	}
	return n
}

func envInt32Or(key string, def int32) int32 {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}

	n, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		panic(fmt.Errorf("invalid int value for %s: %w", key, err))
	}
	return int32(n)
}

func envInt64Or(key string, def int64) int64 {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}

	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid int value for %s: %w", key, err))
	}
	return n
}

func envDurationOr(key string, def time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		panic(fmt.Errorf("invalid duration value for %s: %w", key, err))
	}
	return d
}

func (cfg *Config) Logger() *log.Logger {
	return log.New(cfg.Log.Level, cfg.Log.Format, cfg.Service.Name)
}

func (cfg *Config) PoolDB(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return pool, nil
}
