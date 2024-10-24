package agent

import (
	"fmt"

	"github.com/plasmatrip/metriq/internal/storage"
)

type Controller struct {
	Repo storage.Repository
}

func NewSender(repo storage.Repository) *Controller {
	return &Controller{Repo: repo}
}

func (s *Controller) SendMetrics() {
	fmt.Println("SendMetrics")
}

func (s *Controller) UpdateMetrics() {
	fmt.Println("UpdateMetrics")
}
