package storage

import (
	"errors"
	"maps"
	"sync"

	"github.com/plasmatrip/metriq/internal/types"
)

type storage map[string]types.Metric

type MemStorage struct {
	Mu      sync.RWMutex
	Storage storage
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Mu:      sync.RWMutex{},
		Storage: make(storage),
	}
}

func (ms *MemStorage) Update(mName string, metric types.Metric) error {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()
	var err error
	switch metric.MetricType {
	case types.Gauge:
		if err := metric.Check(); err != nil {
			return err
		}
		ms.Storage[mName] = metric
		err = ms.updateCounter(types.PollCount, types.Metric{MetricType: types.Counter, Value: int64(1)})
	case types.Counter:
		err = ms.updateCounter(mName, metric)
	}
	return err
}

func (ms *MemStorage) updateCounter(mName string, metric types.Metric) error {
	if oldMetric, ok := ms.Storage[mName]; ok {
		oldValue, ok := oldMetric.Value.(int64)
		if !ok {
			return errors.New("failed to cast stored value to type int64")
		}
		newValue, ok := metric.Value.(int64)
		if !ok {
			return errors.New("failed to cast the received value to type int64")
		}
		ms.Storage[mName] = types.Metric{MetricType: metric.MetricType, Value: oldValue + newValue}
		return nil
	}
	ms.Storage[mName] = metric
	return nil
}

func (ms *MemStorage) Get(key string) (types.Metric, bool) {
	ms.Mu.RLock()
	defer ms.Mu.RUnlock()
	metric, ok := ms.Storage[key]
	return metric, ok
}

func (ms *MemStorage) GetAll() map[string]types.Metric {
	ms.Mu.RLock()
	defer ms.Mu.RUnlock()
	copyStorage := make(storage, len(ms.Storage))
	maps.Copy(copyStorage, ms.Storage)
	return copyStorage
}
