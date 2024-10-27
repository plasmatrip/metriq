package storage

import (
	"errors"
	"sync"

	"github.com/plasmatrip/metriq/internal/types"
)

type Metric struct {
	MetricType string
	Value      any
}

type MemStorage struct {
	Mu      sync.Mutex
	Storage map[string]Metric
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Mu:      sync.Mutex{},
		Storage: map[string]Metric{},
	}
}

func (ms *MemStorage) Update(key string, metric Metric) error {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()
	switch metric.MetricType {
	case types.Gauge:
		ms.Storage[key] = metric
		ms.updateCounter(types.PollCount, Metric{MetricType: types.Counter, Value: int64(1)})
	case types.Counter:
		ms.updateCounter(key, metric)
	}
	return nil
}

func (ms *MemStorage) updateCounter(key string, metric Metric) error {
	if oldMetric, ok := ms.Storage[key]; ok {
		oldValue, ok := oldMetric.Value.(int64)
		if !ok {
			return errors.New("failed to cast stored value to type int64")
		}
		newValue, ok := metric.Value.(int64)
		if !ok {
			return errors.New("failed to cast the received value to type int64")
		}
		ms.Storage[key] = Metric{MetricType: metric.MetricType, Value: oldValue + newValue}
		return nil
	}
	ms.Storage[key] = metric
	return nil
}

func (ms *MemStorage) Get(key string) (Metric, bool) {
	if metric, ok := ms.Storage[key]; ok {
		return metric, true
	}
	return Metric{}, false
}

func (ms *MemStorage) GetAll() map[string]Metric {
	return ms.Storage
}
