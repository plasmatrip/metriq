package handlers

import (
	"github.com/plasmatrip/metriq/internal/server/config"
	"github.com/plasmatrip/metriq/internal/storage"
)

type Handlers struct {
	Repo   storage.Repository
	config config.Config
}

func NewHandlers(repo storage.Repository, config config.Config) *Handlers {
	return &Handlers{Repo: repo, config: config}
}
