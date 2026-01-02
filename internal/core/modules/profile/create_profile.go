package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (s Service) CreateProfile(ctx context.Context, userID uuid.UUID, username string) (models.Profile, error) {
	profile, err := s.repo.CreateProfile(ctx, userID, username)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("creating profile for user '%s': %w", userID, err),
		)
	}

	return profile, nil
}
