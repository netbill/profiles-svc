package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
)

const ProfileAvatarUploadTTL time.Duration = 1 * time.Hour
const ProfileAvatarMaxLength int64 = 5 * 1024 * 1024 // 5 MB

func CreateTempProfileAvatarKey(accountID, sessionID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/temp/%s", accountID, sessionID)
}

func CreateProfileAvatarKey(accountID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s", accountID)
}

var allowedProfileAvatarContentTypes = []string{
	"image/png",
	"image/jpeg",
	"image/jpg",
	"image/img",
	"image/gif",
}

func getAllowedProfileAvatarContentTypes() []string {
	return allowedProfileAvatarContentTypes
}

type Bucket struct {
	awsx3 awsx3
}

func New(awsx3 awsx3) Bucket {
	return Bucket{
		awsx3: awsx3,
	}
}

type awsx3 interface {
	PresignPut(
		ctx context.Context,
		key string,
		contentLength int64,
		ttl time.Duration,
	) (uploadURL, getUrl string, error error)

	HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error)
	CopyObject(ctx context.Context, tmplKey, finalKey string) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

func (r Bucket) GetPreloadLinkForUpdateProfileAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) (uploadURL, getUrl string, error error) {
	uploadURL, getURL, err := r.awsx3.PresignPut(
		ctx,
		CreateTempProfileAvatarKey(accountID, sessionID),
		ProfileAvatarMaxLength,
		ProfileAvatarUploadTTL,
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to presign put object for profile avatar: %w", err,
		)
	}

	return uploadURL, getURL, nil
}

func (r Bucket) AcceptUpdateProfileAvatar(ctx context.Context, accountID, sessionID uuid.UUID) (string, error) {
	obj, err := r.awsx3.HeadObject(ctx, CreateTempProfileAvatarKey(accountID, sessionID))
	if err != nil {
		return "", fmt.Errorf(
			"failed to head object for profile avatar: %w", err,
		)
	}

	ct := *obj.ContentType

	allowed := false
	for _, act := range getAllowedProfileAvatarContentTypes() {
		if ct == act {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", errx.ErrorContentTypeIsNotAllowed.Raise(
			fmt.Errorf(
				"profile avatar extension %s not allowed, allowed only: %s",
				ct, getAllowedProfileAvatarContentTypes(),
			),
		)
	}

	res, err := r.awsx3.CopyObject(ctx,
		CreateTempProfileAvatarKey(accountID, sessionID),
		CreateProfileAvatarKey(accountID),
	)
	if err != nil {
		return "", fmt.Errorf(
			"failed to copy object for profile avatar: %w", err,
		)
	}

	return res, nil
}

func (r Bucket) CancelUpdateProfileAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) error {
	err := r.awsx3.DeleteObject(ctx, CreateTempProfileAvatarKey(accountID, sessionID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete temp object for profile avatar: %w", err,
		)
	}

	return nil
}
