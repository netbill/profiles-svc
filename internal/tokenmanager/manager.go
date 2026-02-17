package tokenmanager

import (
	"time"
)

const Issuer = "profiles-svc"

type Config struct {
	AccountAccess struct {
		SecretKey string `mapstructure:"secret_key" reqquire:"true"`
	} `mapstructure:"account_access" reqquire:"true"`
}

type Manager struct {
	issuer   string
	accessSK string

	mediaTTL time.Duration
}

func New(config Config) *Manager {
	return &Manager{
		issuer:   Issuer,
		accessSK: config.AccountAccess.SecretKey,
	}
}
