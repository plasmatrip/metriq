package storage

import "sync"

type MemStorage struct {
	Mu             sync.Mutex
	GaugeStorage   map[string]Gauge
	CounterStorage map[string]Counter
	// PollCount Counter
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Mu:             sync.Mutex{},
		GaugeStorage:   make(map[string]Gauge),
		CounterStorage: make(map[string]Counter),
		// PollCount:    0,
	}
}

func (ms *MemStorage) UpdateGauge(key string, value Gauge) {
	ms.Mu.Lock()
	defer func() {
		ms.Mu.Unlock()
		ms.UpdateCounter(PollCount, 1)
	}()
	ms.GaugeStorage[key] = value
}

func (ms *MemStorage) UpdateCounter(key string, value Counter) {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()
	if v, ok := ms.CounterStorage[key]; ok {
		ms.CounterStorage[key] = v + value
		return
	}
	ms.CounterStorage[key] = value
}

func (ms *MemStorage) GetGauges() map[string]Gauge {
	return ms.GaugeStorage
}

func (ms *MemStorage) GetCounters() map[string]Counter {
	return ms.CounterStorage
}

func (ms *MemStorage) GetGauge(key string) Gauge {
	return ms.GaugeStorage[key]
}

func (ms *MemStorage) GetCounter(key string) Counter {
	return ms.CounterStorage[key]
}
