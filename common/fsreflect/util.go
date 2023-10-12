/**
* @file: util.go
* @copyright: 2016 fantasysky
* @author: fanky
* @version: 1.0
* @date: 2018-08-30
 */

// 通用工具套件
package fsreflect

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

// -------------------------------------------------------------------
// IsExposed 判断指定名称是否为可导出名称（大写字母开头）
func IsExposed(name string) bool {
	first, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(first)
}

// IsExposedOrBuiltinType 判断指定类型是否是可导出类型，或者是匿名类型
func IsExposedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return IsExposed(t.Name()) || t.PkgPath() == ""
}

// -------------------------------------------------------------------
// 获取 obj 参数的原始类型的 reflect.Type 封装
// 如果传入的 t 为 nil，则返回值也是 nil
func BaseRefType(t reflect.Type) reflect.Type {
	if t == nil { return t }
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// 获取参数引用的最末端 reflect.Value
// 传出参数有可能：
//   参数1 == -1：则返回值2 IsValid() == false
//   参数1 == 0 ：则返回值2 IsNil() == true
//   参数1 == 1 ：则返回值2 为非 Ptr Value
func BaseRefValue(v reflect.Value) (int, reflect.Value) {
	for v.IsValid() {
		if v.Type().Kind() != reflect.Ptr {
			break
		}
		if v.IsNil() {
			return 0, v
		}
		v = v.Elem()
	}
	if v.IsValid() {
		return 1, v
	}
	return -1, v
}
