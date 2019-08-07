/**
@copyright: fantasysky 2016
@brief: 实现错误相关的功能函数
@author: fanky
@version: 1.0
@date: 2019-01-02
**/

package fserror

import "fmt"
import "errors"

// StrErrorf 新建一个格式化字符串的错误
func StrErrorf(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}
