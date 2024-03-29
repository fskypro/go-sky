/**
@copyright: fantasysky 2016
@brief: null 值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package fsjson

import "bufio"

type S_Null struct {
	s_Base
}

func NewNull() *S_Null {
	return &S_Null{s_Base: createBase(TNull)}
}

func (this *S_Null) V() interface{} {
	return nil
}

func (this *S_Null) AsNull() *S_Null {
	return this
}

func (this *S_Null) WriteTo(w *bufio.Writer) (int, error) {
	return w.WriteString("null")
}

func (this *S_Null) String() string {
	return "null"
}

func (this *S_Null) FmtString() string {
	return "null"
}
