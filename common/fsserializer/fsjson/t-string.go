/**
@copyright: fantasysky 2016
@brief: 字符串值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package fsjson

import (
	"bufio"
	"fmt"

	"fsky.pro/fsstr/convert"
	"fsky.pro/fstype"
)

type S_String struct {
	s_Base
	value []byte
}

func NewString(value string) *S_String {
	return &S_String{
		s_Base: createBase(TString),
		value:  []byte(value),
	}
}

func newString(bstr []byte) *S_String {
	return &S_String{
		s_Base: createBase(TString),
		value:  bstr,
	}
}

func (this *S_String) V() string {
	return convert.Bytes2String(this.value)
}

func (this *S_String) Bytes() []byte {
	return this.value
}

func (this *S_String) AsString() *S_String {
	return this
}

func (this *S_String) WriteTo(w *bufio.Writer) (int, error) {
	return w.WriteString(this.String())
}

func (this *S_String) String() string {
	return convert.Bytes2String(this.value)
}

func (this *S_String) FmtString() string {
	return fmt.Sprintf("%q", convert.Bytes2String(this.value))
}

func JStringTo[T fstype.T_AllString](jstr *S_String) T {
	str := string(jstr.value)
	return T(str)
}
