/**
@copyright: fantasysky 2016
@brief: 列表值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package fsjson

import (
	"bufio"
	"fmt"
)

// -------------------------------------------------------------------
// ListIter
// -------------------------------------------------------------------
type S_ListIter struct {
	index int
	owner *S_List
}

func newListIter(owner *S_List) *S_ListIter {
	return &S_ListIter{-1, owner}
}

func (this *S_ListIter) Next() (I_Value, bool) {
	this.index++
	if this.index < len(this.owner.elems) {
		return this.owner.elems[this.index], true
	}
	return nil, false
}

// -------------------------------------------------------------------
// List
// -------------------------------------------------------------------
type S_List struct {
	s_Base
	elems []I_Value
}

func NewList() *S_List {
	return &S_List{
		s_Base: createBase(TList),
		elems:  make([]I_Value, 0),
	}
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

// 清除所有元素
func (this *S_List) Clear() {
	this.elems = []I_Value{}
}

// 获取指定索引处的值（支持负索引）
func (this *S_List) Get(index int) I_Value {
	if index >= 0 && index < len(this.elems) {
		return this.elems[index]
	} else if index < 0 && index >= -len(this.elems) {
		return this.elems[len(this.elems)+index]
	}
	return nil
}

// 元素个数
func (this *S_List) Count() int {
	return len(this.elems)
}

// 获取迭代器
func (this *S_List) Iter() *S_ListIter {
	return newListIter(this)
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

func (this *S_List) WriteTo(w *bufio.Writer) (int, error) {
	err := w.WriteByte('[')
	if err != nil {
		return 0, err
	}
	count := 1
	var c int
	for index, elem := range this.elems {
		if index > 0 {
			err := w.WriteByte(',')
			if err != nil {
				return count, err
			}
			count += 1
		}

		c, err = elem.WriteTo(w)
		if err != nil {
			return count, err
		}
		count += c
	}
	err = w.WriteByte(']')
	if err == nil {
		count += 1
	}
	return count, err
}

func (this *S_List) String() string {
	return this.FmtString()
}

func (this *S_List) FmtString() string {
	str := ""
	for index, elem := range this.elems {
		if index == 0 {
			str = elem.String()
		} else {
			str = str + ", " + elem.FmtString()
		}
	}
	return fmt.Sprintf("[%s]", str)
}
