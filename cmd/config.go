package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Name string `mapstructure:"name"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type RestConfig struct {
	Port     string `mapstructure:"port"`
	Timeouts struct {
		Read       time.Duration `mapstructure:"read"`
		ReadHeader time.Duration `mapstructure:"read_header"`
		Write      time.Duration `mapstructure:"write"`
		Idle       time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeouts"`
}

type DatabaseConfig struct {
	SQL struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"sql"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
}

type AuthConfig struct {
	Account struct {
		Token struct {
			Access struct {
				SecretKey string `mapstructure:"secret_key"`
			} `mapstructure:"access"`
		} `mapstructure:"token"`
	} `mapstructure:"account"`
}

type S3Config struct {
	AWS struct {
		BucketName      string `mapstructure:"bucket_name"`
		Region          string `mapstructure:"region"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	} `mapstructure:"aws"`

	Upload struct {
		Token struct {
			SecretKey string `mapstructure:"secret_key"`
			TTL       struct {
				Profile time.Duration `mapstructure:"profile_avatar"`
			} `mapstructure:"ttl"`
		} `mapstructure:"token"`

		Profile struct {
			Avatar struct {
				AllowedContentTypes []string `mapstructure:"allowed_content_types"`
				AllowedFormats      []string `mapstructure:"allowed_formats"`
				MaxWidth            uint     `mapstructure:"max_width"`
				MaxHeight           uint     `mapstructure:"max_height"`
				ContentLengthMax    uint     `mapstructure:"content_length_max"`
			} `mapstructure:"avatar"`
		} `mapstructure:"profile"`
	} `mapstructure:"upload"`
}

type Config struct {
	Service  ServiceConfig  `mapstructure:"service"`
	Log      LogConfig      `mapstructure:"log"`
	Rest     RestConfig     `mapstructure:"rest"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Database DatabaseConfig `mapstructure:"database"`
	S3       S3Config       `mapstructure:"s3"`
}

func LoadConfig() (Config, error) {
	configPath := os.Getenv("KV_VIPER_FILE")
	if configPath == "" {
		return Config{}, fmt.Errorf("KV_VIPER_FILE env var is not set")
	}

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("error reading config file: %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config: %s", err)
	}

	return config, nil
}
