package utils

import (
	"time"
)

func BeginningOfDay(cTime time.Time) time.Time {
	year, month, day := cTime.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, cTime.Location())
}

func EndOfDay(cTime time.Time) time.Time {
	year, month, day := cTime.Date()
	return time.Date(year, month, day, 23, 59, 59, 59, cTime.Location())
}
