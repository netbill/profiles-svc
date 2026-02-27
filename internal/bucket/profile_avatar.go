package bucket

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/netbill/awsx"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func CreateTempProfileAvatarKey(profileID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/temp/%s", profileID, uuid.New().String())
}

func CreateProfileAvatarKey(profileID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/%s", profileID, uuid.New().String())
}

func (s *Storage) CreateProfileAvatarUploadMediaLinks(
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

func (s *Storage) ValidateUploadProfileAvatar(
	ctx context.Context,
	profileID uuid.UUID,
	key string,
) error {
	err := validateTempProfileAvatarKey(profileID, key)
	if err != nil {
		return err
	}

	out, err := s.s3.GetObjectRange(ctx, key, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("profile profile avatar not found for key: %s", key),
		)
	case err != nil:
		return fmt.Errorf("get object range for profile profile avatar: %w", err)
	}
	defer out.Body.Close()

	if err = s.config.ProfileAvatar.Validate(out); err != nil {
		switch {
		case errors.Is(err, awsx.ErrorNoContentUploaded):
			return errx.ErrorNoContentUploaded.Raise(err)
		case errors.Is(err, awsx.ErrorSizeExceedsMax):
			return errx.ErrorProfileAvatarContentIsExceedsMax.Raise(err)
		case errors.Is(err, awsx.ErrorResolutionIsInvalid):
			return errx.ErrorProfileAvatarResolutionIsInvalid.Raise(err)
		case errors.Is(err, awsx.ErrorFormatNotAllowed):
			return errx.ErrorProfileAvatarFormatIsNotAllowed.Raise(err)
		default:
			return fmt.Errorf("validate profile avatar content: %w", err)
		}
	}

	return nil
}

func (s *Storage) DeleteUploadProfileAvatar(
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

func (s *Storage) DeleteProfileAvatar(
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

func (s *Storage) UpdateProfileAvatar(
	ctx context.Context,
	profileID uuid.UUID,
	key string,
) (string, error) {
	if err := validateTempProfileAvatarKey(profileID, key); err != nil {
		return "", err
	}

	finalKey := CreateProfileAvatarKey(profileID)

	if err := s.s3.CopyObject(ctx, key, finalKey); err != nil {
		return "", fmt.Errorf("copying object for profile avatar: %w", err)
	}

	return finalKey, nil
}

var (
	tempProfileAvatarKeyRe = regexp.MustCompile(
		`^profile/avatar/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
	)

	finalProfileAvatarKeyRe = regexp.MustCompile(
		`^profile/avatar/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
	)
)

func validateTempProfileAvatarKey(profileID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempProfileAvatarKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("key %s does not match temp profile profile avatar key pattern", key))
	}

	if matches[1] != profileID.String() {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to profile profile %s", key, profileID))
	}

	return nil
}

func validateFinalProfileAvatarKey(profileID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalProfileAvatarKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("key %s does not match final profile profile avatar key pattern", key))
	}

	if matches[1] != profileID.String() {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("key %s does not belong to profile profile %s", key, profileID))
	}

	return nil
}
