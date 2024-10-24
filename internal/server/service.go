package server

import "strconv"

func CheckType(mType string) bool {
	return mType == gauge || mType == counter
}

func CheckName(mName string) bool {
	return true
}

func CheckValue(mType, mValue string) error {
	switch mType {
	case gauge:
		if _, err := strconv.ParseFloat(mValue, 64); err != nil {
			return err
		}

	case counter:
		if _, err := strconv.ParseInt(mValue, 36, 64); err != nil {
			return err
		}
	}
	return nil
}
