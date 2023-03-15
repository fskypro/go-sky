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
func arrayGetItemPtr(varr reflect.Value, index int) (pv reflect.Value, err error) {
	if varr.Type().Kind() != reflect.Ptr {
		err = errors.New("the first argument must be a slice object pointer")
		return
	}
	if varr.IsNil() {
		err = errors.New("the array object pointer is nil")
		return
	}
	varr = varr.Elem()
	if index < 0 || index >= varr.Len() {
		err = errors.New("index is out of range.")
		return
	}
	vv := varr.Index(index)

	// 如果元素本身就是指针，则返回本身
	if vv.Type().Kind() == reflect.Ptr {
		pv = vv
	} else if vv.CanAddr() {
		pv = reflect.NewAt(vv.Type(), unsafe.Pointer(vv.UnsafeAddr()))
	} else {
		err = fmt.Errorf("element of array %v is not accessable", varr.Type())
	}
	return
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 获取 array 中指定索引处的值
func ArrayGet(arr interface{}, index int) (v interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("get array element fail: %s", e)
		}
	}()

	if arr == nil {
		err = errors.New("array object must't be nil.")
		return
	}
	if reflect.TypeOf(arr).Kind() != reflect.Array {
		err = errors.New("the first argument must be an array object.")
		return
	}

	va := reflect.ValueOf(arr)
	if index < 0 || index >= va.Len() {
		err = errors.New("index is out of range.")
		return
	}
	v = va.Index(index).Interface()
	return
}

// 设置 array 元素
// cvr 表示，如果追加元素类型不匹配，是否进行强制转换
func ArraySet(parr interface{}, index int, value interface{}, cvr bool) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("set array element fail: %s", e)
		}
	}()

	if parr == nil {
		err = errors.New("the array pointer argunent must be nil")
		return
	}

	tparr := reflect.TypeOf(parr)
	vparr := reflect.ValueOf(parr)
	if tparr.Kind() != reflect.Ptr {
		err = errors.New("the first argument must be an array object pointer.")
		return
	}

	tarr := tparr.Elem()
	varr := vparr.Elem()
	if tarr.Kind() != reflect.Array {
		err = errors.New("the first argument must be an array object pointer.")
		return
	}

	if index < 0 || index >= varr.Len() {
		err = errors.New("index is out of range.")
		return
	}

	tv := tarr.Elem()
	var vvalue reflect.Value
	if value == nil {
		vvalue = reflect.Zero(tv)
	} else {
		vvalue = reflect.ValueOf(value)
	}
	if vvalue.Type() != tv {
		if !cvr {
			err = fmt.Errorf("array item type is %v, value type is %v, they are not match.", tv, vvalue.Type())
			return
		}
		if cv, ok := hardConvert(value, tv); ok {
			vvalue = cv
		} else {
			err = fmt.Errorf("value type is %v, it can't be converted to the array item type %v", vvalue.Type(), tv)
			return
		}
	}

	old := varr.Index(index)
	pold := reflect.NewAt(tv, unsafe.Pointer(old.UnsafeAddr()))
	pold.Elem().Set(vvalue)
	return
}
