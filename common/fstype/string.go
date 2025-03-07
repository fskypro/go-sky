/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: string
@author: fanky
@version: 1.0
@date: 2025-01-01
**/

package fstype

import (
	"reflect"
)

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
type T_String interface {
	~string
}

func IsString(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	return t.Kind() == reflect.String
}

// ---------------------------------------------------------
// 字符串与字符数组
type T_AllString interface {
	~string | ~[]byte | ~[]rune
}

func IsAllString(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	kind := t.Kind()
	switch kind {
	case reflect.String:
		return true // 底层类型是 string
	case reflect.Slice:
		// 判断是否是 []byte 或 []rune
		elem := t.Elem()
		if elem.Kind() == reflect.Uint8 { // []byte
			return true
		}
		if elem.Kind() == reflect.Int32 { // []rune
			return true
		}
	default:
		return false
	}
	return false
}
