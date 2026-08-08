package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/profiles-svc/internal/errx"
	"github.com/netbill/profiles-svc/internal/models"
)

const operationAccountCreated = "account_created"

func (h *Handler) AccountCreated(
	ctx context.Context,
	event eventbox.OutboxEvent,
) error {
	var payload evtypes.AccountCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationAccountCreated).
		With(slog.String("account_id", payload.Account.ID.String()))

	err := h.modules.Account.Create(ctx, models.Account{
		ID:        payload.Account.ID,
		Role:      payload.Account.Role,
		CreatedAt: payload.Account.CreatedAt,
	})
	switch {
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
	event eventbox.OutboxEvent,
) error {
	var payload evtypes.AccountDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithOperation(operationAccountDeleted).
		With(slog.String("account_id", payload.Account.ID.String()))

	err := h.modules.Account.Delete(ctx, payload.Account.ID)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to delete account: %v", err)
		return err
	default:
		log.Debug("account deleted successfully")
		return nil
	}
}
