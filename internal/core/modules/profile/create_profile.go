package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (s Service) CreateProfile(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error) {
	profile, err := s.repo.GetProfileByAccountID(ctx, accountID)
	switch {
	case errors.Is(err, errx.ErrorProfileNotFound):
		// continue to create profile
	case err != nil:
		return models.Profile{}, err
	default:
		return profile, nil
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.InsertProfile(ctx, accountID, username)
		if err != nil {
			return err
		}

		err = s.messanger.WriteProfileCreated(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
