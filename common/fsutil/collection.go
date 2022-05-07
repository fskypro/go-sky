/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: collection functions
@author: fanky
@version: 1.0
@date: 2022-05-07
**/

package fsutil

// 判断 slice 中是否存在指定元素
func SliceHasItem[T comparable](items []T, value T) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

// 找出指定元素在 slice 中的索引，不存在则返回 -1
func SliceIndexOf[T comparable](items []T, value T) int {
	for index, item := range items {
		if item == value {
			return index
		}
	}
	return -1
}

// 获取两个 slice 的交集部分
func SliceIntersection[T comparable](items1 []T, items2 []T) []T {
	items := make([]T, 0)
	for _, item1 := range items1 {
		for _, item2 := range items2 {
			if item1 == item2 {
				items = append(items, item1)
			}
		}
	}
	return items
}

// 获取 items1 中存在，items2 中不存在的集合(即 items2 的补集)
func SilceDifference[T comparable](items1 []T, items2 []T) []T {
	items := make([]T, 0)
	for _, item1 := range items1 {
		for _, item2 := range items2 {
			if item1 == item2 {
				continue
			}
		}
		items = append(items, item1)
	}
	return items
}

// -------------------------------------------------------------------
// 判断 map 中是否存在指定的 key
func MapHasKey[K comparable, V any](items map[K]V, key K) bool {
	_, ok := items[key]
	return ok
}

// 判断 map 中是否存在指定的 value
// 注意：要求 map 的值类型，必须是可比较类型
func MapHasValue[K comparable, V comparable](items map[K]V, value V) bool {
	for _, v := range items {
		if v == value {
			return true
		}
	}
	return false
}

// 用 map m2 更新 map m1
func MapUpdate[K comparable, V any](m1 map[K]V, m2 map[K]V) {
	for k, v := range m2 {
		m1[k] = v
	}
}
