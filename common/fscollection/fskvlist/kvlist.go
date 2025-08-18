/**
@copyright: yxht 2015
@website: http://www.exht.com
@brief: definations
@author: fanky
@version: 1.0
@date: 2025-07-04
**/

package fskvlist

type S_KVItem[K, V any] struct {
	Key   K `lua:"key" json:"key"`
	Value V `lua:"value" json:"value"`
}

type T_KVList[K, V any] []*S_KVItem[K, V]

// 让每个元素的 key 和 value 交换，并返回新的 KVList
func (self T_KVList[K, V]) KVSwapList() T_KVList[V, K] {
	items := []*S_KVItem[V, K]{}
	for _, item := range self {
		items = append(items, &S_KVItem[V, K]{
			Key:   item.Value,
			Value: item.Key,
		})
	}
	return items
}

// -------------------------------------------------------------------
func ToMap[K comparable, V any](list T_KVList[K, V]) map[K]V {
	m := map[K]V{}
	for _, item := range list {
		m[item.Key] = item.Value
	}
	return m
}

func FromMap[K comparable, V any](m map[K]V) []*S_KVItem[K, V] {
	kvlist := []*S_KVItem[K, V]{}
	for k, v := range m {
		kvlist = append(kvlist, &S_KVItem[K, V]{k, v})
	}
	return kvlist
}
