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
// ---------------------------------------------------------
// 时间转换
// ---------------------------------------------------------
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

// ---------------------------------------------------------
// 起始日期相关
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
