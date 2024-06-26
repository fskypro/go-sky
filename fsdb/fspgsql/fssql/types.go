/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: null types
@author: fanky
@version: 1.0
@date: 2024-03-07
**/

package fssql

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// -------------------------------------------------------------------
// String
// -------------------------------------------------------------------
type String string

func (this *String) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case string:
		*this = String(value.(string))
	case []byte:
		*this = String(string(value.([]byte)))
	case time.Time:
		*this = String(value.(time.Time).Format(time.RFC3339Nano))
	default:
		return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
	}
	return nil
}

// -------------------------------------------------------------------
// Int
// -------------------------------------------------------------------
type Int int

func (this *Int) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case int64:
		*this = Int(value.(int64))
		return nil
	}
	return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
}

// ---------------------------------------------------------
type Int8 int8

func (this *Int8) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case int64:
		*this = Int8(value.(int64))
		return nil
	}
	return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
}

// ---------------------------------------------------------
type Int16 int16

func (this *Int16) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case int64:
		*this = Int16(value.(int64))
		return nil
	}
	return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
}

// ---------------------------------------------------------
type Int32 int32

func (this *Int32) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case int64:
		*this = Int32(value.(int64))
		return nil
	}
	return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
}

// ---------------------------------------------------------
type Int64 int64

func (this *Int64) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case int64:
		*this = Int64(value.(int64))
		return nil
	}
	return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
}

// -------------------------------------------------------------------
// Float
// -------------------------------------------------------------------
type Float32 float32

func (this *Float32) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case float64:
		*this = Float32(value.(float64))
	case []uint8:
		str := string(value.([]uint8))
		v, err := strconv.ParseFloat(str, 64)
		if err == nil {
			*this = Float32(v)
			return nil
		}
	}
	return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
}

// ---------------------------------------------------------
type Float64 float64

func (this *Float64) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case float64:
		*this = Float64(value.(float64))
	case []uint8:
		str := string(value.([]uint8))
		v, err := strconv.ParseFloat(str, 64)
		if err == nil {
			*this = Float64(v)
			return nil
		}
	}
	return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
}

// -------------------------------------------------------------------
// Bool
// -------------------------------------------------------------------
type Bool bool

func (this *Bool) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case bool:
		*this = Bool(value.(bool))
	case int64:
		*this = value.(int64) != 0
	default:
		return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
	}
	return nil
}

// -------------------------------------------------------------------
// Time
// -------------------------------------------------------------------
type DateTime time.Time

func (this *DateTime) Scan(value any) error {
	if value == nil { return nil }
	switch value.(type) {
	case time.Time:
		*this = DateTime(value.(time.Time))
		return nil
	}
	return fmt.Errorf("can't convert type %v value to %v", reflect.TypeOf(value), reflect.TypeOf(*this))
}
