package bucket

import (
	"context"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"
)

type Bucket struct {
	s3                     storage
	profileAvatarValidator ObjectValidator
	tokensTTL              UploadTokensTTL
}

type UploadTokensTTL struct {
	ProfileAvatar time.Duration
}

type Config struct {
	Storage                storage
	ProfileAvatarValidator ObjectValidator
	UploadTokensTTL        UploadTokensTTL
}

func New(config Config) Bucket {
	return Bucket{
		s3:                     config.Storage,
		tokensTTL:              config.UploadTokensTTL,
		profileAvatarValidator: config.ProfileAvatarValidator,
	}
}

type storage interface {
	PresignPut(
		ctx context.Context,
		key string,
		ttl time.Duration,
	) (uploadURL, getUrl string, error error)

	GetObjectRange(
		ctx context.Context,
		key string,
		bytes int64,
	) (body io.ReadCloser, size int64, err error)
	CopyObject(ctx context.Context, tmplKey, finalKey string) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

type ObjectValidator interface {
	ValidateImage(data []byte) (bool, error)
	ValidateImageFormat(data []byte) (bool, error)
	ValidateImageContentType(data []byte) (bool, error)
	ValidateImageSize(size uint) (bool, error)
}
