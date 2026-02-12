package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) GetByAccountID(ctx context.Context, userID uuid.UUID) (models.Profile, error) {
	return m.repo.GetProfileByAccountID(ctx, userID)
}

func (m *Module) GetByUsername(ctx context.Context, username string) (models.Profile, error) {
	return m.repo.GetProfileByUsername(ctx, username)
}
