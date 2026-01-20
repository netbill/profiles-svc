package contracts

import (
	"time"

	"github.com/google/uuid"
)

const ProfileCreatedEvent = "profile.created"

type AccountProfileCreatedPayload struct {
	Data      AccountProfileCreatedPayloadData `json:"data"`
	Timestamp time.Time                        `json:"timestamp"`
}

type AccountProfileCreatedPayloadData struct {
	AccountID   uuid.UUID `json:"account_id"`
	Username    string    `json:"username"`
	Official    bool      `json:"official"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
}

const ProfileUpdatedEvent = "profile.updated"

type AccountProfileUpdatedPayload struct {
	Data      AccountProfileUpdatedPayloadData `json:"data"`
	Timestamp time.Time                        `json:"timestamp"`
}

type AccountProfileUpdatedPayloadData struct {
	AccountID   uuid.UUID `json:"account_id"`
	Username    string    `json:"username"`
	Official    bool      `json:"official"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
}

const ProfileUsernameUpdatedEvent = "profile.username.updated"

type AccountProfileUsernameUpdatedPayload struct {
	Data      AccountProfileUsernameUpdatedPayloadData `json:"data"`
	Timestamp time.Time                                `json:"timestamp"`
}

type AccountProfileUsernameUpdatedPayloadData struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
}

const ProfileOfficialUpdatedEvent = "profile.official.updated"

type AccountProfileOfficialUpdatedPayload struct {
	Data      AccountProfileOfficialUpdatedPayloadData `json:"data"`
	Timestamp time.Time                                `json:"timestamp"`
}

type AccountProfileOfficialUpdatedPayloadData struct {
	AccountID uuid.UUID `json:"account_id"`
	Official  bool      `json:"official"`
}
