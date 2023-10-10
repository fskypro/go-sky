/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: types definations
@author: fanky
@version: 1.0
@date: 2022-07-29
**/

package fstype

import (
	"reflect"

	"fsky.pro/fscollection"
)

// -------------------------------------------------------------------
// 类型定义
// -------------------------------------------------------------------
// 数值类型
type T_Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

var number_kinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	reflect.Float32, reflect.Float64,
}

func IsNumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil { return false }
	return fscollection.SliceHas(number_kinds, t.Kind())
}

// ---------------------------------------------------------
// 整数
type T_IUNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

var iunumber_kinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
}

func IsIUNumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil { return false }
	return fscollection.SliceHas(iunumber_kinds, t.Kind())
}

// ---------------------------------------------------------
// 有符号整数
type T_INumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

var inumber_kinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int16,
}

func IsINumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil { return false }
	return fscollection.SliceHas(iunumber_kinds, t.Kind())
}

// ---------------------------------------------------------
// 无符号整数
type T_UNumber interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

var unumber_kinds = []reflect.Kind{
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64,
}

func IsUNumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil { return false }
	return fscollection.SliceHas(unumber_kinds, t.Kind())
}

// ---------------------------------------------------------
// 浮点数
type T_FNumber interface {
	~float32 | ~float64
}

var fnumber_kinds = []reflect.Kind{
	reflect.Float32, reflect.Float64,
}

func IsFNumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil { return false }
	return fscollection.SliceHas(fnumber_kinds, t.Kind())
}

// ---------------------------------------------------------
// 字符串与字符数组
type T_AllString interface {
	string | []byte | []rune
}

var allstring_kinds = []reflect.Kind{
	reflect.String,
	reflect.SliceOf(reflect.TypeOf(new(byte)).Elem()).Kind(),
	reflect.SliceOf(reflect.TypeOf(new(rune)).Elem()).Kind(),
}

func IsAllString(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil { return false }
	return fscollection.SliceHas(allstring_kinds, t.Kind())
}
