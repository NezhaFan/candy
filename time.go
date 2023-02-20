package util

import (
	"time"
)

// go1.20新增
const (
	DateTime = "2006-01-02 15:04:05"
	DateOnly = "2006-01-02"
	TimeOnly = "15:04:05"
)

func Now() string {
	return time.Now().Format(DateTime)
}

func StringToTime(s string) time.Time {
	t, _ := time.ParseInLocation(DateTime, s, time.Local)
	return t
}

// 相对日期开始。   例：昨日凌点 DayStart(-1)
func DayStart(add int) time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d+add, 0, 0, 0, 0, time.Local)
}

// 相对日期结束。   例：今日结束 DayStart(0)
func DayEnd(add int) time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d+add, 23, 59, 59, 0, time.Local)
}

// 相对月初
func MonthStart(add int) time.Time {
	y, m, _ := time.Now().Date()
	return time.Date(y, m+time.Month(add), 1, 0, 0, 0, 0, time.Local)
}

// 相对月末
func MonthEnd(add int) time.Time {
	y, m, _ := time.Now().Date()
	return time.Date(y, m+time.Month(add)+1, 1, 0, 0, 0, 0, time.Local).Add(-1 * time.Second)
}
