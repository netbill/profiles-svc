package account

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
)

type accountMessenger interface {
	WriteProfileCreated(ctx context.Context, profile models.Profile) error
	WriteProfileUpdated(ctx context.Context, profile models.Profile) error
	WriteProfileDeleted(ctx context.Context, profile models.Profile) error
}

type usernameReg interface {
	Validate(username string) error
	Generate() string
}

type Service struct {
	account accountRepo
	profile profileRepo
	tx      transaction

	messenger accountMessenger
	username  usernameReg
}

type ServiceDeps struct {
	AccountRepo accountRepo
	ProfileRepo profileRepo
	Transaction transaction

	Messenger accountMessenger
	Username  usernameReg
}

func NewAccountModule(deps ServiceDeps) *Service {
	return &Service{
		account: deps.AccountRepo,
		profile: deps.ProfileRepo,
		tx:      deps.Transaction,

		messenger: deps.Messenger,
		username:  deps.Username,
	}
}

func (m *Service) Create(
	ctx context.Context,
	account models.Account,
) error {
	candidate := m.username.Generate()

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if _, err := m.account.Create(ctx, account); err != nil {
			return err
		}

		profile, err := m.profile.Create(ctx, account.ID, candidate)
		if err != nil {
			return err
		}

		return m.messenger.WriteProfileCreated(ctx, profile)
	})
}

// Delete soft-deletes both mirror rows. profile.Delete gates the whole
// operation: if it reports errx.ErrorProfileNotExists (already deleted by
// an earlier delivery of the same event), the rest is skipped as a no-op —
// account.Delete alone can't tell "already deleted" from "not created yet"
// apart (both match zero rows), so it must not run first.
func (m *Service) Delete(ctx context.Context, accountID uuid.UUID) error {
	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		profile, err := m.profile.Delete(ctx, accountID)
		switch {
		case errors.Is(err, errx.ErrorProfileNotExists):
			return nil
		case err != nil:
			return err
		}

		if _, err = m.account.Delete(ctx, accountID); err != nil {
			return err
		}

		return m.messenger.WriteProfileDeleted(ctx, profile)
	})
}
