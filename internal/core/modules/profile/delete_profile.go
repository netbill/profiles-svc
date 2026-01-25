package profile

import (
	"context"

	"github.com/google/uuid"
)

func (s Service) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.DeleteProfile(ctx, userID)
		if err != nil {
			return err
		}

		err = s.messanger.WriteProfileDeleted(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})
}
