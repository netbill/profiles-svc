package bucket

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func CreateTempProfileAvatarKey(accountID, sessionID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/temp/%s", accountID, sessionID)
}

func CreateProfileAvatarKey(accountID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s", accountID)
}

func (b Bucket) GetPreloadLinkForProfileMedia(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) (models.UpdateProfileMediaLinks, error) {
	uploadURL, getURL, err := b.s3.PresignPut(
		ctx,
		CreateTempProfileAvatarKey(accountID, sessionID),
		b.config.Profile.Avatar.TTL,
	)
	if err != nil {
		return models.UpdateProfileMediaLinks{}, fmt.Errorf(
			"failed to presign put object for profile avatar: %w", err,
		)
	}

	return models.UpdateProfileMediaLinks{
		UploadURL: uploadURL,
		GetURL:    getURL,
	}, nil
}

func (b Bucket) UpdateProfileAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) (string, error) {
	tempKey := CreateTempProfileAvatarKey(accountID, sessionID)
	finalKey := CreateProfileAvatarKey(accountID)

	head, err := b.s3.HeadObject(ctx, tempKey)
	if err != nil {
		return "", fmt.Errorf("failed to head object for profile avatar: %w", err)
	}

	if head.ContentLength == nil || *head.ContentLength == 0 {
		return "", errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("no content uploaded for profile avatar in session %s", sessionID),
		)
	}

	rc, err := b.s3.GetObjectRange(ctx, tempKey, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to get object range for profile avatar: %w", err)
	}
	defer rc.Close()

	probe, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("failed to read avatar probe bytes: %w", err)
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(probe))
	if err != nil {
		return "", fmt.Errorf("decode config: %w", err)
	}

	if b.config.Profile.Avatar.MaxWidth > 0 && config.Width > b.config.Profile.Avatar.MaxWidth {
		return "", errx.ErrorProfileAvatarContentIsInvalid.Raise(
			fmt.Errorf("uploaded profile avatar width %d exceeds the maximum allowed width", config.Width),
		)
	}
	if b.config.Profile.Avatar.MaxHeight > 0 && config.Height > b.config.Profile.Avatar.MaxHeight {
		return "", errx.ErrorProfileAvatarContentIsInvalid.Raise(
			fmt.Errorf("uploaded profile avatar height %d exceeds the maximum allowed height", config.Height),
		)
	}

	access := func(values []string, target string) bool {
		for _, v := range values {
			if v == target {
				return true
			}
		}
		return false
	}

	if !access(b.config.Profile.Avatar.AllowedFormats, format) {
		return "", errx.ErrorProfileAvatarContentIsInvalid.Raise(
			fmt.Errorf("uploaded profile avatar format %s is not allowed", format),
		)
	}

	res, err := b.s3.CopyObject(ctx, tempKey, finalKey)
	if err != nil {
		return "", fmt.Errorf("failed to copy object for profile avatar: %w", err)
	}

	return res, nil
}

func (b Bucket) DeleteProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	err := b.s3.DeleteObject(ctx, CreateProfileAvatarKey(accountID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete object for profile avatar: %w", err,
		)
	}

	return nil
}

func (b Bucket) CleanProfileMediaSession(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) error {
	err := b.s3.DeleteObject(ctx, CreateTempProfileAvatarKey(accountID, sessionID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete temp object for profile avatar: %w", err,
		)
	}

	return nil
}

func (b Bucket) CancelUpdateProfileAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) error {
	err := b.s3.DeleteObject(ctx, CreateTempProfileAvatarKey(accountID, sessionID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete temp object for profile avatar: %w", err,
		)
	}

	return nil
}
