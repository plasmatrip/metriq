package types

import (
	"errors"
	"strconv"
	"strings"
)

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
	if len(mType) == 0 || (strings.ToLower(mType) != Gauge && strings.ToLower(mType) != Counter) {
		return errors.New(`the type of the metric is not defined`)
	}
	return nil
}
