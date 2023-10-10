/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: string to number or boolean
@author: fanky
@version: 1.0
@date: 2023-09-13
**/

package fsreflect

import (
	"reflect"
	"strconv"
)

func str2Int[T int|int8|int16|int32|int64](str string) (rv reflect.Value, ok bool) {
	v , err := strconv.ParseInt(str, 10, 64)
	ok = err == nil
	if !ok { return }
	rv = reflect.ValueOf(T(v))
	return
}

func str2Uint[T uint|uint8|uint16|uint32|uint64](str string) (rv reflect.Value, ok bool) {
	v , err := strconv.ParseUint(str, 10, 64)
	ok = err == nil
	if !ok { return }
	rv = reflect.ValueOf(T(v))
	return
}

func str2Float[T float32|float64](str string) (rv reflect.Value, ok bool) {
	v , err := strconv.ParseFloat(str, 10)
	ok = err == nil
	if !ok { return }
	rv = reflect.ValueOf(T(v))
	return
}

func str2Bool(str string) (rv reflect.Value, ok bool) {
	v , err := strconv.ParseBool(str)
	ok = err == nil
	if ! ok { return }
	rv = reflect.ValueOf(v)
	return
}

var convertors = map[reflect.Kind]func(string) (reflect.Value, bool) {}

func init() {
	convertors[reflect.Int] = str2Int[int]
	convertors[reflect.Int8] = str2Int[int8]
	convertors[reflect.Int16] = str2Int[int16]
	convertors[reflect.Int32] = str2Int[int32]
	convertors[reflect.Int64] = str2Int[int64]

	convertors[reflect.Uint] = str2Uint[uint]
	convertors[reflect.Uint8] = str2Uint[uint8]
	convertors[reflect.Uint16] = str2Uint[uint16]
	convertors[reflect.Uint32] = str2Uint[uint32]
	convertors[reflect.Uint64] = str2Uint[uint64]

	convertors[reflect.Float32] = str2Float[float32]
	convertors[reflect.Float64] = str2Float[float64]

	convertors[reflect.Bool] = str2Bool
}

func strTo(str string, t reflect.Type) (rv reflect.Value, ok bool) {
	if convertors[t.Kind()] == nil { return }
	rv , yes := convertors[t.Kind()](str)
	if !yes { return }
	rv = rv.Convert(t)
	return
}
