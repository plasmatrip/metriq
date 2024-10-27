package storage

import (
	"testing"

	"github.com/plasmatrip/metriq/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_Update(t *testing.T) {
	storage := NewStorage()

	tests := []struct {
		name       string
		key        string
		value      any
		want       any
		metricType string
	}{
		{
			name:       "Increment counter",
			key:        "key",
			value:      int64(1),
			want:       int64(1),
			metricType: types.Counter,
		},
		{
			name:       "Gauge counter",
			key:        "key",
			value:      float64(1),
			want:       float64(1),
			metricType: types.Gauge,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.metricType {
			case types.Counter:
				value := tt.value.(int64)
				storage.Update(tt.key, Metric{MetricType: tt.metricType, Value: value})
			case types.Gauge:
				value := tt.value.(float64)
				storage.Update(tt.key, Metric{MetricType: tt.metricType, Value: value})
			}
			metric, ok := storage.Get(tt.key)
			require.True(t, ok)
			assert.Equal(t, tt.want, metric.Value)
		})
	}
}
