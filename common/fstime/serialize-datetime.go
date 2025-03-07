/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: serialize time
@author: fanky
@version: 1.0
@date: 2024-06-03
**/

package fstime

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

type T_SerDateTime time.Time

func NowSerDateTime() T_SerDateTime {
	return T_SerDateTime(time.Now())
}

func (this *T_SerDateTime) Origin() time.Time {
	return time.Time(*this)
}

func (self T_SerDateTime) String() string {
	return time.Time(self).Format(time.DateTime)
}

func (self T_SerDateTime) GoString() string {
	return time.Time(self).Format(time.DateTime)
}

func (this *T_SerDateTime) Update(t time.Time) {
	*this = T_SerDateTime(t)
}

// GobEncoder 接口，可以让 gob 序列化
func (self T_SerDateTime) GobEncode() ([]byte, error) {
	ts := time.Time(self)
	return ts.GobEncode()
}

// GobDecoder 接口，可以让 gob 反序列化
func (this *T_SerDateTime) GobDecode(data []byte) error {
	var ts time.Time
	err := ts.GobDecode(data)
	if err != nil {
		return err
	}
	*this = T_SerDateTime(ts)
	return nil
}

// json unmarshal
func (this *T_SerDateTime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		var t T_SerDateTime
		*this = t
		return nil
	}
	b = bytes.Trim(b, `"`)
	t, err := time.ParseInLocation(time.DateTime, string(b), time.Local)
	if err != nil {
		return fmt.Errorf("time format error, it must be like %q", time.DateTime)
	}
	*this = T_SerDateTime(t)
	return nil
}

// json marshal
func (this *T_SerDateTime) MarshalJSON() ([]byte, error) {
	str := `"` + time.Time(*this).Format(time.DateTime) + `"`
	return []byte(str), nil
}

// 解释给数据库(存入数据库时调用)
func (self T_SerDateTime) Value() (driver.Value, error) {
	return time.Time(self), nil
}
