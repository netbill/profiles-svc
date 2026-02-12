package profile

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) DeleteUploadAvatar(
	ctx context.Context,
	accountID, sessionID uuid.UUID,
) error {
	err := m.bucket.CancelUpdateProfileAvatar(ctx, accountID, sessionID)
	if err != nil {
		return err
	}

	return nil
}
