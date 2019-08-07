/**
@copyright: fantasysky 2016
@brief: 关于时间的定义
@author: fanky
@version: 1.0
@date: 2019-03-04
**/

package fstime

import "time"

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
//	WeekStarSunday、WeekStarMonday
var CWeekStart T_WeekStart = WeekStartMonday
