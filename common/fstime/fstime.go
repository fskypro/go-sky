/**
@copyright: fantasysky 2016
@brief: 实现时间相关函数
@author: fanky
@version: 1.0
@date: 2019-02-21
**/

package fstime

import (
	"time"
)

// -----------------------------------------------------------------------------
// inner
// -----------------------------------------------------------------------------
// 指定星期值距离周起始日经历的天数
func _weekdayCount(weekday time.Weekday) int {
	if CWeekStart == WeekStartSunday { // 周日为起始日
		return int(weekday) - int(CWeekStart)
	}
	return (int(weekday) + 6) % 7
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
// -------------------------------------------------------------------
// 时间转换
// -------------------------------------------------------------------
// Days2Seconds 将天数转换为秒数
func Days2Seconds(days int) int64 {
	return int64(days) * OneDaySeconds
}

// 秒数转换为天数+时间
func Seconds2DaysTime(snds int) (d, h, m, s int) {
	d = snds / 86400
	snds = snds - d*86400
	h = snds / 3600
	snds = snds - h*3600
	m = snds / 60
	s = snds - m*60
	return
}

// -------------------------------------------------------------------
// 起始日期相关
// -------------------------------------------------------------------
// 获取指定年月的最后一天是本月的第几天
func LastDayOfMon(year int, mon time.Month) int {
	t := time.Date(year, mon, 1, 0, 0, 0, 0, time.Local)
	t = t.AddDate(0, 1, 0).Add(-time.Hour)
	return t.Day()
}

// ---------------------------------------------------------
// Dawn 获取指定时间的当天凌晨时间
func Dawn(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// WeekStart 获取指定时间所在星期的星期起始时间
func WeekStart(t time.Time) time.Time {
	days := _weekdayCount(t.Weekday())
	t = Dawn(t)
	return t.AddDate(0, 0, -days)
}

// 获取指定时间所在月份 1 号的凌晨时间
func MonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// ---------------------------------------------------------
// 获取两个时间点之间的距离天数
// 如果 t2 > t1 则返回正，如果 t2 < t1 则返回负
func DaysBetween(t1, t2 time.Time) float64 {
	return float64(t2.Sub(t1)) / float64(time.Hour*24)
}
