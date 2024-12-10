package types

import (
	"errors"
	"strconv"
	"strings"

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

func CheckValue(mType, mValue string) (any, error) {
	switch mType {
	case Gauge:
		value, err := strconv.ParseFloat(mValue, 64)
		return value, err
	case Counter:
		value, err := strconv.ParseInt(mValue, 10, 64)
		return value, err
	}
	return nil, errors.New("undefined metric type")
}

func CheckMetricType(mType string) error {
	if len(mType) == 0 {
		return errors.New(`empty metric type name`)
	}
	if len(mType) == 0 || (strings.ToLower(mType) != Gauge && strings.ToLower(mType) != Counter) {
		return errors.New(`the type of the metric is not defined`)
	}
	return nil
}
