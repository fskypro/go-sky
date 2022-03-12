/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: unix format time
@author: fanky
@version: 1.0
@date: 2021-12-24
**/

package mytypes

import (
	"time"
)

type T_UnixTime int64

func NewUnixTime(Y int, M time.Month, D int, h, m, s int, loc *time.Location) T_UnixTime {
	if loc == nil {
		loc = time.Local
	}
	return T_UnixTime(time.Date(Y, M, D, h, m, s, 0, loc).Unix())
}

func NowDateTime() T_UnixTime {
	return T_UnixTime(time.Now().Unix())
}

func NowLocalUnixTime() T_UnixTime {
	return T_UnixTime(time.Now().Local().Unix())
}

func NowUTCUnixTime() T_UnixTime {
	return T_UnixTime(time.Now().UTC().Unix())
}

// -------------------------------------------------------------------
// methods
// -------------------------------------------------------------------
func (self T_UnixTime) GoTime(loc *time.Location) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", self.String(), loc)
	return t
}

func (self T_UnixTime) UTCGoTime() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", self.String())
	return t
}

func (self T_UnixTime) LocalGoTime() time.Time {
	return time.Unix(int64(self), 0)
}

func (self T_UnixTime) String() string {
	return time.Unix(int64(self), 0).Format("2006-01-02 15:04:05")
}

func (self T_UnixTime) Format(layout string) string {
	return self.LocalGoTime().Format(layout)
}

func (self T_UnixTime) Unix() int64 {
	return int64(self)
}

// 是否等于传入值
func (self T_UnixTime) Eq(dt int64) bool {
	return dt == int64(self)
}

// 是否大于传入值
func (self T_UnixTime) Lg(dt int64) bool {
	return int64(self) > dt
}

// 是否大于等于传入值
func (self T_UnixTime) LgEq(dt int64) bool {
	return int64(self) >= dt
}

// 是否小于传入值
func (self T_UnixTime) Le(dt int64) bool {
	return int64(self) < dt
}

// 是否小于等于传入值
func (self T_UnixTime) LeEq(dt int64) bool {
	return int64(self) <= dt
}

// 增时
func (self T_UnixTime) Add(dt int64) T_UnixTime {
	return T_UnixTime(int64(self) + dt)
}

// 减时
func (self T_UnixTime) Sub(dt int64) T_UnixTime {
	return T_UnixTime(int64(self) - dt)
}

// -------------------------------------------------------------------
// 与 go Time 比较
// -------------------------------------------------------------------
func (self T_UnixTime) Before(v time.Time) bool {
	return self.GoTime(v.Location()).Before(v)
}

func (self T_UnixTime) After(v time.Time) bool {
	return self.GoTime(v.Location()).After(v)
}
