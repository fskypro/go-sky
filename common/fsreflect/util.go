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
