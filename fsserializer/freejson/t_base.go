/**
@copyright: fantasysky 2016
@brief: 值接口
@author: fanky
@version: 1.0
@date: 2019-05-30
**/

package freejson

import (
	"bufio"
	"fmt"
)

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
type s_Base struct {
	vtype JType
}

func createBase(t JType) s_Base {
	return s_Base{t}
}

func (self s_Base) Type() JType {
	return self.vtype
}

func (self s_Base) Name() string {
	return typeNames[self.vtype]
}

func (self s_Base) AsNull() *S_Null {
	panic(fmt.Sprintf("json value type is a %s but not nil.", self.Name()))
	return NewNull()
}

func (self s_Base) AsObject() *S_Object {
	panic(fmt.Sprintf("json value type is a %s but not an object.", self.Name()))
	return NewObject()
}

func (self s_Base) AsList() *S_List {
	panic(fmt.Sprintf("json value type is a %s but not a list.", self.Name()))
	return NewList()
}

func (self s_Base) AsString() *S_String {
	panic(fmt.Sprintf("json value type is a %s but not a string.", self.Name()))
	return NewString("")
}

func (self s_Base) AsInt64() *S_Int64 {
	panic(fmt.Sprintf("json value type is a %s but not an int64.", self.Name()))
	return NewInt64(0)
}

func (self s_Base) AsUInt64() *S_UInt64 {
	panic(fmt.Sprintf("json value type is a %s but not an uint64.", self.Name()))
	return NewUInt64(0)
}

func (self s_Base) AsFloat64() *S_Float64 {
	panic(fmt.Sprintf("json value type is a %s but not a float64.", self.Name()))
	return NewFloat64(0.0)
}

func (self s_Base) AsBool() *S_Bool {
	panic(fmt.Sprintf("json value type is a %s but not a bool.", self.Name()))
	return NewBool(false)
}
