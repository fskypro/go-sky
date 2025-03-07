/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: string convertor
@author: fanky
@version: 1.0
@date: 2021-06-12
**/

package fsstrconv

import (
	"strconv"
	"unsafe"
)

type T_Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type T_Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
type T_Float interface{ ~float32 | ~float64 }

// -------------------------------------------------------------------
// convertor interface
// -------------------------------------------------------------------
type i_Convertor interface {
	convert(string) (any, error)
}

var convertors = map[any]i_Convertor{}

func addConvertor[T T_Signed | T_Unsigned | T_Float | bool](convertor i_Convertor) {
	var v T
	convertors[v] = convertor
}

// -------------------------------------------------------------------
// convertors
// -------------------------------------------------------------------
type s_Str2Int[T T_Signed] struct{}

func (self s_Str2Int[T]) convert(str string) (any, error) {
	v, err := strconv.ParseInt(str, 10, int(unsafe.Sizeof(uint(0)))*8)
	return T(v), err
}

type s_Str2Uint[T T_Unsigned] struct{}

func (self s_Str2Uint[T]) convert(str string) (any, error) {
	v, err := strconv.ParseUint(str, 10, int(unsafe.Sizeof(uint(0)))*8)
	return T(v), err
}

type s_Str2Float[T T_Float] struct{}

func (self s_Str2Float[T]) convert(str string) (any, error) {
	v, err := strconv.ParseFloat(str, 32)
	return T(v), err
}

// ---------------------------------------------------------
// 为 true 的字符串：true、TRUE、T、t、1
// 其他全为 false
type s_Str2Bool struct{}

func (self s_Str2Bool) convert(str string) (any, error) {
	return strconv.ParseBool(str)
}

func init() {
	addConvertor[int](s_Str2Int[int]{})
	addConvertor[int8](s_Str2Int[int8]{})
	addConvertor[int16](s_Str2Int[int16]{})
	addConvertor[int32](s_Str2Int[int32]{})
	addConvertor[int64](s_Str2Int[int64]{})

	addConvertor[uint](s_Str2Uint[uint]{})
	addConvertor[uint8](s_Str2Uint[uint8]{})
	addConvertor[uint16](s_Str2Uint[uint16]{})
	addConvertor[uint32](s_Str2Uint[uint32]{})
	addConvertor[uint64](s_Str2Uint[uint64]{})

	addConvertor[float32](s_Str2Float[float32]{})
	addConvertor[float64](s_Str2Float[float64]{})

	addConvertor[bool](s_Str2Bool{})
}

// -------------------------------------------------------------------
func StrTo[T T_Signed | T_Unsigned | T_Float | bool](str string) (T, error) {
	var tv T
	v, err := convertors[tv].convert(str)
	return v.(T), err
}

func StrToAnyType[T T_Signed | T_Unsigned | T_Float | bool](str string) (any, error) {
	return StrTo[T](str)
}
