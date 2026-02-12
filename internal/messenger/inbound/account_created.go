package inbound

import (
	"context"
	"encoding/json"

	"github.com/netbill/profiles-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (i *Inbound) AccountCreated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload contracts.AccountCreatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	return i.modules.profile.Create(ctx, payload.AccountID, payload.Username)
}
