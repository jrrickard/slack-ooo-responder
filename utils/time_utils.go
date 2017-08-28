package utils

import (
	"strconv"
	"time"
)

func ConvertTimestamp(timestamp string) (time.Time, error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	timeObj := time.Unix(i, 0)
	return timeObj, err
}
