package storage

import (
	"context"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

type Repository interface {
	SetMetrics(ctx context.Context, metrics []models.Metrics) error
	SetMetric(ctx context.Context, mName string, metric types.Metric) error
	Metric(ctx context.Context, mName string) (types.Metric, error)
	Metrics(context.Context) (map[string]types.Metric, error)
	SetBackup(chan struct{})
	Ping(context.Context) error
	Close()
}
