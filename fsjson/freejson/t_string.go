/**
@copyright: fantasysky 2016
@brief: 字符串值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package freejson

import "fmt"
import "fsky.pro/fsstr/convert"

type S_String struct {
	s_Base
	value []byte
}

func NewString(value string) *S_String {
	return &S_String{value: []byte(value)}
}

func newString(bstr []byte) *S_String {
	return &S_String{value: bstr}
}

func (*S_String) Type() JType {
	return TString
}

func (this *S_String) Name() string {
	return typeNames[this.Type()]
}

func (this *S_String) V() string {
	return convert.Bytes2String(this.value)
}

func (this *S_String) AsString() *S_String {
	return this
}

func (this *S_String) String() string {
	return fmt.Sprintf("%q", convert.Bytes2String(this.value))
}
