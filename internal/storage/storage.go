package storage

import "github.com/plasmatrip/metriq/internal/types"

type Repository interface {
	SetMetric(key string, metric types.Metric) error
	Metric(key string) (types.Metric, bool)
	Metrics() map[string]types.Metric
	SetBackup(chan bool)
}
