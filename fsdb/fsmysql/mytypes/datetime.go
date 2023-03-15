/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: date time wrap for mysql db
@author: fanky
@version: 1.0
@date: 2021-12-09
**/

package mytypes

import (
	"fmt"
	"time"
)

type T_DateTime string

func NewDateTimeViaGoTime(t time.Time) T_DateTime {
	return T_DateTime(t.Format("2006-01-02 15:04:05"))
}

func NowLocalDateTime() T_DateTime {
	return T_DateTime(time.Now().Format("2006-01-02 15:04:05"))
}

func NowUTCDateTime() T_DateTime {
	return T_DateTime(time.Now().UTC().Format("2006-01-02 15:04:05"))
}

func NewZeroDateTime() T_DateTime {
	return T_DateTime("0000-00-00 00:00:00")
}

func NewUnixStartDateTime() T_DateTime {
	return T_DateTime(time.Unix(0, 0).Format("2006-01-02 15:04:05"))
}

// ---------------------------------------------------------
func (self T_DateTime) String() string {
	if string(self) == "" {
		return "0000-00-00 00:00:00"
	}
	return string(self)
}

// 返回单引号括回的字符串，格式为：'2021-12-9 00:00:00'
func (self T_DateTime) Quote() string {
	return fmt.Sprintf("'%v'", self)
}

// loc 表示你认为 self 是什么时区
func (self T_DateTime) GoTime(loc *time.Location) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", string(self), loc)
}

func (self T_DateTime) MustGoTime(loc *time.Location) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", string(self), loc)
	if err != nil {
		panic(err)
	}
	return t
}

func (self T_DateTime) UTCGoTime() (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", string(self))
}

func (self T_DateTime) MustUTCGoTime() time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", string(self))
	if err != nil {
		panic(err)
	}
	return t
}

func (self T_DateTime) LocalGoTime() (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", string(self), time.Local)
}

func (self T_DateTime) MustLocalGoTime() time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", string(self), time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func (self T_DateTime) IsZeroTime() bool {
	return string(self) == "0000-00-00 00:00:00"
}

// -------------------------------------------------------------------
// 与 T_DateTime 时间比较
// -------------------------------------------------------------------
func (self T_DateTime) Eq(v T_DateTime) bool {
	return self == v
}

func (self T_DateTime) Le(v T_DateTime) bool {
	t, err := self.LocalGoTime()
	if err != nil {
		return false
	}
	tv, err := v.LocalGoTime()
	if err != nil {
		return false
	}
	return t.Before(tv)
}

func (self T_DateTime) LeEq(v T_DateTime) bool {
	if self == v {
		return true
	}
	return self.Le(v)
}

func (self T_DateTime) Lg(v T_DateTime) bool {
	t, err := self.LocalGoTime()
	if err != nil {
		return false
	}
	tv, err := v.LocalGoTime()
	if err != nil {
		return false
	}
	return t.After(tv)
}

func (self T_DateTime) LgEq(v T_DateTime) bool {
	if self == v {
		return true
	}
	return self.Lg(v)
}

// -------------------------------------------------------------------
// 与 go Time 比较
// -------------------------------------------------------------------
func (self T_DateTime) Before(v time.Time) bool {
	t, err := self.GoTime(v.Location())
	if err != nil {
		return false
	}
	return t.Before(v)
}

func (self T_DateTime) After(v time.Time) bool {
	t, err := self.GoTime(v.Location())
	if err != nil {
		return false
	}
	return t.After(v)
}
