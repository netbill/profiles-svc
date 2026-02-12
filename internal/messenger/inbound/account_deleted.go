package inbound

import (
	"context"
	"encoding/json"

	"github.com/netbill/profiles-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (i *Inbound) AccountDeleted(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload contracts.AccountDeletedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	if err := i.domain.Delete(ctx, payload.AccountID); err != nil {
		return err
	}

	return nil
}
