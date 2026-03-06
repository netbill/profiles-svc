package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/profiles-svc/internal/core/account"
	"github.com/netbill/profiles-svc/internal/errx"
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
	case errors.Is(err, errx.ErrorAccountDeleted):
		log.Debug("received account already deleted account")
		return nil
	case errors.Is(err, errx.ErrorAccountAlreadyExists):
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
	case errors.Is(err, errx.ErrorAccountDeleted):
		log.Debug("received account deleted event for already deleted account")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to delete account: %v", err)
		return err
	default:
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
	case errors.Is(err, errx.ErrorAccountDeleted):
		log.Debug("received account username updated event for already deleted account")
		return nil
	case err != nil:
		log.WithError(err).Error("failed to update account username: %v", err)
		return err
	default:
		log.Debug("account username updated successfully")
		return nil
	}
}
