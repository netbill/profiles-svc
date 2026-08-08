package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type profileRepo interface {
	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)

	Update(
		ctx context.Context,
		accountID uuid.UUID,
		params UpdateParams,
	) (models.Profile, error)

	Delete(ctx context.Context, accountID uuid.UUID) error

	Filter(
		ctx context.Context,
		params FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Profile], error)
}

type transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
