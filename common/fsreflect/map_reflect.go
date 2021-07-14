/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 对 map 进程 reflect 操作
@author: fanky
@version: 1.0
@date: 2021-06-21
**/

package fsreflect

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
// 获取值的指针
func mapGetValuePtr(vm reflect.Value, key interface{}) (pv reflect.Value, err error) {
	if vm.Type().Kind() != reflect.Ptr {
		err = errors.New("vm must be a map object pointer")
		return
	}
	if vm.IsNil() {
		err = errors.New("map value is nil")
		return
	}
	vm = vm.Elem()
	tm := vm.Type()

	var vkey reflect.Value
	if key == nil {
		vkey = reflect.Zero(tm.Key())
	} else {
		vkey = reflect.ValueOf(key)
	}
	tkey := vkey.Type()

	var v reflect.Value
	if tkey == tm.Key() {
		v = vm.MapIndex(vkey)
	} else if ck, ok := hardConvert(key, tm.Key()); ok {
		v = vm.MapIndex(ck)
	} else {
		err = fmt.Errorf("type of argument 'key' is %v can't convert to the map key type %v", vkey.Type(), tkey)
		return
	}

	// 如果元素本身就是指针，则返回本身
	if v.Type().Kind() == reflect.Ptr {
		pv = v
	} else if v.CanAddr() {
		pv = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr()))
	} else {
		err = fmt.Errorf("map value which key=%v in %v is unaccessable", key, vm.Type())
	}
	return
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 在 map 中获取键为 key 的值
// hardcrv 表示是否对 key 进行类型强转
// m 必须是 map 对象或者 map 对象的指针
func MapGet(m interface{}, key interface{}, hardcvr bool) (v interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			v = nil
			err = fmt.Errorf("get map value fail: %s", e)
		}
	}()

	if m == nil {
		err = errors.New("the map object argument can't be a nil value.")
		return
	}
	tm := reflect.TypeOf(m)
	vm := reflect.ValueOf(m)
	if tm.Kind() == reflect.Ptr {
		tm = tm.Elem()
		vm = vm.Elem()
	}
	if tm.Kind() != reflect.Map {
		err = fmt.Errorf("the first arguemnt must be a map object. but not %v", tm)
		return
	}

	var vkey reflect.Value
	if key == nil {
		vkey = reflect.Zero(tm.Key())
	} else {
		vkey = reflect.ValueOf(key)
	}
	tkey := vkey.Type()

	if tkey == tm.Key() {
		v = vm.MapIndex(vkey).Interface()
		return
	}
	if !hardcvr {
		err = fmt.Errorf("the map key type is %v but not %v", tkey, vkey.Type())
		return
	}
	if ck, yes := hardConvert(key, tm.Key()); yes {
		v = vm.MapIndex(ck).Interface()
	} else {
		err = fmt.Errorf("type of argument 'key' is %v can't convert to the map key type %v", vkey.Type(), tkey)
	}
	return
}

// 设置 map 中指定 key 处的值
// m 必须是 map 对象，或者是 map 对象的指针
func MapSet(m interface{}, key, value interface{}, kcvr, vcvr bool) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("set map value fail: %s", e)
		}
	}()

	if m == nil {
		err = errors.New("map object is not allow to be a nil value.")
		return
	}

	tm := reflect.TypeOf(m)
	vm := reflect.ValueOf(m)
	if tm.Kind() == reflect.Ptr {
		tm = tm.Elem()
		vm = vm.Elem()
	}
	if tm.Kind() != reflect.Map {
		err = fmt.Errorf("the first arguemnt must be a map object. but not %v", tm)
		return
	}

	tkey := tm.Key()
	var vkey reflect.Value
	if key == nil {
		vkey = reflect.Zero(tkey)
	} else {
		vkey = reflect.ValueOf(key)
	}
	// 传入的 key 值类型不正确
	if vkey.Type() != tkey {
		if !kcvr { // 不允许对 key 进行强转
			err = fmt.Errorf("map key type is %v but not %v", tkey, vkey.Type())
			return
		}
		if ck, ok := hardConvert(key, tkey); ok {
			vkey = ck
		} else {
			err = fmt.Errorf("type of argument 'key' is %v can't convert to the map key type %v", vkey.Type(), tkey)
			return
		}
	}

	tvalue := tm.Elem()
	var vvalue reflect.Value
	if value == nil {
		vvalue = reflect.Zero(tvalue)
	} else {
		vvalue = reflect.ValueOf(value)
	}
	// 传入的 value 值类型不正确
	if vvalue.Type() != tvalue {
		if !vcvr { // 不允许对 value 进行强转
			err = fmt.Errorf("map value type is %v but not %v", tvalue, vvalue.Type())
			return err
		}
		if cv, ok := hardConvert(value, tvalue); ok {
			vvalue = cv
		} else {
			err = fmt.Errorf("type of argument 'value' is %v can't convert to the map value type %v", vvalue.Type(), tvalue)
			return
		}
	}
	vm.SetMapIndex(vkey, vvalue)
	return
}
