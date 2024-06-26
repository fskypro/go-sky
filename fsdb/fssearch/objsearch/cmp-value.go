/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: comare value
@author: fanky
@version: 1.0
@date: 2024-05-30
**/

package objsearch

import (
	"reflect"
	"strconv"
	"time"

	"fsky.pro/fstype"
)

var (
	refString  reflect.Type
	refInt64   reflect.Type
	refUint64  reflect.Type
	refFloat64 reflect.Type
	refTime    reflect.Type
)

func init() {
	var str string
	refString = reflect.TypeOf(str)
	var i64 int64
	refInt64 = reflect.TypeOf(i64)
	var u64 uint64
	refUint64 = reflect.TypeOf(u64)
	var f64 float64
	refFloat64 = reflect.TypeOf(f64)
	var t time.Time
	refTime = reflect.TypeOf(t)
}

// -------------------------------------------------------------------
// CmpValue
// -------------------------------------------------------------------
type s_CmpValue struct {
	value any
	vtype reflect.Type

	_pString  *string
	_pInt64   *int64
	_pUint64  *uint64
	_pFluat64 *float64
	_pTime    *time.Time
	_array    []*s_CmpValue
}

func newCmpValue(value any) *s_CmpValue {
	return &s_CmpValue{
		value: value,
		vtype: reflect.TypeOf(value),
	}
}

// 转换为字符串
func (this *s_CmpValue) asString() (string, bool) {
	if this._pString != nil {
		return *this._pString, true
	}
	if fstype.IsType[string](this.value) {
		value := this.value.(string)
		this._pString = &value
		return value, true
	}
	return "", false
}

// 转换为整形
func (this *s_CmpValue) asInt64() (int64, bool) {
	if this._pInt64 != nil {
		return *this._pInt64, true
	}

	switch this.value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		var v int64
		rv := reflect.ValueOf(this.value).Convert(reflect.TypeOf(v))
		value := rv.Interface().(int64)
		this._pInt64 = &value
		return value, true
	case string:
		v, err := strconv.ParseInt(this.value.(string), 10, 64)
		if err == nil {
			this._pInt64 = &v
			return v, true
		}
	}
	return 0, false
}

// 转换为无符号整形
func (this *s_CmpValue) asUint64() (uint64, bool) {
	if this._pUint64 != nil {
		return *this._pUint64, true
	}

	switch this.value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		rv := reflect.ValueOf(this.value).Convert(refUint64)
		value := rv.Interface().(uint64)
		this._pUint64 = &value
		return value, true
	case string:
		v, err := strconv.ParseUint(this.value.(string), 10, 64)
		if err == nil {
			this._pUint64 = &v
			return v, true
		}
	}
	return 0, false
}

// 转换为浮点型
func (this *s_CmpValue) asFloat64() (float64, bool) {
	if this._pFluat64 != nil {
		return *this._pFluat64, true
	}

	switch this.value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		rv := reflect.ValueOf(this.value).Convert(refFloat64)
		value := rv.Interface().(float64)
		this._pFluat64 = &value
		return value, true
	case string:
		v, err := strconv.ParseFloat(this.value.(string), 64)
		if err == nil {
			this._pFluat64 = &v
			return v, true
		}
	}
	return 0, false
}

// 转换为时间
func (this *s_CmpValue) asTime() (ti time.Time, ok bool) {
	if this._pTime != nil {
		return *this._pTime, true
	}

	if !fstype.IsType[string](this.value) {
		return
	}
	t, err := time.ParseInLocation(time.DateTime, this.value.(string), time.Local)
	if err != nil {
		t, err = time.ParseInLocation(time.DateOnly, this.value.(string), time.Local)
	}
	if err != nil {
		return
	} else {
		this._pTime = &t
	}
	return t, true
}

// 转换为数组
func (this *s_CmpValue) asArray() ([]*s_CmpValue, bool) {
	if this._array != nil {
		return this._array, true
	}
	if !fstype.IsType[[]any](this.value) {
		return nil, false
	}
	this._array = make([]*s_CmpValue, 0)
	for _, item := range this.value.([]any) {
		this._array = append(this._array, newCmpValue(item))
	}
	return this._array, true
}
