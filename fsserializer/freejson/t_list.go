/**
@copyright: fantasysky 2016
@brief: 列表值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package freejson

import "fmt"

type S_List struct {
	s_Base
	elems []I_Value
}

func NewList() *S_List {
	return &S_List{
		elems: make([]I_Value, 0),
	}
}

func (*S_List) Type() JType {
	return TList
}

func (this *S_List) Name() string {
	return typeNames[this.Type()]
}

// 添加元素
func (this *S_List) Add(elem I_Value) {
	this.elems = append(this.elems, elem)
}

// 删除元素
func (this *S_List) Del(elem I_Value) {
	for index, e := range this.elems {
		if e == elem {
			this.elems = append(this.elems[:index], this.elems[index+1:]...)
			break
		}
	}
}

// 获取指定索引处的值
func (this *S_List) Get(index int) I_Value {
	if index > 0 && index < len(this.elems) {
		return this.elems[index]
	}
	return nil
}

// 遍历
// 处理函数 func 返回 false，则停止遍历
func (this *S_List) For(fun F_Elem) {
	for index, value := range this.elems {
		if !fun(index, value) {
			break
		}
	}
}

// ------------------------------------------------------------------
func (this *S_List) AsList() *S_List {
	return this
}

func (this *S_List) String() string {
	str := ""
	for index, elem := range this.elems {
		if index == 0 {
			str = elem.String()
		} else {
			str = str + ", " + elem.String()
		}
	}
	return fmt.Sprintf("[%s]", str)
}
