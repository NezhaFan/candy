package candy

import (
	"time"
)

// go1.20新增
const (
	DateOnly  = "2006-01-02"
	TimeOnly  = "15:04:05"
	DateTime  = DateOnly + " " + TimeOnly
	DaySecond = 86400
)

var (
	location *time.Location = time.Local
)

func SetDefaultLocation(name string) (err error) {
	location, err = time.LoadLocation(name)
	return
}

type Timer struct {
	t time.Time
}

func StringTime(s string) Timer {
	t, _ := time.ParseInLocation(DateTime, s, location)
	return Timer{t: t}
}

func UnixTime[T Unsigned](i T) Timer {
	t := time.Unix(int64(i), 0).In(location)
	return Timer{t: t}
}

func Now() Timer {
	return Timer{t: time.Now().In(location)}
}

func (t Timer) Time() time.Time {
	return t.t
}

func (t Timer) String() string {
	return t.t.Format(DateTime)
}

func (t Timer) Year() int {
	return t.t.Year()
}

func (t Timer) Month() int {
	return int(t.t.Month()) + 1
}

func (t Timer) Day() int {
	return t.t.Day()
}

func (t Timer) Hour() int {
	return t.t.Hour()
}

func (t Timer) Minute() int {
	return t.t.Minute()
}

func (t Timer) Second() int {
	return t.t.Second()
}

func (t Timer) Add(d time.Duration) Timer {
	t.t = t.t.Add(d)
	return t
}

func (t Timer) DayStart(add int) Timer {
	y, m, d := t.t.Date()
	t.t = time.Date(y, m, d+add, 0, 0, 0, 0, location)
	return t
}

func (t Timer) DayEnd(add int) Timer {
	return t.DayStart(add + 1).Add(-time.Nanosecond)
}

func (t Timer) MonthStart(add int) Timer {
	y, m, _ := t.t.Date()
	t.t = time.Date(y, m+time.Month(add), 1, 0, 0, 0, 0, location)
	return t
}

func (t Timer) MonthEnd(add int) Timer {
	return t.MonthStart(add + 1).Add(-time.Nanosecond)
}
