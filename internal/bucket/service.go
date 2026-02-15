package bucket

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/netbill/awsx"
)

type Config struct {
	Link struct {
		TTL time.Duration `json:"ttl"`
	} `json:"link"`
	Profile struct {
		Avatar struct {
			AllowedFormats   []string `mapstructure:"allowed_formats" required:"true"`
			MaxWidth         int      `mapstructure:"max_width" required:"true"`
			MaxHeight        int      `mapstructure:"max_height" required:"true"`
			ContentLengthMax int64    `mapstructure:"content_length_max" required:"true"`
		} `mapstructure:"avatar"`
	} `mapstructure:"profile"`
}

type Bucket struct {
	s3     awsx.Bucket
	config Config
}

func New(s3 awsx.Bucket, config Config) Bucket {
	return Bucket{
		s3:     s3,
		config: config,
	}
}
