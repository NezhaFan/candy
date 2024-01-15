package candy

import "time"

const (
	Datetime = "2006-01-02 15:04:05"
)

type Time struct {
	t time.Time
}

func NewTime[T string | int64 | int | time.Time](init T) Time {
	var v any = init
	switch v := v.(type) {
	case string:
		t, _ := time.ParseInLocation(time.DateTime, v, time.Local)
		return Time{t: t}
	case int:
		return Time{t: time.Unix(int64(v), 0)}
	case int64:
		return Time{t: time.Unix(v, 0)}
	case time.Time:
		return Time{t: v}
	}
	return Time{t: time.Now()}
}

func Now() Time {
	return Time{t: time.Now()}
}

// timex转Time
func (t Time) Time() time.Time {
	return t.t
}

// timex转字符串
func (t Time) String() string {
	return t.t.Format(Datetime)
}

// timex转时间戳
func (t Time) Unix() int64 {
	return t.t.Unix()
}

// 加 秒/分/时
func (t Time) Add(d time.Duration) Time {
	t.t = t.t.Add(d)
	return t
}

// 加天数
func (t Time) AddDay(i int) Time {
	t.t = t.t.AddDate(0, 0, i)
	return t
}

// 加月份
func (t Time) AddMonth(i int) Time {
	t.t = t.t.AddDate(0, i, 0)
	return t
}

// 加年
func (t Time) AddYear(i int) Time {
	t.t = t.t.AddDate(i, 0, 0)
	return t
}

// 当日的开始时间 (00:00:00)
func (t Time) DayStart() Time {
	y, m, d := t.t.Date()
	t.t = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	return t
}

// 当日的结束时间 (23:59:59)
func (t Time) DayEnd() Time {
	return t.DayStart().Add(time.Hour*24 - time.Nanosecond)
}

// 当周一的开始时间
func (t Time) WeekStart() Time {
	w := t.t.Weekday()
	if w == 0 {
		w = 7
	}
	t.t = t.t.Add(-time.Duration(w-1) * time.Hour * 24)
	return t.DayStart()
}

// 当周日的结束时间
func (t Time) WeekEnd() Time {
	return t.WeekStart().Add(time.Hour*24*7 - time.Nanosecond)
}

// 当月的开始时间 (00:00:00)
func (t Time) MonthStart() Time {
	y, m, _ := t.t.Date()
	t.t = time.Date(y, m, 1, 0, 0, 0, 0, t.t.Location())
	return t
}

// 当月的结束时间 (23:59:59)
func (t Time) MonthEnd() Time {
	return t.MonthStart().AddMonth(1).Add(-time.Nanosecond)
}
