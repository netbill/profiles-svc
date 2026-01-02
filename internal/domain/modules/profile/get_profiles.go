package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/domain/errx"
	"github.com/netbill/profiles-svc/internal/domain/models"
)

func (s Service) GetProfileByID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	profile, err := s.repo.GetProfileByAccountID(ctx, userID)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting profile for user '%s': %w", userID, err),
		)
	}

	if profile.IsNil() {
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile for user '%s' does not exist", userID),
		)
	}

	return profile, nil
}

func (s Service) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	profile, err := s.repo.GetProfileByUsername(ctx, username)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("getting profile with username '%s': %w", username, err),
		)
	}

	if profile.IsNil() {
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile with username '%s' does not exist", username),
		)
	}

	return profile, nil
}
