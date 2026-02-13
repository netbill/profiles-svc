package tokenmanager

import (
	"time"
)

type Manager struct {
	Issuer string

	UploadSK string
	AccessSK string

	ProfileMediaUploadTTL time.Duration
}

const (
	ProfileResource = "profile"
)
