package account

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
)

type accountMessenger interface {
	WriteProfileCreated(ctx context.Context, profile models.Profile) error
	WriteProfileUpdated(ctx context.Context, profile models.Profile) error
	WriteProfileDeleted(ctx context.Context, accountID uuid.UUID) error
}

type Service struct {
	account accountRepo
	profile profileRepo
	tx      transaction

	messenger accountMessenger
}

type ServiceDeps struct {
	AccountRepo accountRepo
	ProfileRepo profileRepo
	Transaction transaction

	Messenger accountMessenger
}

func NewAccountModule(deps ServiceDeps) *Service {
	return &Service{
		account: deps.AccountRepo,
		profile: deps.ProfileRepo,
		tx:      deps.Transaction,

		messenger: deps.Messenger,
	}
}

type CreateAccountParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`

	CreatedAt time.Time `json:"created_at"`
}

// Create relies on accounts.id's PRIMARY KEY to reject a duplicate/replayed
// AccountCreated atomically (see AccountRepo.Create) — id is never reused,
// so any existing row, active or soft-deleted, means "already applied".
func (m *Service) Create(
	ctx context.Context,
	params CreateAccountParams,
) error {
	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		_, err := m.account.Create(ctx, params)
		if err != nil {
			return err
		}

		profile, err := m.profile.Create(ctx, params.ID, params.Username)
		if err != nil {
			return err
		}

		return m.messenger.WriteProfileCreated(ctx, profile)
	})
}

type UpdateUsernameParams struct {
	Username  string
	Version   int32
	UpdatedAt time.Time
}

// UpdateUsername applies the update only if AccountRepo.UpdateUsername's own
// WHERE (version < params.Version AND deleted_at IS NULL) matches — a stale/
// replayed event or an account deleted concurrently both surface as
// errx.ErrorAccountNotExists and are silently skipped, same as before.
func (m *Service) UpdateUsername(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateUsernameParams,
) error {
	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		_, err := m.account.UpdateUsername(ctx, accountID, params)
		switch {
		case errors.Is(err, errx.ErrorAccountNotExists):
			return nil
		case err != nil:
			return err
		}

		profile, err := m.profile.UpdateUsername(ctx, accountID, params.Username)
		if err != nil {
			return err
		}

		return m.messenger.WriteProfileUpdated(ctx, profile)
	})
}

// Delete soft-deletes both mirror rows. profile.Delete gates the whole
// operation: if it reports errx.ErrorProfileNotExists (already deleted by
// an earlier delivery of the same event), the rest is skipped as a no-op.
func (m *Service) Delete(ctx context.Context, accountID uuid.UUID) error {
	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		_, err := m.profile.Delete(ctx, accountID)
		switch {
		case errors.Is(err, errx.ErrorProfileNotExists):
			return nil
		case err != nil:
			return err
		}

		if _, err = m.account.Delete(ctx, accountID); err != nil {
			return err
		}

		return m.messenger.WriteProfileDeleted(ctx, accountID)
	})
}
