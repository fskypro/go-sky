/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: reflect util
@author: fanky
@version: 1.0
@date: 2021-06-21
**/

package fsreflect

import (
	"fmt"
	"reflect"
)

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
// 将 v 强制转换为类型 t 的值，尽可能的转换，如：
//   string("123") -> int(123)
//   int(123) -> "123"
//   int64(456) -> int32(456)
func hardConvert(v interface{}, t reflect.Type) (reflect.Value, bool) {
	if v == nil {
		return reflect.Zero(t), true
	}
	tv := reflect.TypeOf(v)
	if tv == t {
		return reflect.ValueOf(v), true
	} else if t.Kind() == reflect.String {
		return reflect.ValueOf(fmt.Sprintf("%v", v)), true
	} else if tv.Kind() == reflect.String {
		if cv, ok := strTo(v.(string), t); ok {
			return cv, true
		}
		return reflect.Zero(t), false
	} else if tv.ConvertibleTo(t) {
		return reflect.ValueOf(v).Convert(t), true
	}
	return reflect.Zero(t), false
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 判断参数 a 能否转换为 b 的类型
func CanConvertToTypeOf(a interface{}, b interface{}) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	return ta.ConvertibleTo(tb) && tb.ConvertibleTo(ta)
}
