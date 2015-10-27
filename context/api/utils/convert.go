package utils

import "strconv"

func ConvertToInt64(intString string) int64 {
	intValue, err := strconv.ParseInt(intString, 10, 64)
	if err != nil {
		return 0
	}
	return intValue
}
