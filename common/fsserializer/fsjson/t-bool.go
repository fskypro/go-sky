/**
@copyright: fantasysky 2016
@brief: 布尔类型值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package fsjson

import (
	"bufio"
	"fmt"
)

type S_Bool struct {
	s_Base
	value bool
}

func NewBool(value bool) *S_Bool {
	return &S_Bool{
		s_Base: createBase(TBool),
		value:  value,
	}
}

func (this *S_Bool) V() bool {
	return this.value
}

// --------------------------------------------------------
func (this *S_Bool) AsBool() *S_Bool {
	return this
}

func (this *S_Bool) WriteTo(w *bufio.Writer) (int, error) {
	return w.WriteString(this.String())
}

func (this *S_Bool) String() string {
	return fmt.Sprintf("%v", this.value)
}

func (this *S_Bool) FmtString() string {
	return fmt.Sprintf("%v", this.value)
}
