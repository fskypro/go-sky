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
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type T_Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// 高 16 位表示时
// 低 8 位表示秒，中八位表示分
type T_DayTime uint32

func NewDayTime[T T_Int](h, m, s T) T_DayTime {
	return T_DayTime(0).Add(int(h), int(m), int(s))
}

// 自动从 json 值中解释
func (this *T_DayTime) UnmarshalJSON(b []byte) (err error) {
	str := strings.Trim(string(b), `"`)
	*this, err = ParseDayTime(str)
	return
}

// 序列号到 json 中
func (this *T_DayTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%02d:%02d:%02d"`, this.Hour(), this.Minute(), this.Second())), nil
}

// 字符串形式写入数据库
func (self T_DayTime) Value() (driver.Value, error) {
	return self.String(), nil
}

// 从数据库扫描
func (this *T_DayTime) Scan(value any) error {
	if value == nil {
		return nil
	}

	var str string
	switch value.(type) {
	case []byte:
		str = string(value.([]byte))
	case string:
		str = value.(string)
	default:
		return fmt.Errorf("can't convert type %v value %q to %v", reflect.TypeOf(value), value, reflect.TypeOf(*this))
	}
	t, err := ParseDayTime(str)
	if err != nil {
		return fmt.Errorf("can't convert type %v value %v to %v", reflect.TypeOf(value), value, reflect.TypeOf(*this))
	}
	*this = t

	return nil
}

// 格式为：00:00:00
// 允许空字符串，空字符串返回 00:00:00
func ParseDayTime(text string) (T_DayTime, error) {
	if strings.TrimSpace(text) == "" {
		return ZeroDayTime(), nil
	}
	hms := strings.Split(text, ":")
	if len(hms) != 3 {
		return 0, fmt.Errorf("error day time format string %q, it must be like '00:00:00'", text)
	}
	h, err := strconv.Atoi(hms[0])
	if err != nil {
		return 0, fmt.Errorf("error day time format string %q, it must be like '00:00:00'", text)
	}
	m, err := strconv.Atoi(hms[1])
	if err != nil {
		return 0, fmt.Errorf("error day time format string %q, it must be like '00:00:00'", text)
	}
	s, err := strconv.Atoi(hms[2])
	if err != nil {
		return 0, fmt.Errorf("error day time format string %q, it must be like '00:00:00'", text)
	}
	return NewDayTime(h, m, s), nil
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
	return fmt.Sprintf("%02d:%02d:%02d", self.Hour(), self.Minute(), self.Second())
}

// 与指定时间的日期部分合并
func (self T_DayTime) WithGoTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, self.Hour(), self.Minute(), self.Second(), 0, t.Location())
}

// 该时间点在今天的时间值
func (self T_DayTime) TodayTime() time.Time {
	return self.WithGoTime(time.Now())
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
