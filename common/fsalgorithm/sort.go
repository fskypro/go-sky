/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: sort utils
@author: fanky
@version: 1.0
@date: 2023-10-06
**/

package fsalgorithm

import "sort"

type s_SortItemsWrap[T any] struct { 
	items []T
	less func(T, T) bool
}
func (self s_SortItemsWrap[T]) Len() int { return len(self.items) }
func (self s_SortItemsWrap[T]) Swap(i, j int) { self.items[i], self.items[j] = self.items[j], self.items[i]}
func (self s_SortItemsWrap[T]) Less(i, j int) bool { return self.less(self.items[i], self.items[j]) }

// 排序任何对象列表
func SortFunc[T any](items []T, less func(T,T)bool) {
	sort.Sort(s_SortItemsWrap[T]{items, less})
}
