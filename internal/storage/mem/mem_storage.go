package mem

import (
	"context"
	"errors"
	"maps"
	"sync"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

type storage map[string]types.Metric

type MemStorage struct {
	Mu      sync.RWMutex
	Storage storage
	bkp     backup
}

type backup struct {
	do bool
	c  chan struct{}
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Mu:      sync.RWMutex{},
		Storage: make(storage),
		bkp:     backup{do: false, c: nil},
	}
}

func (ms *MemStorage) Ping(_ context.Context) error {
	return nil
}

func (ms *MemStorage) Close() {
}

func (ms *MemStorage) SetMetrics(ctx context.Context, metrics []models.Metrics) error {
	for _, metric := range metrics {
		switch metric.MType {
		case types.Gauge:
			err := ms.SetMetric(ctx, metric.ID, types.Metric{MetricType: metric.MType, Value: *metric.Value})
			if err != nil {
				return err
			}
		case types.Counter:
			err := ms.SetMetric(ctx, metric.ID, types.Metric{MetricType: metric.MType, Value: *metric.Delta})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ms *MemStorage) SetMetric(ctx context.Context, mName string, metric types.Metric) error {
	ms.Mu.Lock()
	switch metric.MetricType {
	case types.Gauge:
		if err := metric.Check(); err != nil {
			ms.Mu.Unlock()
			return err
		}
		ms.Storage[mName] = metric
		err := ms.setCounter(ctx, types.PollCount, types.Metric{MetricType: types.Counter, Value: int64(1)})
		if err != nil {
			ms.Mu.Unlock()
			return err
		}
	case types.Counter:
		err := ms.setCounter(ctx, mName, metric)
		if err != nil {
			ms.Mu.Unlock()
			return err
		}
	}

	ms.Mu.Unlock()

	if ms.bkp.do {
		ms.bkp.c <- struct{}{}
		select {
		case <-ctx.Done():
			close(ms.bkp.c)
		default:
		}
	}

	return nil
}

func (ms *MemStorage) setCounter(_ context.Context, mName string, metric types.Metric) error {
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

func (ms *MemStorage) SetBackup(c chan struct{}) {
	ms.bkp.do = true
	ms.bkp.c = c
}

func (ms *MemStorage) Metric(_ context.Context, key string) (types.Metric, error) {
	ms.Mu.RLock()
	defer ms.Mu.RUnlock()
	metric, ok := ms.Storage[key]
	if !ok {
		return types.Metric{}, errors.New("metric not found")
	}
	return metric, nil
}

func (ms *MemStorage) Metrics(_ context.Context) (map[string]types.Metric, error) {
	ms.Mu.RLock()
	defer ms.Mu.RUnlock()
	copyStorage := make(storage, len(ms.Storage))
	maps.Copy(copyStorage, ms.Storage)
	return copyStorage, nil
}
