package profile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		err = s.messanger.WriteProfileUpdated(ctx, profile)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("sending profile updated event for user '%s': %w", accountID, err),
			)
		}

		profile, err = s.repo.UpdateProfile(ctx, accountID, input)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating profile for user '%s': %w", accountID, err),
			)
		}

		return nil
	})

	return profile, nil
}

func (s Service) UpdateProfileOfficial(ctx context.Context, accountID uuid.UUID, official bool) (models.Profile, error) {
	profile, err := s.GetProfileByID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}

	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		err = s.messanger.WriteProfileUpdated(ctx, profile)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("sending profile updated event for user '%s': %w", accountID, err),
			)
		}

		profile, err = s.repo.UpdateProfileOfficial(ctx, accountID, official)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("updating profile for user '%s': %w", accountID, err),
			)
		}

		return nil
	})

	return profile, nil
}

func (s Service) UpdateProfileUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error) {
	profile, err := s.repo.UpdateProfileUsername(ctx, accountID, username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Profile{}, errx.ErrorProfileNotFound.Raise(
				fmt.Errorf("profile for user '%s' does not exist", accountID),
			)
		default:
			return models.Profile{}, errx.ErrorInternal.Raise(
				fmt.Errorf("updating username for user '%s': %w", accountID, err),
			)
		}
	}

	return profile, nil
}
