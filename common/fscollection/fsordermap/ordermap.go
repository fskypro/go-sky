/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: implement an order map
@author: fanky
@version: 1.0
@date: 2024-12-28
**/

// 非线程安全

package fsordermap

import "slices"

type F_IterKeys[K any] func(func(int, K) bool)
type F_IterValues[V any] func(func(int, V) bool)
type F_IterItems[K, V any] func(func(K, V) bool)

// -----------------------------------------------------------------------------
// OrderMap
// -----------------------------------------------------------------------------
type S_OrderMap[K comparable, V any] struct {
	mdata map[K]V
	lkeys []K
}

func NewOrderMap[K comparable, V any]() *S_OrderMap[K, V] {
	return &S_OrderMap[K, V]{
		mdata: map[K]V{},
		lkeys: make([]K, 0),
	}
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func (this *S_OrderMap[K, V]) delKey(k K) {
	keys := []K{}
	for _, kk := range this.lkeys {
		if kk != k {
			keys = append(keys, kk)
		}
	}
	this.lkeys = keys
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_OrderMap[K, V]) Keys() []K {
	return slices.Clone(this.lkeys)
}

func (this *S_OrderMap[K, V]) Values() []V {
	values := []V{}
	for _, v := range this.mdata {
		values = append(values, v)
	}
	return values
}

// 获取指定 key 的值
func (this *S_OrderMap[K, V]) Get(k K) V {
	return this.mdata[k]
}

// 获取指定 key 的值，并检查 key 是否存在
func (this *S_OrderMap[K, V]) CheckGet(k K) (V, bool) {
	v, ok := this.mdata[k]
	return v, ok
}

// 修改或添加元素
// 如果 key 已经存在，则不会调整顺序
func (this *S_OrderMap[K, V]) Set(k K, v V) {
	_, ok := this.mdata[k]
	this.mdata[k] = v
	if !ok {
		this.lkeys = append(this.lkeys, k)
	}
}

// 修改或追加元素
// 不管 key 是否已经存在，都会把新加的元素放最后面
func (this *S_OrderMap[K, V]) Append(k K, v V) {
	_, ok := this.mdata[k]
	this.mdata[k] = v
	if ok {
		this.delKey(k)
	}
	this.lkeys = append(this.lkeys, k)
}

// 删除元素，如果元素存在，则返回 true，否则返回 false
func (this *S_OrderMap[K, V]) Del(k K) bool {
	_, ok := this.mdata[k]
	if ok {
		delete(this.mdata, k)
		this.delKey(k)
	}
	return ok
}

// 删除元素，并如果元素存在则返回被删除元素和 true，否则返回被删除元素类型的默认值和 false
func (this *S_OrderMap[K, V]) Remove(k K) (v V, ok bool) {
	v, ok = this.mdata[k]
	if ok {
		delete(this.mdata, k)
		this.delKey(k)
	}
	return
}

// 判断是否存在指定的 key
func (this *S_OrderMap[K, V]) HasKey(k K) bool {
	_, ok := this.mdata[k]
	return ok
}

// key 迭代器
func (this *S_OrderMap[K, V]) IterKeys() F_IterKeys[K] {
	return func(fun func(int, K) bool) {
		for i, k := range this.lkeys {
			if !fun(i, k) {
				break
			}
		}
	}
}

// value 迭代器
func (this *S_OrderMap[K, V]) IterValues() F_IterValues[V] {
	return func(yield func(int, V) bool) {
		for i, k := range this.lkeys {
			if !yield(i, this.mdata[k]) {
				break
			}
		}
	}
}

// key/value 迭代器
func (this *S_OrderMap[K, V]) IterItems() F_IterItems[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range this.lkeys {
			if !yield(k, this.mdata[k]) {
				break
			}
		}
	}
}
