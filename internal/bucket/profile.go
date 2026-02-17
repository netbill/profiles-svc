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

func CreateTempProfileAvatarKey(accountID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/temp/%s", accountID, uuid.New())
}

func CreateFinalProfileAvatarKey(accountID uuid.UUID) string {
	return fmt.Sprintf("profile/avatar/%s/%s", accountID, uuid.New())
}

func (b Bucket) CreateProfileAvatarUploadMediaLinks(
	ctx context.Context,
	accountID uuid.UUID,
) (models.UploadMediaLink, error) {
	key := CreateTempProfileAvatarKey(accountID)

	uploadLink, getLink, err := b.s3.PresignPut(ctx, key, b.config.Media.Link.TTL)
	if err != nil {
		return models.UploadMediaLink{}, fmt.Errorf("presign put object for profile avatar: %w", err)
	}

	return models.UploadMediaLink{
		Key:        key,
		PreloadUrl: getLink,
		UploadURL:  uploadLink,
	}, nil
}

func (b Bucket) ValidateProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
	tempKey string,
) error {
	if err := validateTempProfileAvatarKey(accountID, tempKey); err != nil {
		return err
	}

	out, err := b.s3.GetObjectRange(ctx, tempKey, 64*1024)
	switch {
	case errors.Is(err, awsx.ErrNotFound):
		return errx.ErrorNoContentUploaded.Raise(
			fmt.Errorf("profile avatar not found for key: %s", tempKey),
		)
	case err != nil:
		return fmt.Errorf("get object range for profile avatar: %w", err)
	}
	defer out.Body.Close()

	if err = b.config.Media.Profile.Avatar.Validate(out); err != nil {
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

func (b Bucket) DeleteUploadProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
	tempKey string,
) error {
	if err := validateTempProfileAvatarKey(accountID, tempKey); err != nil {
		return err
	}

	if err := b.s3.DeleteObject(ctx, tempKey); err != nil {
		return fmt.Errorf("delete temp profile avatar: %w", err)
	}

	return nil
}

func (b Bucket) DeleteProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
	finalKey string,
) error {
	if err := validateFinalProfileAvatarKey(accountID, finalKey); err != nil {
		return err
	}

	if err := b.s3.DeleteObject(ctx, finalKey); err != nil {
		return fmt.Errorf("delete profile avatar: %w", err)
	}

	return nil
}

func (b Bucket) UpdateProfileAvatar(
	ctx context.Context,
	accountID uuid.UUID,
	oldFinalKey *string,
	tempKey *string,
) (*string, error) {
	if ptrStrEq(oldFinalKey, tempKey) {
		return oldFinalKey, nil
	}

	if tempKey == nil {
		return nil, b.DeleteProfileAvatar(ctx, accountID, *oldFinalKey)
	}

	if err := b.ValidateProfileAvatar(ctx, accountID, *tempKey); err != nil {
		return nil, err
	}

	finalKey := CreateFinalProfileAvatarKey(accountID)

	if err := b.s3.CopyObject(ctx, *tempKey, finalKey); err != nil {
		return nil, fmt.Errorf("copy object for profile avatar: %w", err)
	}

	err := b.s3.DeleteObject(ctx, *tempKey)
	if err != nil {
		return nil, fmt.Errorf("delete temp profile avatar: %w", err)
	}

	if oldFinalKey != nil {
		if err = b.DeleteProfileAvatar(ctx, accountID, *oldFinalKey); err != nil {
			return nil, err
		}
	}

	return &finalKey, nil
}

var (
	tempAvatarKeyRe = regexp.MustCompile(
		`^profile/avatar/([0-9a-fA-F-]{36})/temp/([0-9a-fA-F-]{36})$`,
	)
	finalAvatarKeyRe = regexp.MustCompile(
		`^profile/avatar/([0-9a-fA-F-]{36})/([0-9a-fA-F-]{36})$`,
	)
)

func validateTempProfileAvatarKey(accountID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := tempAvatarKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("invalid key format"))
	}

	if matches[1] != accountID.String() {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("key does not belong to the account"))
	}

	return nil
}

func validateFinalProfileAvatarKey(accountID uuid.UUID, key string) error {
	if key == "" {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("empty key"))
	}

	matches := finalAvatarKeyRe.FindStringSubmatch(key)
	if matches == nil {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("invalid key format"))
	}

	if matches[1] != accountID.String() {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("key does not belong to the account"))
	}

	if tempAvatarKeyRe.MatchString(key) {
		return errx.ErrorProfileAvatarKeyIsInvalid.Raise(fmt.Errorf("final key cannot be temp key"))
	}

	return nil
}
