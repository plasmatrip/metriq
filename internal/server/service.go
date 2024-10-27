package server

import (
	"errors"
	"strconv"

	"github.com/plasmatrip/metriq/internal/config"
)

func CheckType(mType string) bool {
	return mType == config.Gauge || mType == config.Counter
}

func CheckValue(mType, mValue string) error {
	switch mType {
	case config.Gauge:
		if _, err := strconv.ParseFloat(mValue, 64); err != nil {
			return err
		}

	case config.Counter:
		if _, err := strconv.ParseInt(mValue, 10, 64); err != nil {
			return err
		}
	}
	return nil
}

func CheckMetricType(metricType string) error {
	if len(metricType) == 0 || !CheckType(metricType) {
		return errors.New(`the type of the metric is not defined`)
	}
	return nil
}
