/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: types converter
@author: fanky
@version: 1.0
@date: 2023-09-12
**/

package fsconv

import (
	"fmt"
	"reflect"

	"fsky.pro/fstype"
)

// 布尔型转换为整型
func BoolToNumber[T fstype.T_Number](v bool) T {
	if v { return T(1) }
	return T(0)
}

// 将 T1 类型的值转换为 T2
func Convert[T1, T2 any](v T1) (T2, error) {
	v1 :=reflect.ValueOf(v)
	var v2 T2
	t2 := reflect.TypeOf(v2)
	if !v1.IsValid() || t2 == nil {
		return v2, fmt.Errorf("no type nil value can't converted to any type")
	}
	if !v1.CanConvert(t2)  {
		return v2, fmt.Errorf("value of type %v con't convert to %v", v1.Type(), t2)
	}
	rv := v1.Convert(t2)
	if rv.CanInterface() {
		v2 = rv.Interface().(T2)
		return v2, nil
	}
	return v2, fmt.Errorf("value of type %v con't convert to %v", v1.Type(), t2)
}
