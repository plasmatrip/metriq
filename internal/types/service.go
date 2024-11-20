package types

import (
	"errors"
	"strconv"
	"strings"
)

func checkType(mType string) bool {
	return strings.ToLower(mType) == Gauge || strings.ToLower(mType) == Counter
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

func CheckMetricType(metricType string) error {
	if len(metricType) == 0 || !checkType(metricType) {
		return errors.New(`the type of the metric is not defined`)
	}
	return nil
}
