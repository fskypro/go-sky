/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: sql types
@author: fanky
@version: 1.0
@date: 2025-06-20
**/

package fssql

import "reflect"

// json 类型
type JSONB []byte

var jsonbType = reflect.TypeOf(JSONB{})
