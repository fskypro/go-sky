/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: definations
@author: fanky
@version: 1.0
@date: 2024-12-09
**/

package fslua

import (
	"reflect"

	glua "github.com/yuin/gopher-lua"
)

var rtError = reflect.TypeOf((*error)(nil)).Elem()
var rtString = reflect.TypeOf("")
var rtNumnber = reflect.TypeOf(float64(0))
var rtBool = reflect.TypeOf(true)
var rtLValue = reflect.TypeOf((*glua.LValue)(nil)).Elem()
