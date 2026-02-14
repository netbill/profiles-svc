package inbound

import (
	"context"
	"encoding/json"

	"github.com/netbill/profiles-svc/internal/core/modules/account"
	"github.com/netbill/profiles-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

func (i *Inbound) AccountUsernameUpdated(
	ctx context.Context,
	message kafka.Message,
) error {
	var payload evtypes.AccountUsernameUpdatedPayload
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		return err
	}

	return i.modules.profile.UpdateUsername(ctx, payload.AccountID, account.UpdateUsernameParams{
		Username:  payload.NewUsername,
		Version:   payload.Version,
		UpdatedAt: payload.UpdatedAt,
	})
}
