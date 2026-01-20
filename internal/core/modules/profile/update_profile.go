package profile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"unicode"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

type UpdateParams struct {
	Pseudonym   *string
	Description *string
	Avatar      *string
}

func (s Service) UpdateProfile(ctx context.Context, accountID uuid.UUID, input UpdateParams) (models.Profile, error) {
	profile, err := s.GetProfileByID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.UpdateProfile(ctx, accountID, input)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating profile for user '%s': %w", accountID, err),
			)
		}

		err = s.messanger.WriteProfileUpdated(ctx, profile)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("sending profile updated event for user '%s': %w", accountID, err),
			)
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (s Service) UpdateProfileOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error) {
	profile, err := s.GetProfileByID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.UpdateProfileOfficial(ctx, accountID, official)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating profile for user '%s': %w", accountID, err),
			)
		}

		err = s.messanger.WriteProfileOfficialUpdated(ctx, profile)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("sending profile updated event for user '%s': %w", accountID, err),
			)
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (s Service) checkUsernameRequirements(ctx context.Context, username string) error {
	res, err := s.repo.GetProfileByUsername(ctx, username)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("fetching user '%s' from db: %w", username, err),
		)
	}
	if !res.IsNil() {
		return errx.ErrorUsernameAlreadyTaken.Raise(
			fmt.Errorf("user '%s' is already taken", username),
		)
	}

	if len(username) < 3 || len(username) > 32 {
		return errx.ErrorUsernameIsNotAllowed.Raise(
			fmt.Errorf("username must be between 3 and 32 characters"),
		)
	}

	for _, r := range username {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-') {
			return errx.ErrorUsernameIsNotAllowed.Raise(
				fmt.Errorf("username contains invalid characters %s", string(r)),
			)
		}
	}

	return nil
}

func (s Service) UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error) {
	profile, err := s.GetProfileByID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}

	if err = s.checkUsernameRequirements(ctx, username); err != nil {
		return models.Profile{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.UpdateProfileUsername(ctx, accountID, username)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return errx.ErrorProfileNotFound.Raise(
					fmt.Errorf("profile for user '%s' does not exist", accountID),
				)
			default:
				return errx.ErrorInternal.Raise(
					fmt.Errorf("updating username for user '%s': %w", accountID, err),
				)
			}
		}

		err = s.messanger.WriteProfileUsernameUpdated(ctx, profile)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("sending profile updated event for user '%s': %w", accountID, err),
			)
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
