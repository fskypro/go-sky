/**
@copyright: yxht 2015
@website: http://www.exht.com
@brief: lua extension
@author: fanky
@version: 1.0
@date: 2024-12-10
**/

package fslua

import (
	"fmt"
	"os"

	glua "github.com/yuin/gopher-lua"
)

type i_LState interface {
	Get(idx int) glua.LValue

	SetField(obj glua.LValue, key string, value glua.LValue)
	GetField(obj glua.LValue, skey string) glua.LValue

	GetGlobal(name string) glua.LValue
	SetGlobal(name string, value glua.LValue)

	Push(value glua.LValue)
	Pop(n int)
	//GetTop() int
	//RaiseError(format string, args ...any)

	NewTable() *glua.LTable
	NewUserData() *glua.LUserData

	NewFunction(fn glua.LGFunction) *glua.LFunction
	SetFuncs(tb *glua.LTable, funcs map[string]glua.LGFunction, upvalues ...glua.LValue) *glua.LTable
	CallByParam(cp glua.P, args ...glua.LValue) error

	GetTypeMetatable(typ string) glua.LValue
	NewTypeMetatable(typ string) *glua.LTable
	GetMetatable(obj glua.LValue) glua.LValue
	SetMetatable(obj glua.LValue, mt glua.LValue)
}

// -------------------------------------------------------------------
// LState
// -------------------------------------------------------------------
type S_LState struct {
	*glua.LState
}

func (this *S_LState) Close() (err error) {
	defer func() {
		// lua 库有 bug：如果脚本中，某个全局变量引用了一个 table 中不存在的下标，则调用 Close 函数时，会 panic
		if e := recover(); e != nil {
			err = fmt.Errorf("close lua state panic: %v", e)
		}
	}()
	this.LState.Close()
	return
}

func NewState() *S_LState {
	L := &S_LState{
		LState: glua.NewState(),
	}
	return L
}

func (this *S_LState) SetPaths(path ...string) {
	SetPaths(this.LState, path...)
}

func (this *S_LState) AddPath(path ...string) {
	AddPaths(this.LState, path...)
}

// -------------------------------------------------------------------
func (this *S_LState) DoFile(file string) error {
	stat, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("open file %q fail, %v", file, err)
	}
	if stat.IsDir() {
		return fmt.Errorf("file %q is a directory", file)
	}
	if stat.Mode()&0400 == 0 {
		return fmt.Errorf("file %q can't be accessable", file)
	}
	return this.LState.DoFile(file)
}

// 将 lua 全局变量解释到 go 对象
func (this *S_LState) ParseGloal(name string, v any) error {
	lv := this.GetGlobal(name)
	err := ParseLuaValue(lv, v)
	if err != nil {
		return err
	}
	return nil
}

// 设置全局函数
func (this *S_LState) SetGlobalFunc(name string, fun any) error {
	return SetGlobalFunc(this, name, fun)
}

// -------------------------------------------------------------------
func (this *S_LState) CallGlobalFuncN(name string, nret int, args ...any) error {
	lfunc := this.GetGlobal(name)
	err := s_Function{}.callFunc(this, lfunc, nret, args...)
	if err != nil {
		return fmt.Errorf("call object function %q fail, %v", name, err)
	}
	return nil
}

// 调用没有返回值的全局函数
func (this *S_LState) CallGlobalFunc(name string, args ...any) error {
	return this.CallGlobalFuncN(name, 0, args...)
}

// 调用有一个返回值的全局函数
func (this *S_LState) CallGlobalFunc1(name string, args ...any) (ret glua.LValue, err error) {
	err = this.CallGlobalFuncN(name, 1, args...)
	if err != nil {
		return ret, err
	}
	ret = this.Get(-1)
	this.Pop(1)
	return
}

// 调用有两个返回值的全局函数
func (this *S_LState) CallGlobalFunc2(name string, args ...any) (ret1, ret2 glua.LValue, err error) {
	err = this.CallGlobalFuncN(name, 2, args...)
	if err != nil {
		return ret1, ret2, err
	}
	ret1 = this.Get(-2)
	ret2 = this.Get(-1)
	this.Pop(2)
	return
}

// 调用有三个返回值的全局函数
func (this *S_LState) CallGlobalFunc3(name string, args ...any) (ret1, ret2, ret3 glua.LValue, err error) {
	err = this.CallGlobalFuncN(name, 3, args...)
	if err != nil {
		return ret1, ret2, ret3, err
	}

	ret1 = this.Get(-3)
	ret2 = this.Get(-2)
	ret3 = this.Get(-1)
	this.Pop(3)
	return
}

// -------------------------------------------------------------------
// 调用 lua 脚本中的对象的成员函数
func (this *S_LState) CallMethodN(obj glua.LValue, name string, nret int, args ...any) error {
	lfunc := this.GetField(obj, name)
	if lfunc.Type() == glua.LTFunction {
		args = append([]any{obj}, args...)
	}
	err := s_Function{}.callFunc(this, lfunc, nret, args...)
	if err != nil {
		return fmt.Errorf("call object function %q fail, %v", name, err)
	}
	return nil
}

// 调用没有返回值的全局函数
func (this *S_LState) CallMethod(obj glua.LValue, name string, args ...any) error {
	return this.CallMethodN(obj, name, 0, args...)
}

// 调用有一个返回值的全局函数
func (this *S_LState) CallMethod1(obj glua.LValue, name string, args ...any) (ret glua.LValue, err error) {
	err = this.CallMethodN(obj, name, 1, args...)
	if err != nil {
		return ret, err
	}
	ret = this.Get(-1)
	this.Pop(1)
	return
}

// 调用有两个返回值的全局函数
func (this *S_LState) CallMethod2(obj glua.LValue, name string, args ...any) (ret1, ret2 glua.LValue, err error) {
	err = this.CallMethodN(obj, name, 2, args...)
	if err != nil {
		return
	}

	ret1 = this.Get(-2)
	ret2 = this.Get(-1)
	this.Pop(2)
	return
}

// 调用有三个返回值的全局函数
func (this *S_LState) CallMethod3(obj glua.LValue, name string, args ...any) (ret1, ret2, ret3 glua.LValue, err error) {
	err = this.CallMethodN(obj, name, 3, args...)
	if err != nil {
		return ret1, ret2, ret3, err
	}
	ret1 = this.Get(-3)
	ret2 = this.Get(-2)
	ret3 = this.Get(-1)
	this.Pop(3)
	return
}
