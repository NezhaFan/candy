package candy

import "time"

type Timex struct {
	Time time.Time
}

func Now() Timex {
	return Timex{time.Now()}
}

// 注意字符串必须为 日期 + 时间 格式
func NewTimexByString(dt string) Timex {
	t, _ := time.ParseInLocation(time.DateTime, dt, time.Local)
	return Timex{t}
}

func NewTimexByUnix[E int | uint | int64 | uint64](i E) Timex {
	return Timex{time.Unix(int64(i), 0)}
}

// timex转字符串
func (t Timex) ToString() string {
	return t.Time.Format(time.DateTime)
}

// timex转时间戳
func (t Timex) ToUnix() int64 {
	return t.Time.Unix()
}

// 加天数
func (t Timex) AddDay(i int) Timex {
	t.Time = t.Time.AddDate(0, 0, i)
	return t
}

// 加月份
func (t Timex) AddMonth(i int) Timex {
	t.Time = t.Time.AddDate(0, i, 0)
	return t
}

// 加年
func (t Timex) AddYear(i int) Timex {
	t.Time = t.Time.AddDate(i, 0, 0)
	return t
}

// 当日的开始时间 (00:00:00)
func (t Timex) DayStart() Timex {
	y, m, d := t.Time.Date()
	t.Time = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	return t
}

// 当日的结束时间 (23:59:59 999999)
func (t Timex) DayEnd() Timex {
	t.Time = t.DayStart().Time.Add(time.Hour*24 - time.Nanosecond)
	return t
}

// 当周一的开始时间
func (t Timex) WeekStart() Timex {
	w := t.Time.Weekday()
	if w == time.Sunday {
		w = 7
	}
	y, m, d := t.Time.Date()
	t.Time = time.Date(y, m, d-int(w)+1, 0, 0, 0, 0, time.Local)
	return t.DayStart()
}

// 当周日的结束时间
func (t Timex) WeekEnd() Timex {
	t.Time = t.WeekStart().Time.Add(time.Hour*24*7 - time.Nanosecond)
	return t
}

// 当月的开始时间 (00:00:00)
func (t Timex) MonthStart() Timex {
	y, m, _ := t.Time.Date()
	t.Time = time.Date(y, m, 1, 0, 0, 0, 0, t.Time.Location())
	return t
}

// 当月的结束时间 (23:59:59)
func (t Timex) MonthEnd() Timex {
	t.Time = t.MonthStart().AddMonth(1).Time.Add(-time.Nanosecond)
	return t
}

func (t Timex) GetDate() string {
	return t.Time.Format(time.DateOnly)
}

func (t Timex) GetTime() string {
	return t.Time.Format(time.TimeOnly)
}
