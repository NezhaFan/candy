package timex

import "time"

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
func SetLocation(name string) (err error) {
	location, err = time.LoadLocation(name)
	return
}

// 对时间的操作。方法不使用指针，每次对时间的操作都是独立的。
type Timex struct {
	t time.Time
}

// 生成Timex
func New(t interface{}) Timex {
	switch v := t.(type) {
	case int64:
		return Timex{t: time.Unix(v, 0)}
	case string:
		t, _ := time.ParseInLocation(DateTime, v, location)
		return Timex{t: t}
	default:
		return Timex{t: time.Now().In(location)}
	}
}

// Timex转Time
func (t Timex) Time() time.Time {
	return t.t
}

// Timex转字符串
func (t Timex) String() string {
	return t.t.Format(DateTime)
}

// Timex转时间戳
func (t Timex) Unix() int64 {
	return t.t.Unix()
}

// 加 秒/分/时
func (t Timex) Add(d time.Duration) Timex {
	t.t = t.t.Add(d)
	return t
}

// 加天数
func (t Timex) AddDay(i int) Timex {
	t.t = t.t.AddDate(0, 0, i)
	return t
}

// 加月份
func (t Timex) AddMonth(i int) Timex {
	t.t = t.t.AddDate(0, i, 0)
	return t
}

// 加年
func (t Timex) AddYear(i int) Timex {
	t.t = t.t.AddDate(i, 0, 0)
	return t
}

// 当日的开始时间 (00:00:00)
func (t Timex) DayStart() Timex {
	y, m, d := t.t.Date()
	t.t = time.Date(y, m, d, 0, 0, 0, 0, location)
	return t
}

// 当日的结束时间 (23:59:59)
func (t Timex) DayEnd() Timex {
	return t.DayStart().Add(time.Hour*24 - time.Nanosecond)
}

// 当周一的开始时间
func (t Timex) WeekStart() Timex {
	w := t.t.Weekday()
	if w == 0 {
		w = 7
	}
	t.t = t.t.Add(-time.Duration(w-1) * time.Hour * 24)
	return t.DayStart()
}

// 当周日的结束时间
func (t Timex) WeekEnd() Timex {
	return t.WeekStart().Add(time.Hour*24*7 - time.Nanosecond)
}

// 当月的开始时间 (00:00:00)
func (t Timex) MonthStart() Timex {
	y, m, _ := t.t.Date()
	t.t = time.Date(y, m, 1, 0, 0, 0, 0, t.t.Location())
	return t
}

// 当月的结束时间 (23:59:59)
func (t Timex) MonthEnd() Timex {
	return t.MonthStart().AddMonth(1).Add(-time.Nanosecond)
}
