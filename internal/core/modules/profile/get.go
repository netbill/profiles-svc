package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/core/models"
)

func (m *Module) GetMy(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	return m.repo.GetProfileByAccountID(ctx, accountID)
}

func (m *Module) GetByUsername(
	ctx context.Context,
	username string,
) (models.Profile, error) {
	return m.repo.GetProfileByUsername(ctx, username)
}
