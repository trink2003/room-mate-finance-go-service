package utils

import (
	"fmt"
	"time"
)

func BeginningOfMonth(date time.Time) time.Time {
	result := date.AddDate(0, 0, -date.Day()+1)
	fmt.Printf("BeginningOfMonth time zone is %s\n", result.Location())
	return time.Date(result.Year(), result.Month(), result.Day(), 0, 0, 0, result.Nanosecond(), result.Location())
}

func EndOfMonth(date time.Time) time.Time {
	result := date.AddDate(0, 1, -date.Day())
	fmt.Printf("EndOfMonth time zone is %s\n", result.Location().String())
	return time.Date(result.Year(), result.Month(), result.Day(), 23, 59, 59, result.Nanosecond(), result.Location())
}
