package storage

import (
	"sync"
)

const (
	PollCount = "PollCount"
)

type (
	Gauge   float64
	Counter int64
)

type Repository interface {
	UpdateGauge(key string, value Gauge)
	UpdateCounter(value int64)
	GetGauges() map[string]Gauge
	GetGauge(key string) Gauge
	GetCounter() int64
	// Print()
}

type MemStorage struct {
	Mu        sync.Mutex
	Storage   map[string]Gauge
	PollCount Counter
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Mu:        sync.Mutex{},
		Storage:   make(map[string]Gauge),
		PollCount: 0,
	}
}

func (ms *MemStorage) UpdateGauge(key string, value Gauge) {
	ms.Mu.Lock()
	defer func() {
		ms.Mu.Unlock()
		ms.UpdateCounter(1)
	}()
	ms.Storage[key] = value
}

func (ms *MemStorage) UpdateCounter(value int64) {
	ms.Mu.Lock()
	defer ms.Mu.Unlock()
	ms.PollCount += Counter(value)
}

func (ms *MemStorage) GetGauges() map[string]Gauge {
	return ms.Storage
}

func (ms *MemStorage) GetGauge(key string) Gauge {
	return ms.Storage[key]
}

func (ms *MemStorage) GetCounter() int64 {
	return int64(ms.PollCount)
}

// func (ms *MemStorage) Print() {
// 	fmt.Println("====================")
// 	for k, v := range ms.storage {
// 		fmt.Printf("%s update: %s = %v\r\n", time.Now().Format("15:04:05"), k, v)
// 	}
// }
