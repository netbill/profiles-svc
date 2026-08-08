package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/modules/account"
)

const operationAccountCreated = "account_created"

func (h *Handler) AccountCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.AccountCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationAccountCreated).
		With(slog.String("account_id", payload.AccountID.String()))

	err := h.modules.Account.Create(ctx, account.CreateAccountParams{
		ID:        payload.AccountID,
		Username:  payload.Username,
		Role:      payload.Role,
		CreatedAt: payload.CreatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorAccountAlreadyExists):
		// Covers both a genuine duplicate and a replayed AccountCreated for
		// an account that was since deleted — id is never reused, so either
		// way there's already a row and nothing to do.
		log.Debug("received account created event for already existing account")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to create account: %v", err)
		return err
	default:
		log.Info("account created successfully")
		return nil
	}
}

const operationAccountDeleted = "account_deleted"

func (h *Handler) AccountDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.AccountDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationAccountDeleted).
		With(slog.String("account_id", payload.AccountID.String()))

	err := h.modules.Account.Delete(ctx, payload.AccountID)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to delete account: %v", err)
		return err
	default:
		// Also covers a redelivered AccountDeleted for an already-deleted
		// account — Service.Delete treats that as a silent no-op.
		log.Debug("account deleted successfully")
		return nil
	}
}

const operationAccountUpdated = "account_updated"

func (h *Handler) AccountUsernameUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.AccountUsernameUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationAccountUpdated).
		With(slog.String("account_id", payload.AccountID.String()))

	err := h.modules.Account.UpdateUsername(ctx, payload.AccountID, account.UpdateUsernameParams{
		Username:  payload.Username,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
	switch {
	case err != nil:
		log.WithError(err).Error("failed to update account username: %v", err)
		return err
	default:
		// Also covers a deleted account or a stale/replayed event — both are
		// silently skipped inside Service.UpdateUsername.
		log.Debug("account username updated successfully")
		return nil
	}
}
