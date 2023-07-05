package candy

import (
	"time"
)

// go1.20新增
const (
	DateOnly = "2006-01-02"
	TimeOnly = "15:04:05"
	DateTime = DateOnly + " " + TimeOnly
)

var (
	location *time.Location = time.Local
)

// 设置时区。(非必要不主动设置)
func SetDefaultLocation(name string) (err error) {
	location, err = time.LoadLocation(name)
	return
}

type Timer struct {
	t time.Time
}

// 当前时间Timer
func NowTimer() Timer {
	return Timer{t: time.Now().In(location)}
}

// 字符串转Timer
func StringTimer(s string) Timer {
	t, _ := time.ParseInLocation(DateTime, s, location)
	return Timer{t: t}
}

// 时间戳转Timer
func UnixTimer(i int64) Timer {
	t := time.Unix(i, 0)
	return Timer{t: t}
}

// Timer转Time
func (t Timer) Time() time.Time {
	return t.t
}

// Timer转字符串
func (t Timer) String() string {
	return t.t.Format(DateTime)
}

// Timer转时间戳
func (t Timer) Unix() int64 {
	return t.t.Unix()
}

// 加 秒/分/时
func (t Timer) Add(d time.Duration) Timer {
	t.t = t.t.Add(d)
	return t
}

// 加天数
func (t Timer) AddDay(i int) Timer {
	t.t = t.t.AddDate(0, 0, i)
	return t
}

// 加月份
func (t Timer) AddMonth(i int) Timer {
	t.t = t.t.AddDate(0, i, 0)
	return t
}

// 加年
func (t Timer) AddYear(i int) Timer {
	t.t = t.t.AddDate(i, 0, 0)
	return t
}

// 当日的开始时间 (00:00:00)
func (t Timer) DayStart() Timer {
	y, m, d := t.t.Date()
	t.t = time.Date(y, m, d, 0, 0, 0, 0, location)
	return t
}

// 当日的结束时间 (23:59:59)
func (t Timer) DayEnd() Timer {
	return t.DayStart().Add(time.Hour*24 - time.Nanosecond)
}

// 当月的开始时间 (00:00:00)
func (t Timer) MonthStart() Timer {
	y, m, _ := t.t.Date()
	t.t = time.Date(y, m, 1, 0, 0, 0, 0, location)
	return t
}

// 当月的结束时间 (23:59:59)
func (t Timer) MonthEnd() Timer {
	return t.MonthStart().AddMonth(1).Add(-time.Nanosecond)
}
