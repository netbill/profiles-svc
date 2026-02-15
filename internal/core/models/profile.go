package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	AccountID   uuid.UUID `json:"account_id"`
	Username    string    `json:"username"`
	Official    bool      `json:"official"`
	Pseudonym   *string   `json:"pseudonym,omitempty"`
	Description *string   `json:"description,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
	Version     int32     `json:"version"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UploadMediaLink struct {
	Key        string `json:"key"`
	UploadURL  string `json:"upload_url"`
	PreloadUrl string `json:"preload_url"`
}
