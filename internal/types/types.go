package types

import (
	"errors"

	"github.com/plasmatrip/metriq/internal/models"
)

const (
	Gauge   = "gauge"
	Counter = "counter"

	PollCount = "PollCount"
)

type Metric struct {
	MetricType string
	Value      any
}

func (metric Metric) Check() error {
	err := CheckMetricType(metric.MetricType)
	if err != nil {
		return err
	}
	switch metric.MetricType {
	case Gauge:
		_, ok := metric.Value.(float64)
		if !ok {
			return errors.New("the value is not float64")
		}
	case Counter:
		_, ok := metric.Value.(int64)
		if !ok {
			return errors.New("the value is not int64")
		}
	}
	return nil
}

func (metric Metric) Convert(key string) models.Metrics {
	if metric.MetricType == Gauge {
		value, _ := metric.Value.(float64)
		return models.Metrics{
			ID:    key,
			MType: metric.MetricType,
			Value: &value,
		}
	}
	value, _ := metric.Value.(int64)
	return models.Metrics{
		ID:    key,
		MType: metric.MetricType,
		Delta: &value,
	}
}
