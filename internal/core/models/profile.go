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
	AvatarURL   *string   `json:"avatar_url,omitempty"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (e Profile) IsNil() bool {
	return e.AccountID == uuid.Nil
}

type UpdateProfileAvatar struct {
	UploadURL   string `json:"upload_url"`
	GetURL      string `json:"get_url"`
	UploadToken string `json:"upload_token"`
}

type ProfileAvatar struct {
}
