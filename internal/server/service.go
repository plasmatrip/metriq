package server

import (
	"strconv"

	"github.com/plasmatrip/metriq/internal/storage"
)

func CheckType(mType string) bool {
	return mType == storage.Gauge || mType == storage.Counter
}

func CheckName(mName string) bool {
	return true
}

func CheckValue(mType, mValue string) error {
	switch mType {
	case storage.Gauge:
		if _, err := strconv.ParseFloat(mValue, 64); err != nil {
			return err
		}

	case storage.Counter:
		if _, err := strconv.ParseInt(mValue, 10, 64); err != nil {
			return err
		}
	}
	return nil
}
