package contracts

import (
	"time"

	"github.com/google/uuid"
)

const AccountsTopicV1 = "accounts.v1"

const AccountDeletedEvent = "account.deleted"

type AccountDeletedPayload struct {
	Data      AccountDeletedPayloadData `json:"data"`
	Timestamp time.Time                 `json:"timestamp"`
}

type AccountDeletedPayloadData struct {
	AccountID uuid.UUID `json:"account_id"`
}
