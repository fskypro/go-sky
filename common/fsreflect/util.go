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
	"errors"
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
// 获取 obj 参数的原始类型值的 reflect.Value 封装
func BaseRefValue(obj any) (reflect.Value, error) {
	tobj := reflect.TypeOf(obj)
	if tobj == nil {
		return reflect.ValueOf(nil), errors.New("obj mustn't be a no type nil value")
	}
	vobj := reflect.ValueOf(obj)
	if tobj.Kind() == reflect.Ptr {
		if vobj.IsNil() {
			return vobj, errors.New("obj must be a not nil value")
		}
		return vobj.Elem(), nil
	}
	return vobj, nil
}

// 获取 obj 参数的原始类型的 reflect.Type 封装
func BaseRefType(obj any) (reflect.Type, error) {
	tobj := reflect.TypeOf(obj)
	if tobj == nil {
		return nil, errors.New("obj mustn't be a no type nil value")
	}
	if tobj.Kind() == reflect.Ptr {
		return tobj.Elem(), nil
	}
	return tobj, nil
}

// 获取参数引用的最末端 reflect.Value
func EndRefValue(v any) reflect.Value {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() { return rv }

L:
	if rv.IsNil() { return rv }
	if rv.Type().Kind() == reflect.Ptr {
		rv = rv.Elem()
		goto L
	}
	return rv
}
