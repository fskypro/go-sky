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
	"fsky.pro/fsenv"
)

// -------------------------------------------------------------------
// I_Error
// -------------------------------------------------------------------
type I_Error interface {
	Error() string
	Unwrap() error       // 返回包装的底层错误
	ErrorChain() []error // 获取错误链，第一个元素为最底层发出的第一个错误
	FmtError() string    // 格式化输出错误链
}

// -------------------------------------------------------------------
// S_Error
// -------------------------------------------------------------------
// WarpOutput 表示用 Error 方法获取的错误内容，是否换行和缩进
//    默认为 true，如果设置为 false，则用“|”分隔
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

func (this *S_Error) Error() string {
	return this.msg
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

func (this *S_Error) FmtError() string {
	var base error
	errs := []I_Error{this}
	werr := this.werr
	for {
		if werr == nil {
			break
		}
		if we, ok := werr.(I_Error); ok {
			errs = append(errs, we)
			werr = we.Unwrap()
		} else {
			base = werr
			break
		}
	}

	var msg string
	space := ""
	for _, we := range errs {
		msg = fmt.Sprintf("%s%s%s%s", msg, fsenv.Endline, space, we.Error())
		space = space + "  "
	}
	if base != nil {
		msg = fmt.Sprintf("%s%s%s%s", msg, fsenv.Endline, space, base.Error())
	}
	return msg
}
