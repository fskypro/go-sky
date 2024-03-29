/**
@copyright: fantasysky 2016
@brief: error chain
@author: fanky
@version: 1.0
@date: 2020-05-01
**/

package fserror

import (
	"fmt"

	"fsky.pro/fsdef"
)

// 判断传入的 error 类型参数是否为模板中指定的错误类型
// 注意：“继承” 类型的错误，譬如：
//
//	type Error struct{}
//	type SubError struct{ Error }
//	IsError[Error](SubError{})     // false
func IsError[E error](e error) (same bool) {
	if e == nil { return false }
	defer func() { recover() }()
	var ee E = e.(E)
	func(E) {}(ee)
	same = true
	return
}

// -------------------------------------------------------------------
// I_Error
// -------------------------------------------------------------------
type I_Error interface {
	Unwrap() error       // 返回包装的底层错误
	ErrorChain() []error // 获取错误链，第一个元素为最底层发出的第一个错误
	Error() string       // 格式化输出错误链
	LatestError() string // 最后一个错误
}

// -------------------------------------------------------------------
// S_Error
// -------------------------------------------------------------------
// WarpOutput 表示用 Error 方法获取的错误内容，是否换行和缩进
//
//	默认为 true，如果设置为 false，则用“|”分隔
type S_Error struct {
	msg  string
	werr error
}

func Newf(msg string, args ...interface{}) *S_Error {
	return &S_Error{
		msg: fmt.Sprintf(msg, args...),
	}
}

func Wrapf(err error, msg string, args ...interface{}) *S_Error {
	inst := Newf(msg, args...)
	inst.werr = err
	return inst
}

// ---------------------------------------------------------
func (this *S_Error) Unwrap() error {
	return this.werr
}

// 获取错误链条，返回错误数组中，第一个错误为最早出现的错误
func (this *S_Error) ErrorChain() []error {
	var errs []error
	if e, ok := this.werr.(*S_Error); ok {
		errs = e.ErrorChain()
	} else {
		errs = []error{this.werr}
	}
	return append(errs, this)
}

func (this *S_Error) LatestError() string {
	return this.msg
}

func (this *S_Error) Error() string {
	space := "  "
	msg := this.msg
	werr := this.werr
	for {
		if werr == nil {
			break
		}
		if err, ok := werr.(I_Error); ok {
			msg = msg + fsdef.Endline + space + err.LatestError()
			werr = err.Unwrap()
			space = space + "  "
		} else {
			msg = msg + fsdef.Endline + space + werr.Error()
			break
		}
	}
	return msg
}
