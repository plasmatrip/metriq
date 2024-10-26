package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_UpdateCounter(t *testing.T) {
	storage := NewStorage()

	tests := []struct {
		name  string
		value int64
		want  int64
	}{
		{
			name:  "Increment counter",
			value: 1,
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.UpdateCounter(tt.value)
			assert.Equal(t, tt.want, storage.GetCounter())
		})
	}
}

func TestStorage_UpdateGauge(t *testing.T) {
	storage := NewStorage()

	tests := []struct {
		name  string
		key   string
		value Gauge
		want  Gauge
	}{
		{
			name:  "Update metric",
			key:   "key",
			value: 1,
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.UpdateGauge(tt.key, tt.value)
			assert.Equal(t, tt.want, storage.GetGauge(tt.key))
		})
	}
}
