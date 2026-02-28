package handler

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/profiles-svc/internal/core/errx"
	"github.com/netbill/profiles-svc/internal/core/modules/account"
)

func (h *Handler) AccountCreated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.AccountCreatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	log := h.log.WithInboxEvent(event)

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
		return err
	default:
		log.Debug("account created successfully")
		return nil
	}
}

func (h *Handler) AccountDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.AccountDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	err := h.modules.Account.Delete(ctx, payload.AccountID)
	switch {
	case errors.Is(err, errx.ErrorAccountDeleted):
		h.log.WithInboxEvent(event).Debug("received account deleted event for already deleted account")
		return nil
	case err != nil:
		return err
	default:
		h.log.WithInboxEvent(event).Debug("account deleted successfully")
		return nil
	}
}

func (h *Handler) AccountUsernameUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.AccountUsernameUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	err := h.modules.Account.UpdateUsername(ctx, payload.AccountID, account.UpdateUsernameParams{
		Username:  payload.Username,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
	switch {
	case errors.Is(err, errx.ErrorAccountDeleted):
		h.log.WithInboxEvent(event).Debug("received account username updated event for already deleted account")
		return nil
	case err != nil:
		return err
	default:
		h.log.WithInboxEvent(event).Debug("account username updated successfully")
		return nil
	}
}
