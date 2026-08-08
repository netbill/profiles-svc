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

	// UpdateUsername and Delete filter "deleted_at IS NULL" (and, for
	// UpdateUsername, "version < params.Version") in their own WHERE clause,
	// so a stale/replayed event or a concurrent Delete atomically yields
	// errx.ErrorAccountNotExists instead of a separate check-then-act race.
	UpdateUsername(
		ctx context.Context,
		accountID uuid.UUID,
		params UpdateUsernameParams,
	) (models.Account, error)

	Delete(ctx context.Context, accountID uuid.UUID) (models.Account, error)
}

type profileRepo interface {
	Create(
		ctx context.Context,
		accountID uuid.UUID,
		username string,
	) (models.Profile, error)

	// UpdateUsername and Delete filter "deleted_at IS NULL" in their own
	// WHERE clause — see accountRepo.
	UpdateUsername(
		ctx context.Context,
		accountID uuid.UUID,
		username string,
	) (models.Profile, error)

	Delete(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
}

type transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
