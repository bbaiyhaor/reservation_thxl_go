package utils

import (
	"time"
)

var weekdayMap = map[int]string{
	1: "星期一",
	2: "星期二",
	3: "星期三",
	4: "星期四",
	5: "星期五",
	6: "星期六",
	0: "星期日",
}

func GetChineseWeekday(ts time.Time) string {
	return weekdayMap[int(ts.Weekday())]
}

func ConcatTime(date time.Time, clock time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), clock.Hour(), clock.Minute(),
		clock.Second(), clock.Nanosecond(), time.Local)
}

func BeginOfDay(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
}

func BeginOfYesterday(tm time.Time) time.Time {
	return BeginOfDay(tm.Add(-24 * time.Hour))
}

func BeginOfTomorrow(tm time.Time) time.Time {
	return BeginOfDay(tm.Add(24 * time.Hour))
}
