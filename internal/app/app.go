package app

import (
	"github.com/netbill/profiles-svc/internal/config"
	"github.com/netbill/profiles-svc/pkg/log"
)

type App struct {
	log    *log.Logger
	config *config.Config
}

func New(log *log.Logger, cfg *config.Config) *App {
	return &App{
		log:    log,
		config: cfg,
	}
}
