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

type T_SerDateOnly time.Time

func (this *T_SerDateOnly) Origin() time.Time {
	return time.Time(*this)
}

func (this *T_SerDateOnly) String() string {
	return time.Time(*this).Format(time.DateTime)
}

func (this *T_SerDateOnly) Update(t time.Time) {
	*this = T_SerDateOnly(t)
}

// json unmarshal
func (this *T_SerDateOnly) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		var t T_SerDateOnly
		*this = t
		return nil
	}
	b = bytes.Trim(b, `"`)
	t, err := time.ParseInLocation(time.DateOnly, string(b), time.Local)
	if err != nil {
		return fmt.Errorf("time format error, it must be like %q", time.DateOnly)
	}
	*this = T_SerDateOnly(t)
	return nil
}

// json marshal
func (this *T_SerDateOnly) MarshalJSON() ([]byte, error) {
	str := `"` + time.Time(*this).Format(time.DateOnly) + `"`
	return []byte(str), nil
}
