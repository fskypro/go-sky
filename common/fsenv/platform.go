/**
@copyright: fantasysky 2016
@brief: 平台相关功能函数
@author: fanky
@version: 1.0
@date: 2019-06-25
**/

package fsenv

import "unsafe"

// -------------------------------------------------------------------
// 大小端机判别
// -------------------------------------------------------------------
// 判断当前平台是否是小端字节序
func IsLittleEndian() bool {
	const INT_SIZE = int(unsafe.Sizeof(0))
	var i int = 1
	intBytes := (*[INT_SIZE]byte)(unsafe.Pointer(&i))
	return intBytes[0] > 0
}

// 判断当前平台是否是大端字节序
func IsBigEndian() bool {
	return !IsLittleEndian()
}
