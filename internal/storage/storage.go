package storage

import "fmt"

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Repository interface {
	// Add(key string, value any)
	// Delete(key string)
	Update(mType string, key string, value int64)
}

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (ms *MemStorage) Update(mType string, key string, value int64) {
	fmt.Println(mType)
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Gauge:   map[string]float64{},
		Counter: map[string]int64{},
	}
}
