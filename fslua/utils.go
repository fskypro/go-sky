/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: utils
@author: fanky
@version: 1.0
@date: 2024-12-10
**/

package fslua

import (
	"fmt"
	"reflect"
	"unsafe"

	glua "github.com/yuin/gopher-lua"
	lua "github.com/yuin/gopher-lua"
)

type s_Util struct{}

// 将 go 值转换为 lua 值
func (s_Util) refGoValueToLuaValue(rv reflect.Value) (glua.LValue, error) {
	if !rv.IsValid() {
		return glua.LNil, nil
	}
	if !rv.CanInterface() {
		if !rv.CanAddr() {
			return glua.LNil, fmt.Errorf("value of type %v is unaccessable", rv.Type())
		}
		rv = reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem()
	}

	glvalue := reflect.TypeOf((*glua.LValue)(nil)).Elem()
	if rv.CanConvert(glvalue) {
		rv = rv.Convert(glvalue)
		return rv.Interface().(glua.LValue), nil
	}
	if rv.CanConvert(rtNumnber) {
		rv = rv.Convert(rtNumnber)
		return glua.LNumber(rv.Interface().(float64)), nil
	}
	if rv.CanConvert(rtBool) {
		rv = rv.Convert(rtBool)
		return glua.LBool(rv.Interface().(bool)), nil
	}
	if rv.CanConvert(rtString) {
		rv = rv.Convert(rtString)
		return glua.LString(rv.Interface().(string)), nil
	}

	tb, err := s_Table{}.refMarshalTable(rv)
	if err != nil {
		return glua.LNil, fmt.Errorf("can't convert value which type of %v to lua value", rv.Type())
	}
	return tb, nil
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
// 设置 lua 路径
func SetPaths(L *glua.LState, paths ...string) {
	currentPath := ""
	for _, path := range paths {
		currentPath += ";" + path
	}
	packageTable := L.GetGlobal("package").(*glua.LTable)
	L.SetField(packageTable, "path", lua.LString(currentPath))
}

// 追加 lua 路径
func AddPaths(L *glua.LState, paths ...string) {
	packageTable := L.GetGlobal("package").(*glua.LTable)

	// 获取当前 package.path
	currentPath := L.GetField(packageTable, "path").String()

	for _, path := range paths {
		currentPath += ";" + path
	}
	L.SetField(packageTable, "path", lua.LString(currentPath))
}

// -------------------------------------------------------------------
// 将 lua 值转换为 go 值
func LuaValueToGoValue[T any](lv glua.LValue) (T, error) {
	var rv T
	var rt = reflect.TypeOf(rv)
	if rt == nil {
		// 当 T 当为 interface{} 时，走到这里
		var v any
		if lv.Type() == glua.LTNumber {
			v = float64(lv.(glua.LNumber))
		} else if lv.Type() == glua.LTBool {
			v = bool(lv.(glua.LBool))
		} else if lv.Type() == glua.LTString {
			v = string(lv.(glua.LString))
		} else {
			return rv, fmt.Errorf("can't convert lua value of type %q to go type %v", lv.Type(), rt)
		}
		return v.(T), nil
	}

	if lv.Type() == glua.LTTable {
		if err := UnmarshalTable(lv.(*glua.LTable), &rv); err != nil {
			return rv, fmt.Errorf("can't convert lua table to type %v, %v", rt, err)
		}
		return rv, nil
	}

	v := reflect.ValueOf(lv)
	if !v.CanConvert(rt) {
		return rv, fmt.Errorf("can't convert lua type %v to go type %v", lv.Type(), rt)
	}
	if !v.CanInterface() {
		if !v.CanAddr() {
			return rv, fmt.Errorf("value is unaccessable")
		}
		v = reflect.NewAt(v.Type(), v.UnsafePointer()).Elem()
	}
	return v.Convert(rt).Interface().(T), nil
}

// 将 go 值转换为 lua 值
func GoValueToLuaValue(v any) (glua.LValue, error) {
	return s_Util{}.refGoValueToLuaValue(reflect.ValueOf(v))
}

// 解释 lua value
func ParseLuaValue(lv glua.LValue, v any) error {
	var rv = reflect.ValueOf(v)
	var rt = reflect.TypeOf(v)
	if rt == nil || rt.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("out value must be a non-nil value")
	}
	rt = rt.Elem()
	rv = rv.Elem()

	if lv.Type() != lua.LTTable {
		// 当 T 当为 interface{} 时，走到这里
		var vv reflect.Value
		if lv.Type() == glua.LTNumber {
			vv = reflect.ValueOf(float64(lv.(glua.LNumber)))
		} else if lv.Type() == glua.LTBool {
			vv = reflect.ValueOf(bool(lv.(glua.LBool)))
		} else if lv.Type() == glua.LTString {
			vv = reflect.ValueOf(string(lv.(glua.LString)))
		} else {
			return fmt.Errorf("can't convert lua value %v to go type %v", lv, rt)
		}
		if !vv.CanConvert(rt) {
			return fmt.Errorf("can't convert lua value %v to go type %v", lv, rt)
		}
		if !rv.CanSet() {
			return fmt.Errorf("out value is unaccessable")
		}
		rv.Set(vv.Convert(rt))
	}

	if err := UnmarshalTable(lv.(*glua.LTable), v); err != nil {
		return fmt.Errorf("can't convert lua table to type %v, %v", rt, err)
	}
	return nil
}
