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
