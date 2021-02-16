/**
@copyright: fantasysky 2016
@brief: 值接口
@author: fanky
@version: 1.0
@date: 2019-05-30
**/

package freejson

import "bufio"

// -------------------------------------------------------------------
// json value interface
// -------------------------------------------------------------------
type I_Value interface {
	Type() JType
	Name() string

	AsNull() *S_Null
	AsObject() *S_Object
	AsList() *S_List
	AsString() *S_String
	AsInt64() *S_Int64
	AsUInt64() *S_UInt64
	AsFloat64() *S_Float64
	AsBool() *S_Bool

	WriteTo(*bufio.Writer) (int, error)
	String() string
}

// -------------------------------------------------------------------
// base value
// -------------------------------------------------------------------
type s_Base struct{}

func (s_Base) AsNull() *S_Null {
	return NewNull()
}

func (s_Base) AsObject() *S_Object {
	return NewObject()
}

func (s_Base) AsList() *S_List {
	return NewList()
}

func (s_Base) AsString() *S_String {
	return NewString("")
}

func (s_Base) AsInt64() *S_Int64 {
	return NewInt64(0)
}

func (s_Base) AsUInt64() *S_UInt64 {
	return NewUInt64(0)
}

func (s_Base) AsFloat64() *S_Float64 {
	return NewFloat64(0.0)
}

func (s_Base) AsBool() *S_Bool {
	return NewBool(false)
}
