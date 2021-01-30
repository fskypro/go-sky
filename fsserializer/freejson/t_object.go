/**
@copyright: fantasysky 2016
@brief: Object 类型值
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package freejson

import "fmt"

type S_Object struct {
	s_Base
	keys  []string
	elems map[string]I_Value
}

func NewObject() *S_Object {
	return &S_Object{
		keys:  make([]string, 0),
		elems: make(map[string]I_Value),
	}
}

func (*S_Object) Type() JType {
	return TObject
}

func (this *S_Object) Name() string {
	return typeNames[this.Type()]
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

// 获取指定 key 的值
func (this *S_Object) Get(key string) I_Value {
	value, _ := this.elems[key]
	return value
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
