/**
* @file: define.go
* @copyright: 2016 fantasysky
* @author: fanky
* @version: 1.0
* @date: 2018-09-05
**/

package fserror

import "reflect"

// 标准错误的反射类型
var RTypeStdError = reflect.TypeOf((*error)(nil)).Elem()

// nil 错误的 value 值
var NilErrorValue = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
