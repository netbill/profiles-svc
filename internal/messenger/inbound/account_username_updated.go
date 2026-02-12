package inbound

import (
	"context"
	"encoding/json"

	"github.com/netbill/profiles-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (i *Inbound) AccountUsernameUpdated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload contracts.AccountUsernameUpdatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	if _, err := i.domain.UpdateUsername(ctx, payload.AccountID, payload.NewUsername); err != nil {
		return err
	}

	return nil
}
