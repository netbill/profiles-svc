package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/models"
)

type accountRepo interface {
	Create(ctx context.Context, account models.Account) (models.Account, error)
	Delete(ctx context.Context, accountID uuid.UUID) (models.Account, error)
}

type profileRepo interface {
	Create(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	Delete(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
}

type transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
