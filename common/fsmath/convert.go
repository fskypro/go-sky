/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 进制转换
@author: fanky
@version: 1.0
@date: 2023-09-15
**/

package fsmath

import (
	"math"
	"strings"

	"fsky.pro/fstype"
)

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 256 以内任何进制数转换为 10 进制数
// 注意：不检查溢出
func DecimalFrom(base uint8, nums []uint8) int64 {
	var result int64 = 0
	exp := len(nums) - 1
	if exp < 0 { return 0 }
	for _, n := range nums {
		result += int64(n+1) * int64(math.Pow(float64(base), float64(exp)))
		exp--
	}
	return result - 1
}

// 10 进制数转换为 256 以内的任何进制数
func DecimalTo(base uint8, decimal int64) []uint8 {
	array := []uint8{}
	for decimal >= 0 {
		remainder := decimal % int64(base)
		array = append([]uint8{uint8(remainder)}, array...)
		decimal = decimal/int64(base) - 1
	}
	return array
}

// -------------------------------------------------------------------
// 十六进制数
// -------------------------------------------------------------------
type T_Hex[T fstype.T_IUNumber | fstype.T_UNumber] string

func Hex[T fstype.T_IUNumber | fstype.T_UNumber](decimal T) T_Hex[T] {
	result := ""
	nums := DecimalTo(16, int64(decimal))
	for _, num := range nums {
		if num < 10 {
			result += string('0'+num)
		} else { 
			result += string('a'+num) 
		}
	}
	return T_Hex[T](result)
}

func HEX[T fstype.T_IUNumber | fstype.T_UNumber](decimal T) T_Hex[T] {
	hex := Hex[T](decimal)
	return T_Hex[T](strings.ToUpper(string(hex)))
}

func (self T_Hex[T]) To() T {
	nums := []uint8{}
	for _, c := range []byte( self) {
		if c >= '0' && c <= '9' {
			nums = append(nums, uint8(c-'0'))
		} else if c >= 'A' && c <= 'F' {
			nums = append(nums, c - 'A')
		} else if c >= 'a' && c <= 'f' {
			nums = append(nums, c - 'a')
		}
	}
	return T(DecimalFrom(16, nums))
}
