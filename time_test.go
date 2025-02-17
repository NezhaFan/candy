package candy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTime(t *testing.T) {
	timex := NewTimexByString("2000-01-30 12:00:00")

	// 当日开始
	assert.Equal(t, timex.DayStart().ToString(), "2000-01-30 00:00:00")
	// 当日结束
	assert.Equal(t, timex.DayEnd().ToString(), "2000-01-30 23:59:59")

	// 月初
	assert.Equal(t, timex.MonthStart().GetDate(), "2000-01-01")
	assert.Equal(t, timex.MonthStart().ToString(), "2000-01-01 00:00:00")
	// 月末
	assert.Equal(t, timex.MonthEnd().GetDate(), "2000-01-31")
	assert.Equal(t, timex.MonthEnd().ToString(), "2000-01-31 23:59:59")

	// 当天星期日
	assert.Equal(t, timex.Time.Weekday(), time.Sunday)
	// 当周周一
	assert.Equal(t, timex.WeekStart().GetDate(), "2000-01-24")
	// 当周周日
	assert.Equal(t, timex.WeekEnd().GetDate(), "2000-01-30")

	// 加分
	assert.Equal(t, Timex{timex.Time.Add(time.Minute * 60)}.ToString(), "2000-01-30 13:00:00")

	// 加时
	assert.Equal(t, Timex{timex.Time.Add(time.Hour)}.ToString(), "2000-01-30 13:00:00")

	// 加天
	assert.Equal(t, timex.AddDay(1).GetDate(), "2000-01-31")
	assert.Equal(t, timex.AddDay(2).GetDate(), "2000-02-01")

	// 加月 （没有 02-30，自动到 03-01）
	assert.Equal(t, timex.AddMonth(1).GetDate(), "2000-03-01")

	// 加年
	assert.Equal(t, timex.AddYear(1).GetDate(), "2001-01-30")

	assert.Equal(t, Now().ToUnix(), time.Now().Unix())
}
