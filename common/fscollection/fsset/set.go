/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: implement set
@author: fanky
@version: 1.0
@date: 2022-05-07
**/

package fsset

type S_Set[T comparable] struct {
	items map[T]any
}

func New[T comparable]() *S_Set[T] {
	return &S_Set[T]{
		items: make(map[T]any),
	}
}

func NewWith[T comparable](s *S_Set[T]) *S_Set[T] {
	items := make(map[T]any)
	for v := range s.items {
		items[v] = nil
	}
	return &S_Set[T]{
		items: items,
	}
}

func NewWithSlice[T comparable](s []T) *S_Set[T] {
	items := make(map[T]any)
	for _, v := range s {
		items[v] = nil
	}
	return &S_Set[T]{
		items: items,
	}
}

func NewWithMapKeys[T comparable, V any](m map[T]V) *S_Set[T] {
	items := make(map[T]any)
	for v := range m {
		items[v] = nil
	}
	return &S_Set[T]{
		items: items,
	}
}

func NewWithMapValues[K comparable, T comparable](m map[K]T) *S_Set[T] {
	items := make(map[T]any)
	for _, v := range m {
		items[v] = nil
	}
	return &S_Set[T]{
		items: items,
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_Set[T]) Has(value T) bool {
	_, ok := this.items[value]
	return ok
}

func (this *S_Set[T]) Add(value T) {
	this.items[value] = nil
}

func (this *S_Set[T]) Del(value T) {
	delete(this.items, value)
}

func (this *S_Set[T]) For(f func(T) bool) {
	for v := range this.items {
		if !f(v) {
			break
		}
	}
}

func (this *S_Set[T]) AddSlice(values []T) {
	for _, value := range values {
		this.items[value] = nil 
	}
}

func (this *S_Set[T]) DelSlice(values []T) {
	for _, value := range values {
		delete(this.items, value)
	}
}

func (this *S_Set[T]) ToSlice() []T {
	items := make([]T, 0)
	for v := range this.items {
		items = append(items, v)
	}
	return items
}

// ---------------------------------------------------------
// 求交集
func (this *S_Set[T]) Intersection(s *S_Set[T]) *S_Set[T] {
	ret := New[T]()
	for v := range this.items {
		for vv := range s.items {
			if v == vv {
				ret.items[v] = nil
			}
		}
	}
	return ret
}

// 求并集
func (this *S_Set[T]) Union(s *S_Set[T]) *S_Set[T] {
	ret := New[T]()
	for v := range this.items {
		ret.items[v] = nil
	}
	for v := range s.items {
		ret.items[v] = nil
	}
	return ret
}

// 求补集，s 中存在，this 中不存在
func (this *S_Set[T]) Difference(s *S_Set[T]) *S_Set[T] {
	ret := New[T]()
	for v := range s.items {
		if _, ok := this.items[v]; !ok {
			ret.items[v] = nil
		}
	}
	return ret
}

// this 是否是 s 的超集(即 s 的元素全部在 this 中)
func (this *S_Set[T]) IsSuperSet(s *S_Set[T]) bool {
	for v := range s.items {
		if _, ok := this.items[v]; !ok {
			return false
		}
	}
	return true
}

// this 是否是 s 的子集(即 this 的元素全部在 s 中)
func (this *S_Set[T]) IsSubset(s *S_Set[T]) bool {
	for v := range this.items {
		if _, ok := s.items[v]; !ok {
			return false
		}
	}
	return true
}
