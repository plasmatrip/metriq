package handlers

import (
	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/storage"
)

type Handlers struct {
	Repo   storage.Repository
	config config.Config
	lg     logger.Logger
}

func NewHandlers(repo storage.Repository, config config.Config, lg logger.Logger) *Handlers {
	return &Handlers{Repo: repo, config: config, lg: lg}
}
