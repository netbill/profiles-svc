package profile

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) CancelUpdateSession(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) error {
	err := m.bucket.CleanProfileMediaSession(ctx, accountID, sessionID)
	if err != nil {
		return err
	}

	return nil
}
