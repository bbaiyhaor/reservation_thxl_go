package utils

import (
	"fmt"
	"time"
)

var Weekdays = [...]string{
	"日",
	"一",
	"二",
	"三",
	"四",
	"五",
	"六",
}

const (
	TIME_PATTERN  = "2006-01-02 15:04"
	DATE_PATTERN  = "2006-01-02"
	CLOCK_PATTERN = "15:04"
)

var (
	Location *time.Location
)

func ConcatTime(date time.Time, clock time.Time) time.Time {
	return time.Date(date.In(Location).Year(), date.In(Location).Month(), date.In(Location).Day(),
		clock.In(Location).Hour(), clock.In(Location).Minute(), clock.In(Location).Second(),
		clock.In(Location).Nanosecond(), Location)
}

func GetToday() time.Time {
	now := time.Now().In(Location)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, Location)
}

func GetNow() time.Time {
	return time.Now().In(Location)
}

func ParseClock(clock string) (time.Time, error) {
	return time.ParseInLocation(TIME_PATTERN, fmt.Sprintf("%s %s", DATE_PATTERN, clock), Location)
}
