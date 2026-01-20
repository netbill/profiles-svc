package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/models"
)

type CreateParams struct {
	Username    string  `json:"username"`
	Pseudonym   *string `json:"pseudonym,omitempty"`
	Description *string `json:"description,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}

func (s Service) CreateProfile(ctx context.Context, userID uuid.UUID, params CreateParams) (profile models.Profile, err error) {
	err = s.checkUsernameRequirements(ctx, params.Username)
	if err != nil {
		return models.Profile{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		profile, err = s.repo.CreateProfile(ctx, userID, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("creating profile for user '%s': %w", userID, err),
			)
		}

		err = s.messanger.WriteProfileCreated(ctx, profile)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("creating profile for user '%s': %w", userID, err),
			)
		}

		return nil
	}); err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}
