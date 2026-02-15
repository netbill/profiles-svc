package tokenmanager

import (
	"time"
)

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

func New(issuer string, config Config) *Manager {
	return &Manager{
		issuer:   issuer,
		accessSK: config.AccountAccess.SecretKey,
	}
}
