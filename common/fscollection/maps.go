/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: collection functions
@author: fanky
@version: 1.0
@date: 2022-05-07
**/

package fscollection

// -------------------------------------------------------------------
// map utils
// -------------------------------------------------------------------
// 获取 map 中的值，不存在则返回默认值
func MapGet[K comparable, V any](items map[K]V, key K, def V) V {
	if v, ok := items[key]; ok {
		return v
	}
	return def
}

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

// 将 map 的 key 和 value 交换
func MapSwapKeyValue[K comparable, V comparable](m map[K]V) map[V]K {
	vk := map[V]K{}
	for k, v := range m {
		vk[v] = k
	}
	return vk
}

// 将 map 的所有 key 转换为 slice
func MapKeysToSlice[K comparable, V any](m map[K]V) []K {
	keys := []K{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// 将 map 的所有 value 转换为 slice
func MapValuesToSlice[K comparable, V any](m map[K]V) []V {
	values := []V{}
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// 将指定 map 的 key-value 互换成 value-key
func MapKVExchange[T comparable](m map[T]T) map[T]T {
	newMap := map[T]T{}
	for k, v := range m {
		newMap[v] = k
	}
	return newMap
}

// 复制一个 map
func MapCopy[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	newMap := map[K]V{}
	for k, v := range m {
		newMap[k] = v
	}
	return m
}
