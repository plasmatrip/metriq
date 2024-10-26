package server

import (
	"errors"
	"strconv"

	"github.com/plasmatrip/metriq/internal/storage"
)

func CheckType(mType string) bool {
	return mType == Gauge || mType == Counter
}

func CheckName(mName string) bool {
	_, ok := storage.Metrics[mName]
	return ok
}

func CheckValue(mType, mValue string) error {
	switch mType {
	case Gauge:
		if _, err := strconv.ParseFloat(mValue, 64); err != nil {
			return err
		}

	case Counter:
		if _, err := strconv.ParseInt(mValue, 10, 64); err != nil {
			return err
		}
	}
	return nil
}

func CheckMetricName(name string) error {
	// if len(name) == 0 || !CheckName(name) {
	if len(name) == 0 {
		return errors.New(`the name of the metric is not defined`)
	}
	return nil
}

func CheckMetricType(metricType string) error {
	if len(metricType) == 0 || !CheckType(metricType) {
		return errors.New(`the type of the metric is not defined`)
	}
	return nil
}
