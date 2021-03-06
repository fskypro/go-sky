/**
@copyright: fantasysky 2016
@brief: Object 类型值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package freejson

import (
	"bufio"
	"fmt"
)

type S_Object struct {
	s_Base
	keys  []string
	elems map[string]I_Value
}

func NewObject() *S_Object {
	return &S_Object{
		s_Base: createBase(TObject),
		keys:   make([]string, 0),
		elems:  make(map[string]I_Value),
	}
}

// 判断指定 key 是否存在
func (this *S_Object) Has(key string) bool {
	_, ok := this.elems[key]
	return ok
}

// 添加节点
func (this *S_Object) Add(key string, elem I_Value) {
	if this.Has(key) {
		this.Del(key)
	}
	this.keys = append(this.keys, key)
	this.elems[key] = elem
}

// 删除节点
func (this *S_Object) Del(key string) I_Value {
	for index, k := range this.keys {
		if k == key {
			elem := this.elems[key]
			delete(this.elems, key)
			this.keys = append(this.keys[:index], this.keys[index+1:]...)
			return elem
		}
	}
	return nil
}

// 清除所有元素
func (this *S_Object) Clear() {
	this.keys = []string{}
	this.elems = make(map[string]I_Value)
}

// 获取指定 key 的值
func (this *S_Object) Get(key string) I_Value {
	value, _ := this.elems[key]
	return value
}

// 获取元素个数
func (this *S_Object) Count() int {
	return len(this.keys)
}

// 遍历
// 处理函数 func 返回 false，则停止遍历
func (this *S_Object) For(fun F_KeyValue) {
	for _, key := range this.keys {
		if !fun(key, this.elems[key]) {
			break
		}
	}
}

// ------------------------------------------------------------------
func (this *S_Object) AsObject() *S_Object {
	return this
}

func (this *S_Object) WriteTo(w *bufio.Writer) (int, error) {
	err := w.WriteByte('{')
	if err != nil {
		return 0, err
	}

	count := 1
	var c int
	for index, key := range this.keys {
		if index > 0 {
			err = w.WriteByte(',')
			if err != nil {
				return count, err
			}
			count += 1
		}
		w.WriteString(fmt.Sprintf("%q:", key))
		c, err = this.elems[key].WriteTo(w)
		if err != nil {
			return count, err
		}
		count += c
	}
	err = w.WriteByte('}')
	if err == nil {
		count += 1
	}
	return count, err
}

func (this *S_Object) String() string {
	str := ""
	for index, key := range this.keys {
		if index > 0 {
			str += ", "
		}
		elem := this.elems[key]
		str += fmt.Sprintf("%q: %s", key, elem)
	}
	return fmt.Sprintf("{%s}", str)
}
