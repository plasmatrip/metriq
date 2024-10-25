package server

import (
	"strconv"
)

func CheckType(mType string) bool {
	return mType == Gauge || mType == Counter
}

func CheckName(mName string) bool {
	return true
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
