package bucket

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/netbill/awsx"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func CreateTempProfileAvatarKey(accountID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/temp/%s", accountID, uuid.New())
}

func CreateProfileAvatarKey(accountID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s", accountID)
}

func (b Bucket) GetPreloadLinkForProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempProfileAvatarKey(accountID)

	uploadLink, getLink, err := b.s3.PresignPut(ctx, key, b.config.Link.TTL)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presign put object for profile avatar: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getLink,
		UploadURL:  uploadLink,
	}, nil
}

func (b Bucket) UpdateProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
	key string,
) (string, error) {
	if key == "" {
		return "", nil
	}

	if err := validateTempProfileAvatarKey(accountID, key); err != nil {
		return "", err
	}

	tempKey := key
	finalKey := CreateProfileAvatarKey(accountID)

	head, err := b.s3.HeadObject(ctx, tempKey)
	if err != nil {
		if errors.Is(err, awsx.ErrNotFound) {
			return "", errx.ErrorNoContentUploaded.Raise(
				fmt.Errorf("profile avatar not found for key: %s", tempKey),
			)
		}

		return "", fmt.Errorf("head object for profile avatar: %w", err)
	}

	if head.ContentLength == nil || *head.ContentLength == 0 {
		return "", errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("no content uploaded for profile avatar"),
		)
	}

	if *head.ContentLength > b.config.Profile.Avatar.ContentLengthMax {
		return "", errx.ErrorProfileAvatarContentIsInvalid.Raise(
			fmt.Errorf("uploaded profile avatar size %d exceeds the maximum allowed size", *head.ContentLength),
		)
	}

	rc, err := b.s3.GetObjectRange(ctx, tempKey, 2048)
	if err != nil {
		return "", fmt.Errorf("get object range for profile avatar: %w", err)
	}
	defer rc.Close()

	probe, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("read avatar probe bytes: %w", err)
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(probe))
	if err != nil {
		return "", fmt.Errorf("decode config: %w", err)
	}

	if cfg.Width > b.config.Profile.Avatar.MaxWidth {
		return "", errx.ErrorProfileAvatarContentIsInvalid.Raise(
			fmt.Errorf("uploaded profile avatar width %d exceeds the maximum allowed width", cfg.Width),
		)
	}

	if cfg.Height > b.config.Profile.Avatar.MaxHeight {
		return "", errx.ErrorProfileAvatarContentIsInvalid.Raise(
			fmt.Errorf("uploaded profile avatar height %d exceeds the maximum allowed height", cfg.Height),
		)
	}

	if !contains(b.config.Profile.Avatar.AllowedFormats, format) {
		return "", errx.ErrorProfileAvatarContentIsInvalid.Raise(
			fmt.Errorf("uploaded profile avatar format %s is not allowed", format),
		)
	}

	err = b.s3.CopyObject(ctx, tempKey, finalKey)
	if err != nil {
		return "", fmt.Errorf("copy object for profile avatar: %w", err)
	}

	_ = b.s3.DeleteObject(ctx, tempKey)

	return finalKey, nil
}

func (b Bucket) DeleteUploadProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
	key string,
) error {
	if err := validateTempProfileAvatarKey(accountID, key); err != nil {
		return err
	}

	if err := b.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("delete temp profile avatar: %w", err)
	}

	return nil
}

func validateTempProfileAvatarKey(accountID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	prefix := tempProfileAvatarPrefix(accountID)
	if !strings.HasPrefix(key, prefix) {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(
			fmt.Errorf("key does not belong to the account"),
		)
	}

	return nil
}

func tempProfileAvatarPrefix(accountID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/temp/", accountID.String())
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}
