package storage

import "github.com/plasmatrip/metriq/internal/types"

type Repository interface {
	Update(key string, metric types.Metric) error
	Get(key string) (types.Metric, bool)
	GetAll() map[string]types.Metric
}
