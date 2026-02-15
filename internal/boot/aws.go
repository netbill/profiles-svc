package boot

import (
	"context"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/netbill/awsx"
)

type AwsConfig struct {
	BucketName      string `mapstructure:"bucket_name"`
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	SessionToken    string `mapstructure:"session_token"`
}

func newAws(c AwsConfig) awsx.Bucket {
	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithRegion(c.Region),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				c.AccessKeyID,
				c.SecretAccessKey,
				c.SessionToken,
			),
		),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	return awsx.New(c.BucketName, s3Client, s3.NewPresignClient(s3Client))
}
