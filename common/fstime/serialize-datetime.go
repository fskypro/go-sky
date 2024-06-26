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
	"fmt"
	"time"
)

type T_SerDateTime time.Time

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
