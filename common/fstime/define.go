/**
@copyright: fantasysky 2016
@brief: 关于时间的定义
@author: fanky
@version: 1.0
@date: 2019-03-04
**/

package fstime

import "time"

var ZeroTime, _ = time.Parse("", "") // 0000-01-01 00:00:00 +0000 UTC
var UnixStart = time.Unix(0, 0)      // 1970-01-01 08:00:00 +0800 CST

const (
	OneDaySeconds int64 = 24 * 3600 // 一天拥有的秒数
)

// -------------------------------------------------------------------
// 星期起始日定义
// -------------------------------------------------------------------
type T_WeekStart time.Weekday

const (
	WeekStartSunday T_WeekStart = T_WeekStart(time.Sunday) // 以周日作为星期起始日
	WeekStartMonday T_WeekStart = T_WeekStart(time.Monday) // 以周一作为星期起始日
)

// 默认以周一作为星期起始日
// 可以通过设置该值为：
//
//	WeekStarSunday、WeekStarMonday
var weekStart T_WeekStart = WeekStartMonday

// 设置以哪天作为星期开始日
var SetWeekStart = func(ws T_WeekStart) { weekStart = ws }
