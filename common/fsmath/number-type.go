/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: number types
@author: fanky
@version: 1.0
@date: 2026-07-07
**/

package fsmath

// 数值类型
type T_Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// 有符号数
type T_Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}
