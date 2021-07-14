/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: errors
@author: fanky
@version: 1.0
@date: 2021-02-19
**/

package fsreflect

import (
	"fmt"
	"reflect"
)

// -------------------------------------------------------------------
// 域名字不存在错误
// -------------------------------------------------------------------
type S_FieldError struct {
	stype reflect.Type // 结构体类型名称
	fname string       // 域名称
}

func newFieldError(stype reflect.Type, fname string) *S_FieldError {
	return &S_FieldError{
		stype: stype,
		fname: fname,
	}
}

func (this *S_FieldError) Error() string {
	return fmt.Sprintf("field name %q is not exists in type %v", this.fname, this.stype)
}

// -------------------------------------------------------------------
// 值类型错误
// -------------------------------------------------------------------
type S_ValueError struct {
	stype reflect.Type // 结构体名称
	fname string       // 域名称
	ftype reflect.Type // 域类型名称
	vtype reflect.Type // 设置给域的值
}

func newValueError(stype reflect.Type, fname string, ftype, vtype reflect.Type) *S_ValueError {
	return &S_ValueError{
		stype: stype, // 结构体类型
		fname: fname, // 域名
		ftype: ftype, // 域类型
		vtype: vtype, // 新值类型
	}
}

func (this *S_ValueError) Error() string {
	return fmt.Sprintf("the type of field %q in %q is a %v, but not a %v.", this.fname, this.stype, this.ftype, this.vtype)
}

// -------------------------------------------------------------------
// 方法不存在错误
// -------------------------------------------------------------------
type S_MethodError struct {
	stype reflect.Type // 结构体类型名称
	fname string       // 域名称
}

func newMethodError(stype reflect.Type, fname string) *S_FieldError {
	return &S_FieldError{
		stype: stype,
		fname: fname,
	}
}

func (this *S_MethodError) Error() string {
	return fmt.Sprintf("method name %q is not exists in type %v", this.fname, this.stype)
}
