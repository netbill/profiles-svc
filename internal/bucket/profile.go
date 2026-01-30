package bucket

import (
	"context"
	"fmt"
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
		b.tokensTTL.ProfileAvatar,
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

func (b Bucket) AcceptUpdateProfileMedia(ctx context.Context, accountID, sessionID uuid.UUID) (string, error) {
	tempKey := CreateTempProfileAvatarKey(accountID, sessionID)
	finalKey := CreateProfileAvatarKey(accountID)

	rc, size, err := b.s3.GetObjectRange(ctx, tempKey, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to get object range for profile avatar: %w", err)
	}
	defer rc.Close()

	if size == 0 {
		return "", errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("no content uploaded for profile avatar in session %s", sessionID),
		)
	}

	probe, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("failed to read avatar probe bytes: %w", err)
	}

	valid, err := b.profileAvatarValidator.ValidateImageSize(uint(size))
	if err != nil {
		return "", fmt.Errorf("failed to validate profile avatar image size: %w", err)
	}
	if !valid {
		return "", errx.ErrorProfileAvatarTooLarge.Raise(
			fmt.Errorf("uploaded profile avatar size %d exceeds the maximum allowed size", size),
		)
	}

	valid, err = b.profileAvatarValidator.ValidateImage(probe)
	if err != nil {
		return "", fmt.Errorf("failed to validate profile avatar image: %w", err)
	}
	if !valid {
		return "", errx.ErrorProfileAvatarContentTypeIsNotAllowed.Raise(
			fmt.Errorf("uploaded file is not a valid image"),
		)
	}

	valid, err = b.profileAvatarValidator.ValidateImageFormat(probe)
	if err != nil {
		return "", fmt.Errorf("failed to validate profile avatar image format: %w", err)
	}
	if !valid {
		return "", errx.ErrorProfileAvatarContentFormatIsNotAllowed.Raise(
			fmt.Errorf("profile avatar image format is not allowed"),
		)
	}

	valid, err = b.profileAvatarValidator.ValidateImageContentType(probe)
	if err != nil {
		return "", fmt.Errorf("failed to validate profile avatar content type: %w", err)
	}
	if !valid {
		return "", errx.ErrorProfileAvatarContentTypeIsNotAllowed.Raise(
			fmt.Errorf("profile avatar content type is not allowed"),
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
