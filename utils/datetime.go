package utils

import (
	"time"
)

func UTCDate(year int, month time.Month, Day int) time.Time {
	return time.Date(year, month, Day, 0, 0, 0, 0, time.UTC)
}

func UTCDateFromTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func UTCStartOfMonth() time.Time {
	now := time.Now().UTC()
	year, month, _ := now.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}
func DaysSince(ts *time.Time) int {
	if ts == nil {
		return 0
	}
	return int(time.Now().UTC().Sub(*ts).Hours() / 24)
}
