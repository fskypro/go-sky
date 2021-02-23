/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: reflect utils
@author: fanky
@version: 1.0
@date: 2021-02-19
**/

// ------------------------------------------------------------
// GetFieldValue（获取结构体的私有字段值）
// SetFieldValue（设置结构体的私有字段值）
//
// 获取包的全局变量值、调用包的私有函数、调用结构体的私有函数，不需要
// 反射，可以用 //go:linkname 指令
// ------------------------------------------------------------

package fsreflect

import (
	"reflect"
	"unsafe"
)

// 获取结构体字段值，包括私有字段
func GetFieldValue(obj interface{}, fname string) (fv interface{}, err error) {
	v := reflect.ValueOf(obj)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var field reflect.Value
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name == fname {
			field = v.Field(i)
			break
		}
	}

	if !field.IsValid() {
		err = newFieldError(v.Type().Name(), fname)
		return
	}

	if field.CanInterface() {
		fv = field.Interface()
		return
	}

	// 成员为指针类型
	if field.Type().Kind() == reflect.Ptr {
		if field.IsNil() {
			return
		}
		up := field.Elem().UnsafeAddr()
		vv := reflect.NewAt(field.Elem().Type(), unsafe.Pointer(up))
		fv = vv.Interface()
	} else {
		up := field.UnsafeAddr()
		vv := reflect.NewAt(field.Type(), unsafe.Pointer(up))
		fv = vv.Elem().Interface()
	}

	return
}

// 设置结构体字段值，包括私有字段
func SetFieldValue(obj interface{}, fname string, fv interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var field reflect.Value
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Name == fname {
			field = v.Field(i)
			break
		}
	}

	// 域名不存在
	if !field.IsValid() {
		return newFieldError(v.Type().Name(), fname)
	}

	vv := reflect.ValueOf(fv)
	// 指定的域类型为指针类型
	if field.Type().Kind() == reflect.Ptr {
		// 传入值为 nil
		if fv == nil {
			vv = reflect.Zero(field.Type()) // 创建 nil 值
		} else if vv.Type() != field.Type() { // 设置了不同类型的值
			return newValueError(v.Type().Name(), fname, field.Type().Name(), vv.Type().Name())
		}
		if field.CanSet() {
			field.Set(vv)
		} else {
			upp := unsafe.Pointer(field.UnsafeAddr()) // 指针的指针
			ppv := reflect.NewAt(field.Type(), upp)   // 创建指针的指针对象（该指针的指针，指向 field 原来的位置）
			ppv.Elem().Set(vv)                        // 将指针的指针，指向的位置（即 field 的内存位置）修改为新的值
		}
	} else {
		if fv == nil {
			return newValueError(v.Type().Name(), fname, field.Type().Name(), "nil")
		} else if field.Type() != vv.Type() {
			return newValueError(v.Type().Name(), fname, field.Type().Name(), vv.Type().Name())
		}
		if field.CanSet() {
			field.Set(vv)
		} else {
			up := unsafe.Pointer(field.UnsafeAddr())
			pv := reflect.NewAt(field.Type(), up) // 创建一个新的指针，指向原来 field 的位置
			pv.Elem().Set(vv)
		}
	}
	return nil
}
