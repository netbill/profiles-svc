package pg

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/pgdbx"
	"github.com/netbill/profiles-svc/internal/repository"
)

type tombstones struct {
	db *pgdbx.DB
}

func NewTombstonesQ(db *pgdbx.DB) repository.TombstonesSql {
	return &tombstones{db: db}
}

func (t *tombstones) BuryAccount(ctx context.Context, accountID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
		INSERT INTO tombstones (entity_type, entity_id)
		SELECT 'account', $1
		UNION ALL
		SELECT 'profile', p.account_id FROM profiles p WHERE p.account_id = $1
		ON CONFLICT (entity_type, entity_id) DO NOTHING
	`, accountID)
	if err != nil {
		return fmt.Errorf("burying account: %w", err)
	}

	return nil
}

func (t *tombstones) AccountIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM tombstones
			WHERE entity_type = 'account' AND entity_id = $1
		)
	`, accountID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking account is buried: %w", err)
	}

	return exists, nil
}
