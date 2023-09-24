package utils

import (
	"time"

	"github.com/shopspring/decimal"
)

func DaysToDuration(days int64) time.Duration {
	return time.Duration(days) * time.Hour * 24
}

func DurationToDays(d time.Duration) decimal.Decimal {
	return decimal.New(int64(d), 0).Div(decimal.New(int64(time.Hour*24), 0))
}

func FormatToDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func MonthsBack(from time.Time, months int) time.Time {
	return from.AddDate(0, -months, 0)
}
