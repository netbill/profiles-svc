package inbound

import (
	"context"
	"encoding/json"

	"github.com/netbill/profiles-svc/internal/core/modules/account"
	"github.com/netbill/profiles-svc/internal/messenger/evtypes"
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

	return i.modules.profile.Create(ctx, account.CreateAccountParams{
		ID:        payload.AccountID,
		Username:  payload.Username,
		Role:      payload.Role,
		Version:   payload.Version,
		CreatedAt: payload.CreatedAt,
	})
}
