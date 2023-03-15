/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: types definations
@author: fanky
@version: 1.0
@date: 2022-07-29
**/

package fstype

// -------------------------------------------------------------------
// 类型定义
// -------------------------------------------------------------------
// 数值类型
type T_Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// 整数
type T_IUNumber interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

// 有符号整数
type T_INumber interface {
	int | int8 | int16 | int32 | int64
}

// 无符号整数
type T_UNumber interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// 浮点数
type T_FNumber interface {
	float32 | float64
}
