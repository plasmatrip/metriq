package types

import (
	"testing"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTypes_Check(t *testing.T) {
	tests := []struct {
		name    string
		value   Metric
		wantErr bool
	}{
		{
			name: "Valid metric type gauge",
			value: Metric{
				MetricType: Gauge,
				Value:      float64(1),
			},
			wantErr: false,
		},
		{
			name: "Valid metric type counter",
			value: Metric{
				MetricType: Counter,
				Value:      int64(1),
			},
			wantErr: false,
		},
		{
			name: "Valid metric value gauge",
			value: Metric{
				MetricType: Gauge,
				Value:      float64(1),
			},
			wantErr: false,
		},
		{
			name: "Valid metric value counter",
			value: Metric{
				MetricType: Counter,
				Value:      int64(1),
			},
			wantErr: false,
		},
		{
			name: "Invalid metric type",
			value: Metric{
				MetricType: "SomeMetric",
				Value:      int64(1),
			},
			wantErr: true,
		},
		{
			name: "Invalid metric value gauge",
			value: Metric{
				MetricType: Gauge,
				Value:      int64(1),
			},
			wantErr: true,
		},
		{
			name: "Invalid metric value counter",
			value: Metric{
				MetricType: Counter,
				Value:      float64(1),
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.value.Check()
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestTypes_Convert(t *testing.T) {
	floatValue := float64(1)
	intValue := int64(1)
	tests := []struct {
		name   string
		mName  string
		metric Metric
		want   models.Metrics
	}{
		{
			name:  "Valid gauge convert",
			mName: "SomeMetricName",
			metric: Metric{
				MetricType: Gauge,
				Value:      floatValue,
			},
			want: models.Metrics{
				ID:    "SomeMetricName",
				MType: Gauge,
				Value: &floatValue,
			},
		},
		{
			name:  "Valid counter convert",
			mName: "SomeMetricName",
			metric: Metric{
				MetricType: Counter,
				Value:      intValue,
			},
			want: models.Metrics{
				ID:    "SomeMetricName",
				MType: Counter,
				Delta: &intValue,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.metric.Convert(test.mName))
		})
	}
}

func TestTypes_CheckValue(t *testing.T) {
	t.Run("Float64 value", func(t *testing.T) {
		_, err := CheckValue(Gauge, "100")
		assert.NoError(t, err)
	})
	t.Run("Int64 value", func(t *testing.T) {
		_, err := CheckValue(Counter, "100")
		assert.NoError(t, err)
	})
	t.Run("Wrong value", func(t *testing.T) {
		_, err := CheckValue(Counter, "100.5")
		assert.Error(t, err)
	})
	t.Run("Wrong value", func(t *testing.T) {
		_, err := CheckValue(Gauge, "aa")
		assert.Error(t, err)
	})
	t.Run("Wrong type", func(t *testing.T) {
		_, err := CheckValue("SomeType", "aa")
		assert.Error(t, err)
	})
}

func TestTypes_CheckMetricType(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
		wanrErr bool
	}{
		{
			name:    "Gauge type test",
			value:   Gauge,
			wantErr: false,
		},
		{
			name:    "Counter type test",
			value:   Counter,
			wantErr: false,
		},
		{
			name:    "Wrong type test",
			value:   "SomeType",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckMetricType(test.value)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
