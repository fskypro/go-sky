/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: string convertor
@author: fanky
@version: 1.0
@date: 2021-06-12
**/

package fsconv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"fsky.pro/fserror"
)

var intSize = int(unsafe.Sizeof(int(0)) * 8)

// -------------------------------------------------------------------
// base types
// -------------------------------------------------------------------
type T_Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type T_Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type T_Float interface{ ~float32 | ~float64 }

type T_Base interface {
	~string | T_Signed | T_Unsigned | T_Float | ~bool
}

// -------------------------------------------------------------------
// reflect types
// -------------------------------------------------------------------
var strType = reflect.TypeOf("").Kind()

var intType = reflect.TypeOf(int(0)).Kind()
var int8Type = reflect.TypeOf(int8(0)).Kind()
var int16Type = reflect.TypeOf(int16(0)).Kind()
var int32Type = reflect.TypeOf(int32(0)).Kind()
var int64Type = reflect.TypeOf(int64(0)).Kind()

var uintType = reflect.TypeOf(uint(0)).Kind()
var uint8Type = reflect.TypeOf(uint8(0)).Kind()
var uint16Type = reflect.TypeOf(uint16(0)).Kind()
var uint32Type = reflect.TypeOf(uint32(0)).Kind()
var uint64Type = reflect.TypeOf(uint64(0)).Kind()

var float32Type = reflect.TypeOf(float32(0)).Kind()
var float64Type = reflect.TypeOf(float64(0)).Kind()

var boolType = reflect.TypeOf(true).Kind()

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
type RangeError struct {
	StrNum string
}

func (self RangeError) Error() string {
	return fmt.Sprintf("fsconv.StrTo: parsing %q: value out of range", self.StrNum)
}

type SyntaxError struct {
	StrNum string
}

func (self SyntaxError) Error() string {
	return fmt.Sprintf("fsconv.StrTo: parsing %q: invalid syntax", self.StrNum)
}

func StrTo[T T_Base](str string) (tv T, err error) {
	defer func() {
		e, ok := fserror.AsError[*strconv.NumError](err)
		if ok && errors.Is(e, strconv.ErrRange) {
			err = RangeError{str}
		} else if ok && errors.Is(e, strconv.ErrSyntax) {
			err = SyntaxError{str}
		}
	}()

	var v any
	switch any(tv).(type) {
	case string:
		return any(str).(T), nil

	// -------------------------------------------
	case int:
		vv, e := strconv.ParseInt(str, 10, intSize)
		v, err = int(vv), e
	case int8:
		vv, e := strconv.ParseInt(str, 10, 8)
		v, err = int8(vv), e
	case int16:
		vv, e := strconv.ParseInt(str, 10, 16)
		v, err = int16(vv), e
	case int32:
		vv, e := strconv.ParseInt(str, 10, 32)
		v, err = int32(vv), e
	case int64:
		v, err = strconv.ParseInt(str, 10, 64)

	// -------------------------------------------
	case uint:
		vv, e := strconv.ParseUint(str, 10, intSize)
		v, err = uint(vv), e
	case uint8:
		vv, e := strconv.ParseUint(str, 10, 8)
		v, err = uint8(vv), e
	case uint16:
		vv, e := strconv.ParseUint(str, 10, 16)
		v, err = uint16(vv), e
	case uint32:
		vv, e := strconv.ParseUint(str, 10, 32)
		v, err = uint32(vv), e
	case uint64:
		v, err = strconv.ParseUint(str, 10, 64)

	// -------------------------------------------
	case float32:
		vv, e := strconv.ParseFloat(str, 32)
		v, err = float32(vv), e
	case float64:
		v, err = strconv.ParseFloat(str, 64)

	// -------------------------------------------
	case bool:
		v, err = strconv.ParseBool(str)

	// -------------------------------------------
	default:
		goto L
	}
	return

	// ----------------------------------------------------
L:
	t := reflect.TypeOf(tv)
	switch t.Kind() {
	case strType:
		v = str

	// -------------------------------------------
	case intType:
		v, err = strconv.ParseInt(str, 10, intSize)
	case int8Type:
		v, err = strconv.ParseInt(str, 10, 8)
	case int16Type:
		v, err = strconv.ParseInt(str, 10, 16)
	case int32Type:
		v, err = strconv.ParseInt(str, 10, 32)
	case int64Type:
		v, err = strconv.ParseInt(str, 10, 64)

	// -------------------------------------------
	case uintType:
		v, err = strconv.ParseUint(str, 10, intSize)
	case uint8Type:
		v, err = strconv.ParseUint(str, 10, 8)
	case uint16Type:
		v, err = strconv.ParseUint(str, 10, 16)
	case uint32Type:
		v, err = strconv.ParseUint(str, 10, 32)
	case uint64Type:
		v, err = strconv.ParseUint(str, 10, 64)

	// -------------------------------------------
	case float32Type:
		v, err = strconv.ParseFloat(str, 32)
	case float64Type:
		v, err = strconv.ParseFloat(str, 64)

	// -------------------------------------------
	case boolType:
		v, err = strconv.ParseBool(str)
	}
	if err != nil {
		return
	}
	tv = reflect.ValueOf(v).Convert(t).Interface().(T)
	return
}

func StrToAnyType[T T_Base](str string) (any, error) {
	return StrTo[T](str)
}
