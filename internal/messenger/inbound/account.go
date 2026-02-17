package inbound

import (
	"context"
	"encoding/json"

	"github.com/netbill/profiles-svc/internal/core/modules/account"
	"github.com/netbill/profiles-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

func (i *Inbound) AccountCreated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.AccountCreatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	return i.modules.account.Create(ctx, account.CreateAccountParams{
		ID:        payload.AccountID,
		Username:  payload.Username,
		Role:      payload.Role,
		Version:   payload.Version,
		CreatedAt: payload.CreatedAt,
	})
}

func (i *Inbound) AccountDeleted(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.AccountDeletedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	return i.modules.account.Delete(ctx, payload.AccountID)
}

func (i *Inbound) AccountUsernameUpdated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.AccountUsernameUpdatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	return i.modules.account.UpdateUsername(ctx, payload.AccountID, account.UpdateUsernameParams{
		Username:  payload.NewUsername,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
}
