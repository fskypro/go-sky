/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: time utils
@author: fanky
@version: 1.0
@date: 2022-08-17
**/

package fstime

import "time"

// -----------------------------------------------------------------------------
// inner
// -----------------------------------------------------------------------------
// 本地时间与 UTC 时间差(微秒)
var _localUTCSpace int64

func init() {
	local, _ := time.ParseInLocation("2006-01-02 15:04:04", "2006-01-02 15:04:04", time.Local)
	utc, _ := time.ParseInLocation("2006-01-02 15:04:04", "2006-01-02 15:04:04", time.UTC)
	_localUTCSpace = utc.UnixMicro() - local.UnixMicro()
}

// 指定星期值距离周起始日经历的天数
func _weekdayCount(weekday time.Weekday) int {
	if weekStart == WeekStartSunday { // 周日为起始日
		return int(weekday) - int(weekStart)
	}
	return (int(weekday) + 6) % 7
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
// 传入微秒级 UTC 时间戳，转换为时间
func UTCUnixMicro(umicro int64) time.Time {
	return time.UnixMicro(umicro + _localUTCSpace)
}

// 传入毫秒级 UTC 时间戳，转换为时间
func UTCUnixMilli(umilli int64) time.Time {
	return time.UnixMilli(umilli + (_localUTCSpace / 1000))
}

// 传入秒级 UTC 时间戳，转换为时间
func UTCUnix(usec int64) time.Time {
	return time.Unix(usec+(_localUTCSpace/1000000), 0)
}

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

// 秒数转换为小时+分钟+秒
func Seconds2HoursTime(snds int) (h, m, s int) {
	h = snds / 3600
	m = (snds - h*3600) / 60
	s = snds % 60
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
	days := _weekdayCount(t.Weekday()) // 指定时间与起始日相差的天数
	t = Dawn(t)
	return t.AddDate(0, 0, -days)
}

// 获取指定时间所在月份 1 号的凌晨时间
func MonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// 指定时间所在月的最后一天是第几天(一号为第一天)
func MonthEndDay(t time.Time) int {
	return MonthStart(t.AddDate(0, 1, 0)).Add(-time.Hour).Day()
}

// ---------------------------------------------------------
// 获取两个时间点之间的距离天数
// 如果 t2 > t1 则返回正，如果 t2 < t1 则返回负
func DaysBetween(t1, t2 time.Time) int {
	neg := 1
	if t1.After(t2) {
		t1, t2 = t2, t1
		neg = -1
	}
	days := 0
	for {
		temp := t1.AddDate(200, 0, 0)
		if !temp.Before(t2) {
			days += int(t2.Sub(t1).Hours() / 24)
			break
		}
		days += int(temp.Sub(t1).Hours() / 24)
		t1 = temp
	}
	return days * neg
}

// 获取两个时间点之间的时间间隔
// 如果 t2 > t1 则返回正，如果 t2 < t1 则返回负
func HmsBetween(t1, t2 time.Time) *S_Hms {
	return NewHmsFromDuration(t2.Sub(t1))
}

// ---------------------------------------------------------
// 获取计算机元年到指定是时间的距离天数
func DaysFromUnixTime(t time.Time) int {
	utime := time.Date(1970, 1, 1, 0, 0, 0, 0, t.Location())
	return int(t.Sub(utime).Hours()) / 24
}

// 给出距离计算机元年天数，返回实际时间
func DateToUnixTime(days int, loc *time.Location) time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, loc).AddDate(0, 0, days)
}
