package handler

import (
	"context"
	"encoding/json"

	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
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

	return h.modules.account.Create(ctx, account.CreateAccountParams{
		ID:        payload.AccountID,
		Username:  payload.Username,
		Role:      payload.Role,
		CreatedAt: payload.CreatedAt,
	})
}

func (h *Handler) AccountDeleted(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.AccountDeletedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.account.Delete(ctx, payload.AccountID)
}

func (h *Handler) AccountUsernameUpdated(
	ctx context.Context,
	event eventbox.InboxEvent,
) error {
	var payload evtypes.AccountUsernameUpdatedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	return h.modules.account.UpdateUsername(ctx, payload.AccountID, account.UpdateUsernameParams{
		Username:  payload.Username,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
}
