package storage

import (
	"context"

	"github.com/plasmatrip/metriq/internal/types"
)

type Repository interface {
	SetMetric(mName string, metric types.Metric) error
	Metric(mName string) (types.Metric, error)
	Metrics() (map[string]types.Metric, error)
	SetBackup(chan struct{})
	Ping(ctx context.Context) error
	Close() error
}
