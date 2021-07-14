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
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// ---------------------------------------------------------
func Str2Int(str string) (int, error) {
	v64, err := strconv.ParseInt(str, 10, int(unsafe.Sizeof(uint(0)))*8)
	return int(v64), err
}

func Str2Int8(str string) (int8, error) {
	v64, err := strconv.ParseInt(str, 10, 8)
	return int8(v64), err
}

func Str2Int16(str string) (int16, error) {
	v64, err := strconv.ParseInt(str, 10, 16)
	return int16(v64), err
}

func Str2Int32(str string) (int32, error) {
	v64, err := strconv.ParseInt(str, 10, 32)
	return int32(v64), err
}

func Str2Int64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

// ---------------------------------------------------------
func Str2Uint(str string) (uint, error) {
	v64, err := strconv.ParseUint(str, 10, int(unsafe.Sizeof(uint(0)))*8)
	return uint(v64), err
}

func Str2Uint8(str string) (uint8, error) {
	v64, err := strconv.ParseUint(str, 10, 8)
	return uint8(v64), err
}

func Str2Uint16(str string) (uint16, error) {
	v64, err := strconv.ParseUint(str, 10, 16)
	return uint16(v64), err
}

func Str2Uint32(str string) (uint32, error) {
	v64, err := strconv.ParseUint(str, 10, 32)
	return uint32(v64), err
}

func Str2Uint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

// ---------------------------------------------------------
func Str2Float32(str string) (float32, error) {
	v64, err := strconv.ParseFloat(str, 32)
	return float32(v64), err
}

func Str2Float64(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

// ---------------------------------------------------------
// 为 true 的字符串：true、TRUE、T、t、1
// 其他全为 false
func Str2Bool(str string) bool {
	v, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return v
}

// -------------------------------------------------------------------
// 将 str 转换为与 v 类型一致的值
// 注意：v 必须为基础类型值
func Str2TypeOf(str string, v interface{}) (interface{}, error) {
	switch v.(type) {
	case string:
		return str, nil
	case []byte:
		return []byte(str), nil
	case []rune:
		return []rune(str), nil
	case int:
		return Str2Int(str)
	case int8:
		return Str2Int8(str)
	case int16:
		return Str2Int16(str)
	case int32:
		return Str2Int32(str)
	case int64:
		return Str2Int64(str)
	case uint:
		return Str2Uint(str)
	case uint8:
		return Str2Uint8(str)
	case uint16:
		return Str2Uint16(str)
	case uint32:
		return Str2Uint32(str)
	case uint64:
		return Str2Uint64(str)
	case float32:
		return Str2Float32(str)
	case float64:
		return Str2Float64(str)
	case bool:
		return Str2Bool(str), nil

	// ------------------------------
	case *string:
		return &str, nil
	case *int:
		vv, err := Str2Int(str)
		if err == nil {
			return &vv, err
		}
		return nil, err
	case *int8:
		vv, err := Str2Int8(str)
		if err == nil {
			return &vv, err
		}
		return nil, err
	case *int16:
		vv, err := Str2Int16(str)
		if err == nil {
			return &vv, err
		}
		return nil, err
	case *int32:
		vv, err := Str2Int32(str)
		if err == nil {
			return &vv, err
		}
		return nil, err
	case *int64:
		vv, err := Str2Int64(str)
		if err == nil {
			return &vv, err
		}
		return nil, err
	case *uint:
		vv, err := Str2Uint(str)
		if err != nil {
			return &vv, err
		}
		return nil, err
	case *uint8:
		vv, err := Str2Uint8(str)
		if err != nil {
			return &vv, err
		}
		return nil, err
	case *uint16:
		vv, err := Str2Uint16(str)
		if err != nil {
			return &vv, err
		}
		return nil, err
	case *uint32:
		vv, err := Str2Uint32(str)
		if err != nil {
			return &vv, err
		}
		return nil, err
	case *uint64:
		vv, err := Str2Uint64(str)
		if err != nil {
			return &vv, err
		}
		return nil, err
	case *float32:
		vv, err := Str2Float32(str)
		if err != nil {
			return &vv, err
		}
		return nil, err
	case *float64:
		vv, err := Str2Float64(str)
		if err != nil {
			return &vv, err
		}
		return nil, err
	case *bool:
		vv := Str2Bool(str)
		return &vv, nil
	}
	return nil, fmt.Errorf("string %q can't convert to type %v", str, reflect.TypeOf(v))
}
