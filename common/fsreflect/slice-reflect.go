/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 对 slice 进行 reflect 操作
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
func sliceGetItemPtr(vs reflect.Value, index int) (pv reflect.Value, err error) {
	if vs.Type().Kind() != reflect.Ptr {
		err = errors.New("the first argument must be a slice object pointer")
		return
	}
	if vs.IsNil() {
		err = errors.New("the slice object pointer is nil")
		return
	}
	vs = vs.Elem()
	if index < 0 || index >= vs.Len() {
		err = errors.New("index is out of range.")
		return
	}
	vv := vs.Index(index)

	// 如果元素本身就是指针，则返回本身
	if vv.Type().Kind() == reflect.Ptr {
		pv = vv
	} else if vv.CanAddr() {
		pv = reflect.NewAt(vv.Type(), unsafe.Pointer(vv.UnsafeAddr()))
	} else {
		err = fmt.Errorf("element of slice %v is not accessable", vs.Type())
	}
	return
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 获取 slice 中指定索引处的值
// s 可以是 slice 或 slice 指针
func SliceGet(s interface{}, index int) (v interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("get slice element fail: %s", e)
		}
	}()

	if s == nil {
		err = errors.New("slice object must't be nil.")
		return
	}
	ts := reflect.TypeOf(s)
	vs := reflect.ValueOf(s)
	if ts.Kind() == reflect.Ptr {
		ts = ts.Elem()
		vs = vs.Elem()
	}
	if ts.Kind() != reflect.Slice {
		err = errors.New("the first argument must be a slice object.")
		return
	}

	if index < 0 || index >= vs.Len() {
		err = errors.New("index is out of range.")
		return
	}
	v = vs.Index(index).Interface()
	return
}

// 设置 slice 元素
// cvr 表示，如果追加元素类型不匹配，是否进行强制转换
// ps 必须是 slice 指针
func SliceSet(ps interface{}, index int, value interface{}, cvr bool) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("set slice element fail: %s", e)
		}
	}()

	if ps == nil {
		err = errors.New("the slice pointer argunent must be nil")
		return
	}

	tps := reflect.TypeOf(ps)
	vps := reflect.ValueOf(ps)
	if tps.Kind() != reflect.Ptr {
		err = errors.New("the first argument must be an slice object pointer.")
		return
	}

	ts := tps.Elem()
	vs := vps.Elem()
	if ts.Kind() != reflect.Slice {
		err = errors.New("the first argument must be an slice object pointer.")
		return
	}

	if index < 0 || index >= vs.Len() {
		err = errors.New("index is out of range.")
		return
	}

	tv := ts.Elem()
	var vvalue reflect.Value
	if value == nil {
		vvalue = reflect.Zero(tv)
	} else {
		vvalue = reflect.ValueOf(value)
	}
	if vvalue.Type() != tv {
		if !cvr {
			err = fmt.Errorf("slice item type is %v, value type is %v, they are not match.", tv, vvalue.Type())
			return
		}
		if cv, ok := hardConvert(value, tv); ok {
			vvalue = cv
		} else {
			err = fmt.Errorf("value type is %v, it can't be converted to the slice item type %v", vvalue.Type(), tv)
			return
		}
	}

	old := vs.Index(index)
	pold := reflect.NewAt(tv, unsafe.Pointer(old.UnsafeAddr()))
	pold.Elem().Set(vvalue)
	return
}

// 给 slice 添加元素，ps 必须是 slice 指针
// cvr 表示，如果追加元素类型不匹配，是否进行强制转换
// ps 必须是 slice 指针
func SliceAppend(ps interface{}, value interface{}, cvr bool) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("append fail: %s", e)
		}
	}()

	if ps == nil {
		err = errors.New("slice object pointer must't be nil.")
		return
	}

	tps := reflect.TypeOf(ps)
	vps := reflect.ValueOf(ps)
	if tps.Kind() != reflect.Ptr {
		err = errors.New("the first argument must be a slice object pointer")
		return
	}
	ts := tps.Elem()
	vs := vps.Elem()
	if ts.Kind() != reflect.Slice {
		err = errors.New("the first argument must be a slice object pointer")
		return
	}
	if !vs.IsValid() || vs.IsNil() {
		err = errors.New("the slice pointer point to a nil slice")
		return
	}

	tv := ts.Elem()
	var vvalue reflect.Value
	if value == nil {
		vvalue = reflect.Zero(tv)
	} else {
		vvalue = reflect.ValueOf(value)
	}
	if vvalue.Type() == tv {
		vs.Set(reflect.Append(vs, vvalue))
		return
	}
	if !cvr {
		err = fmt.Errorf("slice element type is %v, value type is %v, they are not match.", tv, vvalue.Type())
		return
	}
	if cv, ok := hardConvert(value, tv); ok {
		vs.Set(reflect.Append(vs, cv))
	} else {
		err = fmt.Errorf("the value type is %v, it can't convert to the type of slice element type %v", vvalue.Type(), tv)
	}
	return
}
