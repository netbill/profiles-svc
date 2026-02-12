package tokenmanager

import (
	"time"
)

type Manager struct {
	uploadSK string
	accessSK string

	profileMediaUploadTTL time.Duration
}

const (
	ProfilesServiceActor = "profiles-svc"
	ProfileResource      = "profile"
)

func New(accessSK, uploadSK string, profileMediaUploadTTL time.Duration) *Manager {
	return &Manager{
		accessSK:              accessSK,
		uploadSK:              uploadSK,
		profileMediaUploadTTL: profileMediaUploadTTL,
	}
}
