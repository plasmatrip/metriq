package server

import (
	"errors"
	"strconv"

	"github.com/plasmatrip/metriq/internal/types"
)

func CheckType(mType string) bool {
	return mType == types.Gauge || mType == types.Counter
}

func CheckValue(mType, mValue string) error {
	switch mType {
	case types.Gauge:
		if _, err := strconv.ParseFloat(mValue, 64); err != nil {
			return err
		}

	case types.Counter:
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
