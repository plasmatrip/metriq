package handlers

import (
	"github.com/plasmatrip/metriq/internal/server"
	"github.com/plasmatrip/metriq/internal/storage"
)

type Handlers struct {
	Repo   storage.Repository
	config server.Config
}

func NewHandlers(repo storage.Repository, config server.Config) *Handlers {
	return &Handlers{Repo: repo, config: config}
}
