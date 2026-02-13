package tokenmanager

import (
	"time"
)

type Manager struct {
	issuer string

	uploadSK string
	accessSK string

	profileMediaUploadTTL time.Duration
}

const (
	ProfileResource = "profile"
)

type Config struct {
	AccessSK string
	UploadSK string

	ProfileMediaUploadTTL time.Duration
}

func New(issuer string, config Config) *Manager {
	return &Manager{
		issuer:                issuer,
		accessSK:              config.AccessSK,
		uploadSK:              config.UploadSK,
		profileMediaUploadTTL: config.ProfileMediaUploadTTL,
	}
}
