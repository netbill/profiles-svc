package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountActor = uuid.UUID

type UploadScope = uuid.UUID

type Account struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	Version  int32     `json:"version"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
