package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TombstoneRow struct {
	ID         uuid.UUID `db:"id"`
	EntityType string    `db:"entity_type"`
	EntityID   uuid.UUID `db:"entity_id"`
	DeletedAt  time.Time `db:"deleted_at"`
}

type TombstonesSql interface {
	BuryAccount(ctx context.Context, accountID uuid.UUID) error
	AccountIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)
}
