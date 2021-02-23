/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: errors
@author: fanky
@version: 1.0
@date: 2021-02-19
**/

package fsreflect

import "fmt"

// -------------------------------------------------------------------
// 域名字不存在错误
// -------------------------------------------------------------------
type S_FieldError struct {
	stype string // 结构体类型名称
	fname string // 域名称
}

func newFieldError(stype, fname string) *S_FieldError {
	return &S_FieldError{
		stype: stype,
		fname: fname,
	}
}

func (this *S_FieldError) Error() string {
	return fmt.Sprintf("field name %q is not exists in type %q", this.fname, this.stype)
}

// -------------------------------------------------------------------
// 值类型错误
// -------------------------------------------------------------------
type S_ValueError struct {
	stype string // 结构体名称
	fname string // 域名称
	ftype string // 域类型名称
	vtype string // 设置给域的值
}

func newValueError(stype, fname, ftype, vtype string) *S_ValueError {
	return &S_ValueError{
		stype: stype,
		fname: fname,
		ftype: ftype,
		vtype: vtype,
	}
}

func (this *S_ValueError) Error() string {
	return fmt.Sprintf("the type of field %q in %q is %q, but not %q.", this.fname, this.stype, this.ftype, this.vtype)
}

// -------------------------------------------------------------------
// 方法不存在错误
// -------------------------------------------------------------------
type S_MethodError struct {
	stype string // 结构体类型名称
	fname string // 域名称
}

func newMethodError(stype, fname string) *S_FieldError {
	return &S_FieldError{
		stype: stype,
		fname: fname,
	}
}

func (this *S_MethodError) Error() string {
	return fmt.Sprintf("method name %q is not exists in type %q", this.fname, this.stype)
}
