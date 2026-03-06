package repository

import (
	"context"

	"github.com/google/uuid"
)

type TombstonesSql interface {
	BuryAccount(ctx context.Context, accountID uuid.UUID) error
	AccountIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)
}
