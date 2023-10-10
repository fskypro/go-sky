/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: errors
@author: fanky
@version: 1.0
@date: 2023-10-03
**/

package fsicmp

import "errors"

var (
	// 超时错误
	ErrTimeout = errors.New("timeout")

	// 脏数据错误
	ErrInvalidPackage = errors.New("dirty icmp package")
)
