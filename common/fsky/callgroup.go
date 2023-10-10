/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 初始化器
@author: fanky
@version: 1.0
@date: 2022-03-19
**/

package fsky

type E_CGLevel int

// CG_L1 优先于 CG_L5 调用
const (
	CG_L1 E_CGLevel = 1
	CG_L2           = 2
	CG_L3           = 3
	CG_L4           = 4
	CG_L5           = 5
	CG_L6           = 6
)

// -------------------------------------------------------------------
// 无参数调用组
// -------------------------------------------------------------------
type CallGroup map[E_CGLevel][]func()

func NewCallGroup() CallGroup {
	return CallGroup(map[E_CGLevel][]func(){})
}

// 添加初始化接收函数
func (self CallGroup) Add(lv E_CGLevel, call func()) {
	calls := self[lv]
	if calls == nil {
		calls = []func(){}
	}
	calls = append(calls, call)
	self[lv] = calls
}

// 不考虑调用顺序
func (self CallGroup) AddUnorder(call func()) {
	self.Add(CG_L4, call)
}

// 调用初始化接收函数
// CG_L1 比 CG_L5 优先调用
func (self CallGroup) Call() {
	for l := CG_L1; l <= CG_L5; l++ {
		calls := self[l]
		if calls == nil {
			continue
		}
		for _, call := range calls {
			call()
		}
	}
}

// -------------------------------------------------------------------
// 单个参数调用组
// -------------------------------------------------------------------
type CallGroup1[A any] map[E_CGLevel][]func(A)

func NewCallGroup1[A any]() CallGroup1[A] {
	return CallGroup1[A](map[E_CGLevel][]func(A){})
}

// 添加初始化接收函数
func (self CallGroup1[A]) Add(lv E_CGLevel, call func(A)) {
	calls := self[lv]
	if calls == nil {
		calls = []func(A){}
	}
	calls = append(calls, call)
	self[lv] = calls
}

// 不考虑调用顺序
func (self CallGroup1[A]) AddUnorder(call func(A)) {
	self.Add(CG_L4, call)
}

// 调用初始化接收函数
// CG_L1 比 CG_L5 优先调用
func (self CallGroup1[A]) Call(arg A) {
	for l := CG_L1; l <= CG_L5; l++ {
		calls := self[l]
		if calls == nil {
			continue
		}
		for _, call := range calls {
			call(arg)
		}
	}
}

// -------------------------------------------------------------------
// 两个参数调用组
// -------------------------------------------------------------------
type CallGroup2[A1, A2 any] map[E_CGLevel][]func(A1, A2)

func NewCallGroup2[A1, A2 any]() CallGroup2[A1, A2] {
	return CallGroup2[A1, A2](map[E_CGLevel][]func(A1, A2){})
}

// 添加初始化接收函数
func (self CallGroup2[A1, A2]) Add(lv E_CGLevel, call func(A1, A2)) {
	calls := self[lv]
	if calls == nil {
		calls = []func(A1, A2){}
	}
	calls = append(calls, call)
	self[lv] = calls
}

// 不考虑调用顺序
func (self CallGroup2[T1, T2]) AddUnorder(call func(T1, T2)) {
	self.Add(CG_L4, call)
}

// 调用初始化接收函数
// CG_L1 比 CG_L5 优先调用
func (self CallGroup2[A1, A2]) Call(arg1 A1, arg2 A2) {
	for l := CG_L1; l <= CG_L5; l++ {
		calls := self[l]
		if calls == nil {
			continue
		}
		for _, call := range calls {
			call(arg1, arg2)
		}
	}
}

// -------------------------------------------------------------------
// 三个参数调用组
// -------------------------------------------------------------------
type CallGroup3[A1, A2, A3 any] map[E_CGLevel][]func(A1, A2, A3)

func NewCallGroup3[A1, A2, A3 any]() CallGroup3[A1, A2, A3] {
	return CallGroup3[A1, A2, A3](map[E_CGLevel][]func(A1, A2, A3){})
}

// 添加初始化接收函数
func (self CallGroup3[A1, A2, A3]) Add(lv E_CGLevel, call func(A1, A2, A3)) {
	calls := self[lv]
	if calls == nil {
		calls = []func(A1, A2, A3){}
	}
	calls = append(calls, call)
	self[lv] = calls
}

// 不考虑调用顺序
func (self CallGroup3[A1, A2, A3]) AddUnorder(call func(A1, A2, A3)) {
	self.Add(CG_L4, call)
}

// 调用初始化接收函数
// CG_L1 比 CG_L5 优先调用
func (self CallGroup3[A1, A2, A3]) Call(arg1 A1, arg2 A2, arg3 A3) {
	for l := CG_L1; l <= CG_L5; l++ {
		calls := self[l]
		if calls == nil {
			continue
		}
		for _, call := range calls {
			call(arg1, arg2, arg3)
		}
	}
}

// -------------------------------------------------------------------
// 四个参数调用组
// -------------------------------------------------------------------
type CallGroup4[A1, A2, A3, A4 any] map[E_CGLevel][]func(A1, A2, A3, A4)

func NewCallGroup4[A1, A2, A3, A4 any]() CallGroup4[A1, A2, A3, A4] {
	return CallGroup4[A1, A2, A3, A4](map[E_CGLevel][]func(A1, A2, A3, A4){})
}

// 添加初始化接收函数
func (self CallGroup4[A1, A2, A3, A4]) Add(lv E_CGLevel, call func(A1, A2, A3, A4)) {
	calls := self[lv]
	if calls == nil {
		calls = []func(A1, A2, A3, A4){}
	}
	calls = append(calls, call)
	self[lv] = calls
}

// 不考虑调用顺序
func (self CallGroup4[A1, A2, A3, A4]) AddUnorder(call func(A1, A2, A3, A4)) {
	self.Add(CG_L4, call)
}

// 调用初始化接收函数
// CG_L1 比 CG_L5 优先调用
func (self CallGroup4[A1, A2, A3, A4]) Call(arg1 A1, arg2 A2, arg3 A3, arg4 A4) {
	for l := CG_L1; l <= CG_L5; l++ {
		calls := self[l]
		if calls == nil {
			continue
		}
		for _, call := range calls {
			call(arg1, arg2, arg3, arg4)
		}
	}
}
