package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/models"
)

type accountRepo interface {
	Create(
		ctx context.Context,
		params CreateAccountParams,
	) (models.Account, error)

	GetByID(ctx context.Context, accountID uuid.UUID) (models.Account, error)
	ExistsByID(ctx context.Context, accountID uuid.UUID) (bool, error)

	UpdateUsername(
		ctx context.Context,
		accountID uuid.UUID,
		params UpdateUsernameParams,
	) (models.Account, error)
}

type profileRepo interface {
	Create(
		ctx context.Context,
		accountID uuid.UUID,
		username string,
	) (models.Profile, error)

	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	ExistsByID(ctx context.Context, accountID uuid.UUID) (bool, error)

	// IsDeleted reports whether this account was already deleted (soft-delete
	// on profiles.deleted_at), so a late/redelivered AccountCreated or
	// AccountUsernameUpdated event doesn't resurrect it.
	IsDeleted(ctx context.Context, accountID uuid.UUID) (bool, error)

	UpdateUsername(
		ctx context.Context,
		accountID uuid.UUID,
		username string,
	) (models.Profile, error)

	Delete(ctx context.Context, accountID uuid.UUID) error
}

type transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
