/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 日时间
@author: fanky
@version: 1.0
@date: 2022-08-16
**/

package fstime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fsky.pro/fstype"
)

type T_DayTime uint

func NewDayTime[T fstype.T_IUNumber](h, m, s T) (T_DayTime, error) {
	if h < 0 || h > 23 {
		return 0, errors.New("hours in day time must be 0~23")
	}
	if m < 0 || m > 59 {
		return 0, errors.New("minutes in day time must be 0~59")
	}
	if s < 0 || s > 59 {
		return 0, errors.New("seconds in day time must be 0~59")
	}
	v := uint(h<<16) + uint(m<<8) + uint(s)
	return T_DayTime(v), nil
}

func MustDayTime[T fstype.T_IUNumber](h, m, s T) T_DayTime {
	t, err := NewDayTime(h, m, s)
	if err != nil {
		panic(err)
	}
	return t
}

// 格式为：00:00:00
// 允许空字符串，空字符串返回 00:00:01
func ParseDayTime(text string) (T_DayTime, error) {
	if strings.TrimSpace(text) == "" {
		return ZeroDayTime(), nil
	}
	hms := strings.Split(text, ":")
	if len(hms) != 3 {
		return 0, errors.New("error day time format string %q, it must be like '00:00:00'")
	}
	h, err := strconv.Atoi(hms[0])
	if err != nil {
		return 0, errors.New("error day time format string %q, it must be like '00:00:00'")
	}
	m, err := strconv.Atoi(hms[1])
	if err != nil {
		return 0, errors.New("error day time format string %q, it must be like '00:00:00'")
	}
	s, err := strconv.Atoi(hms[2])
	if err != nil {
		return 0, errors.New("error day time format string %q, it must be like '00:00:00'")
	}
	return NewDayTime(h, m, s)
}

func ParseDayTimeDefault(text string, def T_DayTime) T_DayTime {
	t, err := ParseDayTime(text)
	if err != nil {
		return def
	}
	return t
}

// 截取 time 的时间部分
func DayTimeFromGoTime(t time.Time) T_DayTime {
	return T_DayTime((t.Hour() << 16) + (t.Minute() << 8) + (t.Second()))
}

func ZeroDayTime() T_DayTime {
	return T_DayTime(0)
}

// -------------------------------------------------------------------
// methods
// -------------------------------------------------------------------
func (self T_DayTime) String() string {
	return fmt.Sprintf("%d:%d:%d", self.Hour(), self.Minute(), self.Second())
}

// 与指定时间的日期部分合并
func (self T_DayTime) WithGoTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, self.Hour(), self.Minute(), self.Second(), 0, t.Location())
}

// ---------------------------------------------------------
func (self T_DayTime) Hour() int {
	return int(self >> 16)
}

func (self T_DayTime) Minute() int {
	return int((self >> 8) & 0xff)
}

func (self T_DayTime) Second() int {
	return int(self & 0xff)
}

// ---------------------------------------------------------
// 距离凌晨时的秒数
func (self T_DayTime) Seconds() int {
	return self.Hour()*3600 + self.Minute()*60 + self.Second()
}

func (self T_DayTime) DuSeconds() time.Duration {
	return time.Duration(self.Seconds())
}

// 距离凌晨时的分钟数
func (self T_DayTime) Minutes() int {
	return self.Hour()*60 + self.Minute()
}

func (self T_DayTime) DuMinute() time.Duration {
	return time.Duration(self.Minutes())
}

// ---------------------------------------------------------
// 添加分时秒，可以是负数
// 自动进退日期，譬如：T_DayTime("00:00:00").Add(1,0,0,) == T_DayTime("23:00:00")
func (self T_DayTime) Add(h, m, s int) T_DayTime {
	addSeconds := h*3600 + m*60 + s
	mySeconds := self.Hour()*3600 + self.Minute()*60 + self.Second()
	newSeconds := addSeconds + mySeconds
	daySeconds := newSeconds % 86400
	if daySeconds < 0 {
		daySeconds += 86400
	}
	h = daySeconds / 3600
	m = (daySeconds - (h * 3600)) / 60
	s = daySeconds - h*3600 - m*60
	return T_DayTime((h << 16) + (m << 8) + s)
}

// 时间差
func (self T_DayTime) Sub(t T_DayTime) *S_Hms {
	snds := self.Seconds() - t.Seconds()
	return NewHmsFromDuration(time.Second * time.Duration(snds))
}
