/**
@copyright: fantasysky 2016
@brief: error chain
@author: fanky
@version: 1.0
@date: 2020-05-01
**/

package fserror

import (
	"errors"
)

// 判断传入的 error 类型参数是否为模板中指定的错误类型
// type E1 struct { error }
// type E2 struct { error }
// errs := errors.Join(E1{errows.New("xxxx")}, E2{errors.New("yyyy")})
// fserrors.IsError[E1](err1)  // true
// fserrors.IsError[E2](err2)  // true
// e1 = E1{ errors.New("zzz") }
// fserrors.IsError[E1](e1)    // true
func IsError[E error](err error) bool {
	var e E
	return errors.As(err, &e)
}


// 把指定错误还原或转换为原错误类型对象
func AsError[E error](err error) (E, bool) {
	var e E
	ok := errors.As(err, &e)
	return e, ok
}
