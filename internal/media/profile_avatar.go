package media

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/netbill/awsx"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
)

func (s *Uploader) CreateProfileAvatarUploadMediaLinks(
	ctx context.Context,
	profileID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempProfileAvatarKey(profileID)

	uploadURL, getURL, err := s.s3.PresignPut(
		ctx,
		key,
		s.config.LinkTTL,
	)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presigning put for profile profile avatar: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getURL,
		UploadURL:  uploadURL,
	}, nil
}

func (s *Uploader) UpdateProfileAvatar(
	ctx context.Context,
	profileID uuid.UUID,
	key string,
) (string, error) {
	err := validateTempProfileAvatarKey(profileID, key)
	if err != nil {
		return "", err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return "", errx.ErrorProfileUploadedAvatarInvalid.Raise(
			fmt.Errorf("profile profile avatar not found for key: %s", key),
		)
	case err != nil:
		return "", fmt.Errorf("get object range for profile profile avatar: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.ProfileAvatar.Validate(out); err != nil {
		return "", errx.ErrorProfileUploadedAvatarInvalid.Raise(
			fmt.Errorf("validating profile profile avatar: %w", err),
		)
	}

	finalKey := CreateProfileAvatarKey(profileID)

	if err = s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for profile avatar: %w", err)
	}

	return finalKey, nil
}

func (s *Uploader) DeleteUploadProfileAvatar(
	ctx context.Context,
	profileID uuid.UUID,
	key string,
) error {
	if err := validateTempProfileAvatarKey(profileID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting temp profile profile avatar object: %w", err)
	}

	return nil
}

func (s *Uploader) DeleteProfileAvatar(
	ctx context.Context,
	profileID uuid.UUID,
	key string,
) error {
	if err := validateFinalProfileAvatarKey(profileID, key); err != nil {
		return err
	}

	if err := s.s3.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("deleting profile profile avatar object: %w", err)
	}

	return nil
}

var tempProfileAvatarKeyRe = regexp.MustCompile(
	`^profile/avatar/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
)

func CreateTempProfileAvatarKey(profileID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/temp/%s", profileID, uuid.New().String())
}

func validateTempProfileAvatarKey(profileID uuid.UUID, key string) error {
	matches := tempProfileAvatarKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorProfileUploadedAvatarInvalid.Raise(fmt.Errorf("key %s does not match temp profile profile avatar key pattern", key))
	}

	if matches[1] != profileID.String() {
		return errx.ErrorProfileUploadedAvatarInvalid.Raise(fmt.Errorf("key %s does not belong to profile profile %s", key, profileID))
	}

	return nil
}

var finalProfileAvatarKeyRe = regexp.MustCompile(
	`^profile/avatar/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
)

func CreateProfileAvatarKey(profileID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/%s", profileID, uuid.New().String())
}

func validateFinalProfileAvatarKey(profileID uuid.UUID, key string) error {
	matches := finalProfileAvatarKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorProfileUploadedAvatarInvalid.Raise(fmt.Errorf("key %s does not match final profile profile avatar key pattern", key))
	}

	if matches[1] != profileID.String() {
		return errx.ErrorProfileUploadedAvatarInvalid.Raise(fmt.Errorf("key %s does not belong to profile profile %s", key, profileID))
	}

	return nil
}
