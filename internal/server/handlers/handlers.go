package handlers

import (
	"github.com/plasmatrip/metriq/internal/storage"
)

type Handlers struct {
	Repo storage.Repository
}

func NewHandlers(repo storage.Repository) *Handlers {
	return &Handlers{Repo: repo}
}
