/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: stack
@author: fanky
@version: 1.0
@date: 2022-10-29
**/

package fsstack

import "errors"

type F_Iter[T any] func(func(T) bool)

// -------------------------------------------------------------------
// Stack
// -------------------------------------------------------------------
type S_Stack[T comparable] struct {
	items []T
}

func New[T comparable]() *S_Stack[T] {
	return &S_Stack[T]{
		items: make([]T, 0),
	}
}

func NewWithSlice[T comparable](s []T) *S_Stack[T] {
	items := make([]T, len(s))
	for _, item := range s {
		items = append(items, item)
	}
	return &S_Stack[T]{items: items}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 元素总数
func (this *S_Stack[T]) Count() int {
	return len(this.items)
}

// 判断栈内是否有指定元素
func (this *S_Stack[T]) Has(value T) bool {
	for _, item := range this.items {
		if item == value {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------
// 获取栈底元素，如果为空栈，则会产生 panic
func (this *S_Stack[T]) MustBottom() T {
	if len(this.items) == 0 {
		panic("stack is empty")
	}
	return this.items[0]
}

// 获取栈底元素
func (this *S_Stack[T]) Bottom() (value T, err error) {
	if len(this.items) > 0 {
		value = this.MustBottom()
	} else {
		err = errors.New("stack is empty")
	}
	return
}

// 获取栈顶元素，如果为空栈，则会产生 panic
func (this *S_Stack[T]) MustTop() T {
	if len(this.items) == 0 {
		panic("stack is empty")
	}
	return this.items[len(this.items)-1]
}

// 获取栈顶元素
func (this *S_Stack[T]) Top() (value T, err error) {
	if len(this.items) > 0 {
		value = this.MustTop()
	} else {
		err = errors.New("stack is empty")
	}
	return

}

// ---------------------------------------------------------
// 入栈一个元素
func (this *S_Stack[T]) Push(item T) {
	this.items = append(this.items, item)
}

// 入栈多个元素
func (this *S_Stack[T]) Pushs(values []T) {
	this.items = append(this.items, values...)
}

// 弹出栈顶元素，如果为空栈，则会产生 panic
func (this *S_Stack[T]) MustPop() T {
	if len(this.items) == 0 {
		panic("stack is empt")
	}
	top := len(this.items) - 1
	item := this.items[top]
	this.items = this.items[:top]
	return item
}

// 弹出栈顶元素
func (this *S_Stack[T]) Pop() (value T, err error) {
	if len(this.items) > 0 {
		value = this.MustPop()
	} else {
		err = errors.New("stack is empty")
	}
	return
}

// ---------------------------------------------------------
// 从栈底开始遍历
func (this *S_Stack[T]) BFor(f func(T) bool) {
	for _, v := range this.items {
		if !f(v) {
			break
		}
	}
}

// 从栈顶开始遍历
func (this *S_Stack[T]) TFor(f func(T) bool) {
	for i := len(this.items) - 1; i >= 0; i-- {
		if !f(this.items[i]) {
			break
		}
	}
}

// 从栈底开始遍历
func (this *S_Stack[T]) BIter() F_Iter[T] {
	return func(yield func(T) bool) {
		for _, v := range this.items {
			if !yield(v) {
				break
			}
		}
	}
}

// 从栈顶开始遍历
func (this *S_Stack[T]) TIter() F_Iter[T] {
	return func(yield func(T) bool) {
		for i := len(this.items) - 1; i >= 0; i-- {
			if !yield(this.items[i]) {
				break
			}
		}
	}
}

// ---------------------------------------------------------
func (this *S_Stack[T]) ToSlice() []T {
	items := make([]T, len(this.items))
	for _, item := range this.items {
		items = append(items, item)
	}
	return items
}
