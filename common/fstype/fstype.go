/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: fstype utils
@author: fanky
@version: 1.0
@date: 2022-10-20
**/

package fstype

import "reflect"

// -------------------------------------------------------------------
// type type
// -------------------------------------------------------------------
type s_T[T any] struct{}

// 判断两个模板类型是否是同一个类型
// 注意：重定义类型认为是不同类型：
// type S string
// SameType[S, string]() == false
func SameType[T1, T2 any]() bool {
	return any(s_T[T1]{}) == any(s_T[T2]{})
}

// 判断两个模板类型是否具有相同的原始类型
// 如，以下类型判断为 true
//   type S string
//   var s S = "xxx"
//   SameOriginType[S, string]() == true
func SameOriginType[T1, T2 any]() bool {
    return reflect.TypeOf(new(T1)).Elem().Kind() == reflect.TypeOf(new(T2)).Elem().Kind()
}

// -------------------------------------------------------------------
// value type
// -------------------------------------------------------------------
// 判断传入的值是否是模板参数中的类型
// 注意，类似下面这种，判断是 false：
//   type S string
//   var s S = "xxx"
//   IsType[string](s) == false
func IsType[T any](v any) bool {
	switch v.(type){
	case T: return true
	}
	return false
}

// 判断传入的值的类型与模板类型是否有相同的原始类型
// 如，以下类型判断为 true
//   type S string
//   var s S = "xxx"
//   IsOriginType[string](s) == true
func IsOriginType[T any](v any) bool{
    return reflect.TypeOf(new(T)).Elem().Kind() == reflect.TypeOf(v).Kind()
}

