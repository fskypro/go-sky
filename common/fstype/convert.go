/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: types converter
@author: fanky
@version: 1.0
@date: 2023-09-12
**/

package fstype

import (
	"fmt"
	"reflect"
	"strconv"
)

// -------------------------------------------------------------------
// 合法转换
// -------------------------------------------------------------------
// 判断参数中的值是否能转换为范型参数指定的类型
func CanConvert[T any](v any) bool {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return false
	}
	var target T
	t := reflect.TypeOf(target)
	return rv.CanConvert(t)
}

// 将参数值转换为范型参数类型
func Convert[T any](v any) (T, error) {
	var target T
	if CanConvert[T](v) {
		t := reflect.TypeOf(target)
		target = reflect.ValueOf(v).Convert(t).Interface().(T)
		return target, nil
	}
	return target, fmt.Errorf("value of type %v con't convert to %v",
		reflect.TypeOf(v), reflect.TypeOf(target))
}

// ---------------------------------------------------------
// 布尔型转换为数值：
// true 转换为 1
// false 转换为 0
func BoolToNumber[T T_Number](v bool) T {
	if v {
		return T(1)
	}
	return T(0)
}

// -------------------------------------------------------------------
// 字符串转换为数值类型
// -------------------------------------------------------------------
// 字符串转换为数值类型，如果转换失败，则返回错误
func StrToNum[N T_Number, S T_AllString](str S) (N, error) {
	var v N
	s := string(str)
	if IsFNumber(v) {
		ret, err := strconv.ParseFloat(s, 64)
		return N(ret), err
	}
	if IsUNumber(v) {
		ret, err := strconv.ParseUint(s, 10, 64)
		return N(ret), err
	}
	if IsINumber(v) {
		ret, err := strconv.ParseInt(s, 10, 64)
		return N(ret), err
	}
	return v, fmt.Errorf("can't convert string %q to type %v", s, reflect.TypeOf(v))
}

// 字符串转换为数值类型，如果转换失败，则返回默认值
func StrToNumOrDef[N T_Number, S T_AllString](str S, v N) N {
	s := string(str)
	if IsFNumber(v) {
		ret, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return N(ret)
		}
	}
	if IsUNumber(v) {
		ret, err := strconv.ParseUint(s, 10, 64)
		if err == nil {
			return N(ret)
		}
	}
	if IsINumber(v) {
		ret, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			return N(ret)
		}
	}
	return v
}
