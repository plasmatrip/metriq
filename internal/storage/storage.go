package storage

import (
	"fmt"
	"sync"
	"time"
)

const (
	PollCount = "PollCount"
)

type (
	Gauge   float64
	counter int64
)

type Repository interface {
	Update(key string, value any)
	GetAll() map[string]any
	Print()
}

// type metric struct {
// 	mType string
// 	value any
// }

type MemStorage struct {
	mu      sync.Mutex
	storage map[string]any
}

func NewStorage() *MemStorage {
	return &MemStorage{
		storage: make(map[string]any),
		mu:      sync.Mutex{},
	}
}

func (ms *MemStorage) Update(key string, value any) {
	ms.mu.Lock()
	defer func() {
		ms.mu.Unlock()
		ms.updateCounter()
	}()
	ms.storage[key] = value
}

func (ms *MemStorage) updateCounter() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	count, ok := ms.storage[PollCount]
	if !ok {
		ms.storage[PollCount] = counter(0)
		return
	}
	ms.storage[PollCount] = count.(counter) + 1
}

func (ms *MemStorage) GetAll() map[string]any {
	return ms.storage
}

func (ms *MemStorage) Print() {
	fmt.Println("====================")
	for k, v := range ms.storage {
		fmt.Printf("%s update: %s = %v\r\n", time.Now().Format("15:04:05"), k, v)
	}
}
