package storage

import (
	"context"

	"github.com/plasmatrip/metriq/internal/types"
)

type Repository interface {
	SetMetric(key string, metric types.Metric) error
	Metric(key string) (types.Metric, bool)
	Metrics() map[string]types.Metric
	SetBackup(chan struct{})
	Ping(ctx context.Context) error
	Close() error
}
