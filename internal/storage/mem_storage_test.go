package storage

import (
	"testing"

	"github.com/plasmatrip/metriq/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_Update(t *testing.T) {
	storage := NewStorage()

	tests := []struct {
		name       string
		key        string
		value      any
		want       any
		errWant    bool
		metricType string
	}{
		{
			name:       "Increment counter correct value",
			key:        "key",
			value:      int64(1),
			want:       int64(1),
			errWant:    false,
			metricType: types.Counter,
		},
		{
			name:       "Gauge counter correct value",
			key:        "key",
			value:      float64(1),
			want:       float64(1),
			errWant:    false,
			metricType: types.Gauge,
		},
		{
			name:       "Increment counter incorrect value",
			key:        "key",
			value:      float64(1),
			errWant:    true,
			metricType: types.Counter,
		},
		{
			name:       "Gauge counter incorrect value",
			key:        "key",
			value:      "a",
			errWant:    true,
			metricType: types.Gauge,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			switch test.metricType {
			case types.Counter:
				err = storage.SetMetric(test.key, types.Metric{MetricType: test.metricType, Value: test.value})
			case types.Gauge:
				err = storage.SetMetric(test.key, types.Metric{MetricType: test.metricType, Value: test.value})
			}
			if test.errWant {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestMemStorage_Get(t *testing.T) {
	storage := NewStorage()

	tests := []struct {
		name       string
		key        string
		metric     types.Metric
		getKey     string
		want       any
		errWant    bool
		metricType string
	}{
		{
			name:       "Increment counter correct value",
			key:        "Counter",
			metric:     types.Metric{MetricType: types.Counter, Value: int64(1)},
			getKey:     "Counter",
			want:       int64(1),
			errWant:    false,
			metricType: types.Counter,
		},
		{
			name:       "Gauge counter correct value",
			key:        "Gouge",
			metric:     types.Metric{MetricType: types.Gauge, Value: float64(1)},
			getKey:     "Counter",
			want:       float64(1),
			errWant:    false,
			metricType: types.Gauge,
		},
		{
			name:       "Increment counter incorrect name",
			key:        "Counter",
			metric:     types.Metric{MetricType: types.Counter, Value: int64(1)},
			getKey:     "aa",
			errWant:    true,
			metricType: types.Counter,
		},
		{
			name:       "Gauge counter incorrect name",
			key:        "Gouge",
			metric:     types.Metric{MetricType: types.Gauge, Value: float64(1)},
			getKey:     "aa",
			errWant:    true,
			metricType: types.Gauge,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.metricType {
			case types.Counter:
				_ = storage.SetMetric(test.key, test.metric)
			case types.Gauge:
				_ = storage.SetMetric(test.key, test.metric)
			}
			_, ok := storage.Metric(test.getKey)
			if test.errWant {
				assert.False(t, ok)
				return
			}
			assert.True(t, ok)
		})
	}
}
