package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (s Service) GetProfileByID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	return s.repo.GetProfileByAccountID(ctx, userID)
}

func (s Service) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	return s.repo.GetProfileByUsername(ctx, username)
}
