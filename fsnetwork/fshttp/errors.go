/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: errors
@author: fanky
@version: 1.0
@date: 2023-04-27
**/

package fshttp

import (
	"fmt"
	"reflect"
)

// -------------------------------------------------------------------
// 没有传入对应请求参数错误
// -------------------------------------------------------------------
type S_NoReqArgError struct {
	error
	ArgName string
}

func newNoReqArgError(argName string) S_NoReqArgError {
	return S_NoReqArgError{ArgName: argName}
}

func (self S_NoReqArgError) Error() string {
	return fmt.Sprintf("request argument %q is not exists", self.ArgName)
}

// -------------------------------------------------------------------
// 请求参数值类型错误
// -------------------------------------------------------------------
type S_ReqArgTypeError struct {
	error
	ArgName string
	argType reflect.Type
}

func newReqArgTypeError(argName string, argType reflect.Type) S_ReqArgTypeError {
	return S_ReqArgTypeError{ArgName: argName, argType: argType}
}

func (self S_ReqArgTypeError) Error() string {
	return fmt.Sprintf("the value type for request argument %q must be %v", self.ArgName, self.argType.Kind())
}
