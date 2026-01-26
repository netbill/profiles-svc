package bucket

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"time"

	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
)

const (
	ProfileAvatarMaxW          = 512
	ProfileAvatarMaxH          = 512
	ProfileAvatarProbeMaxBytes = int64(512 * 1024)

	ProfileContentLengthMax               = 5 * 1024 * 1024 // 5 MB
	ProfileAvatarUploadTTL  time.Duration = 1 * time.Hour
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

func isAllowedDetectedContentType(ct string) bool {
	for _, allowedCT := range allowedProfileAvatarContentTypes {
		if ct == allowedCT {
			return true
		}
	}
	return false
}

var allowedProfileAvatarExtensions = []string{
	"png",
	"jpeg",
	"jpg",
	"gif",
}

func isAllowedImageFormat(format string) bool {
	for _, allowedExt := range allowedProfileAvatarExtensions {
		if format == allowedExt {
			return true
		}
	}
	return false
}

func (r Bucket) GetPreloadLinkForUpdateProfileAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) (uploadURL, getUrl string, error error) {
	uploadURL, getURL, err := r.s3.PresignPut(
		ctx,
		CreateTempProfileAvatarKey(accountID, sessionID),
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
	tempKey := CreateTempProfileAvatarKey(accountID, sessionID)
	finalKey := CreateProfileAvatarKey(accountID)

	obj, err := r.s3.HeadObject(ctx, tempKey)
	if err != nil {
		var respErr *smithyhttp.ResponseError
		if errors.As(err, &respErr) && (respErr.HTTPStatusCode() == 404 || respErr.HTTPStatusCode() == 403) {
			return "", errx.ErrorNoAvatarUpload.Raise(
				fmt.Errorf("avatar upload not found for session %s", sessionID),
			)
		}
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("failed to head object for profile avatar: %w", err),
		)
	}

	cl := obj.ContentLength
	if cl == nil || *cl <= 0 || *cl > ProfileContentLengthMax {
		size := int64(-1)
		if cl != nil {
			size = *cl
		}
		return "", errx.ErrorContentLengthExceed.Raise(
			fmt.Errorf("profile avatar size %d exceeds max allowed size %d", size, ProfileContentLengthMax),
		)
	}

	rc, err := r.s3.GetObjectRange(ctx, tempKey, ProfileAvatarProbeMaxBytes)
	if err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get object range for profile avatar: %w", err),
		)
	}
	defer rc.Close()

	probe, err := io.ReadAll(rc)
	if err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("failed to read avatar probe bytes: %w", err),
		)
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(probe))
	if err != nil {
		return "", errx.ErrorContentTypeIsNotAllowed.Raise(
			fmt.Errorf("uploaded file is not a supported image: %w", err),
		)
	}

	if cfg.Width > ProfileAvatarMaxW || cfg.Height > ProfileAvatarMaxH {
		return "", errx.ErrorContentTypeIsNotAllowed.Raise(
			fmt.Errorf("profile avatar resolution %dx%d exceeds max allowed %dx%d", cfg.Width, cfg.Height, ProfileAvatarMaxW, ProfileAvatarMaxH),
		)
	}

	if !isAllowedImageFormat(format) {
		return "", errx.ErrorContentTypeIsNotAllowed.Raise(
			fmt.Errorf("profile avatar image format %s is not allowed", format),
		)
	}

	detectedCT := http.DetectContentType(probe)
	if !isAllowedDetectedContentType(detectedCT) {
		return "", errx.ErrorContentTypeIsNotAllowed.Raise(
			fmt.Errorf("profile avatar content type %s is not allowed", detectedCT),
		)
	}

	res, err := r.s3.CopyObject(ctx, tempKey, finalKey)
	if err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("failed to copy object for profile avatar: %w", err),
		)
	}

	err = r.s3.DeleteObject(ctx, tempKey)
	if err != nil {
		return "", errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete temp object for profile avatar: %w", err),
		)
	}

	return res, nil
}

func (r Bucket) CancelUpdateProfileAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) error {
	err := r.s3.DeleteObject(ctx, CreateTempProfileAvatarKey(accountID, sessionID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete temp object for profile avatar: %w", err,
		)
	}

	return nil
}

func (r Bucket) DeleteProfileAvatar(ctx context.Context, accountID uuid.UUID) error {
	err := r.s3.DeleteObject(ctx, CreateProfileAvatarKey(accountID))
	if err != nil {
		return fmt.Errorf(
			"failed to delete object for profile avatar: %w", err,
		)
	}

	return nil
}
