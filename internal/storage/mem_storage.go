package storage

import (
	"errors"
	"maps"
	"sync"

	"github.com/plasmatrip/metriq/internal/types"
)

type Metric struct {
	MetricType string
	Value      any
}

type storage map[string]Metric

type MemStorage struct {
	Mu      sync.RWMutex
	Storage storage //map[string]Metric
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Mu:      sync.RWMutex{},
		Storage: make(storage), //map[string]Metric{},
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
	ms.Mu.RLock()
	defer ms.Mu.RUnlock()
	if metric, ok := ms.Storage[key]; ok {
		return metric, true
	}
	return Metric{}, false
}

func (ms *MemStorage) GetAll() map[string]Metric {
	ms.Mu.RLock()
	defer ms.Mu.RUnlock()
	copyStorage := make(storage, len(ms.Storage))
	maps.Copy(copyStorage, ms.Storage)
	return copyStorage //ms.Storage
}
