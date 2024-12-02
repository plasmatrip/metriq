package storage

import (
	"context"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

type Repository interface {
	SetMetrics(context.Context, []models.Metrics) error
	// SetMetrics(context.Context, models.SMetrics) error
	SetMetric(mName string, metric types.Metric) error
	Metric(mName string) (types.Metric, error)
	Metrics() (map[string]types.Metric, error)
	SetBackup(chan struct{})
	Ping(ctx context.Context) error
	Close() error
}
