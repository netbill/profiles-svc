package bucket

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/tokenmanager"
)

const (
	ProfileAvatarMaxW          = 512
	ProfileAvatarMaxH          = 512
	ProfileAvatarProbeMaxBytes = int64(512 * 1024)

	ProfileContentLengthMax = 5 * 1024 * 1024 // 5 MB
)

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
	"image/gif",
}

var allowedProfileAvatarFormats = []string{
	"png",
	"jpeg",
	"jpg",
	"gif",
}

func allowed(value string, allowed []string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || len(allowed) == 0 {
		return false
	}
	for _, a := range allowed {
		a = strings.ToLower(strings.TrimSpace(a))
		if value == a {
			return true
		}
	}
	return false
}

func (b Bucket) GetPreloadLinkForUpdateProfileAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) (uploadURL, getUrl string, error error) {
	uploadURL, getURL, err := b.s3.PresignPut(
		ctx,
		CreateTempProfileAvatarKey(accountID, sessionID),
		tokenmanager.ProfileMediaUploadTTL,
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"failed to presign put object for profile avatar: %w", err,
		)
	}

	return uploadURL, getURL, nil
}

func (b Bucket) AcceptUpdateProfileAvatar(ctx context.Context, accountID, sessionID uuid.UUID) (string, error) {
	tempKey := CreateTempProfileAvatarKey(accountID, sessionID)
	finalKey := CreateProfileAvatarKey(accountID)

	rc, size, err := b.s3.GetObjectRange(ctx, tempKey, ProfileAvatarProbeMaxBytes)
	if err != nil {
		return "", fmt.Errorf("failed to get object range for profile avatar: %w", err)
	}
	defer rc.Close()

	switch {
	case size == 0:
		return "", errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("no content uploaded for profile avatar in session %s", sessionID),
		)
	case size > ProfileContentLengthMax:
		return "", errx.ErrorProfileAvatarTooLarge.Raise(
			fmt.Errorf("profile avatar size %d exceeds max allowed size %d", size, ProfileContentLengthMax),
		)
	}

	probe, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("failed to read avatar probe bytes: %w", err)
	}

	img, format, err := image.DecodeConfig(bytes.NewReader(probe))
	if err != nil {
		return "", errx.ErrorProfileAvatarContentTypeIsNotAllowed.Raise(
			fmt.Errorf("uploaded file is not a supported image: %w", err),
		)
	}

	if img.Width > ProfileAvatarMaxW || img.Height > ProfileAvatarMaxH {
		return "", errx.ErrorProfileAvatarContentTypeIsNotAllowed.Raise(
			fmt.Errorf(
				"profile avatar resolution %dx%d exceeds max allowed %dx%d",
				img.Width, img.Height, ProfileAvatarMaxW, ProfileAvatarMaxH,
			),
		)
	}

	if !allowed(format, allowedProfileAvatarFormats) {
		return "", errx.ErrorProfileAvatarContentFormatIsNotAllowed.Raise(
			fmt.Errorf("profile avatar image format %s is not allowed", format),
		)
	}

	ct := http.DetectContentType(probe)
	if !allowed(ct, allowedProfileAvatarContentTypes) {
		return "", errx.ErrorProfileAvatarContentTypeIsNotAllowed.Raise(
			fmt.Errorf("profile avatar content type %s is not allowed", ct),
		)
	}

	res, err := b.s3.CopyObject(ctx, tempKey, finalKey)
	if err != nil {
		return "", fmt.Errorf("failed to copy object for profile avatar: %w", err)
	}

	return res, nil
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

func (b Bucket) DeleteProfileAvatar(ctx context.Context, accountID uuid.UUID) error {
	err := b.s3.DeleteObject(ctx, CreateProfileAvatarKey(accountID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete object for profile avatar: %w", err,
		)
	}

	return nil
}
