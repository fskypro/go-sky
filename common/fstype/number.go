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
	"slices"
)

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
var number_kinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	reflect.Float32, reflect.Float64,
}

var iunumber_kinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
}

var inumber_kinds = []reflect.Kind{
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int16,
}

var unumber_kinds = []reflect.Kind{
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64,
}

var fnumber_kinds = []reflect.Kind{
	reflect.Float32, reflect.Float64,
}

// -------------------------------------------------------------------
// 类型定义
// -------------------------------------------------------------------
// 数值类型
type T_Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// 是否是数值类型
func IsNumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	return slices.Contains(number_kinds, t.Kind())
}

// ---------------------------------------------------------
// 整数
type T_IUNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func IsIUNumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	return slices.Contains(iunumber_kinds, t.Kind())
}

// ---------------------------------------------------------
// 有符号整数
type T_INumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func IsINumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	return slices.Contains(iunumber_kinds, t.Kind())
}

// ---------------------------------------------------------
// 无符号整数
type T_UNumber interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func IsUNumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	return slices.Contains(unumber_kinds, t.Kind())
}

// ---------------------------------------------------------
// 浮点数
type T_FNumber interface {
	~float32 | ~float64
}

func IsFNumber(v any) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	return slices.Contains(fnumber_kinds, t.Kind())
}
