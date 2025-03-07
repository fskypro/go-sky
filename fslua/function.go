/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: lua function
@author: fanky
@version: 1.0
@date: 2024-12-09
**/

package fslua

import (
	"fmt"
	"reflect"

	"fsky.pro/fsky"
	glua "github.com/yuin/gopher-lua"
)

// -----------------------------------------------------------------------------
// FuncWrap
// -----------------------------------------------------------------------------
type S_FuncWrap struct {
	luaName  string
	refFun   reflect.Value
	isMethod bool
}

func (this *S_FuncWrap) paramFromLua(order int, lv glua.LValue, t reflect.Type) (rv reflect.Value, err error) {
	if lv.Type() == glua.LTNumber {
		v := float64(lv.(glua.LNumber))
		rv = reflect.ValueOf(v)
	} else if lv.Type() == glua.LTBool {
		v := bool(lv.(glua.LBool))
		rv = reflect.ValueOf(v)
	} else if lv.Type() == glua.LTString {
		v := string(lv.(glua.LString))
		rv = reflect.ValueOf(v)
	}
	if rv.IsValid() {
		if rv.CanConvert(t) {
			return rv.Convert(t), nil
		}
		return rv, fmt.Errorf("can't convert number value of argument %d to %v", order, t)
	}
	if lv.Type() != glua.LTTable {
		return rv, fmt.Errorf("argument %v must be a table of %v", order, t)
	}
	v := reflect.New(t).Interface()
	if err := UnmarshalTable(lv.(*glua.LTable), v); err != nil {
		return rv, fmt.Errorf("argument %d(type=%v) can't convert to go object type %v", order, lv.Type(), t)
	}
	return reflect.ValueOf(v), nil
}

// 被 lua 脚本调用
func (this *S_FuncWrap) luaFunc(L *glua.LState) int {
	funcType := this.refFun.Type()
	lcount := L.GetTop()
	var fixArgs = []reflect.Value{}
	fixArgGount := funcType.NumIn()
	if funcType.IsVariadic() {
		// 有不定参数
		fixArgGount -= 1
		if this.isMethod {
			if lcount-1 < fixArgGount {
				L.RaiseError(fmt.Sprintf("call function %q fail, at least %d arguments must be given but not %d", this.luaName, fixArgGount, lcount-1))
				return funcType.NumOut()
			}
		} else if lcount < fixArgGount {
			L.RaiseError(fmt.Sprintf("call function %q fail, at least %d arguments must be given but not %d", this.luaName, fixArgGount, lcount))
			return funcType.NumOut()
		}
	} else {
		// 没有不定参数
		if this.isMethod {
			// 如果是对象的方法，则第一个参数为对象本身 lua.UserData
			if lcount-1 != fixArgGount {
				L.RaiseError(fmt.Sprintf("call function %q fail, %d arguments must be given but not %d", this.luaName, fixArgGount, lcount-1))
				return funcType.NumOut()
			}
		} else if lcount != fixArgGount {
			L.RaiseError(fmt.Sprintf("call function %q fail, %d arguments must be given but not %d", this.luaName, fixArgGount, lcount))
			return funcType.NumOut()
		}
	}

	// 固定参数
	var i = 1
	for ; i <= fixArgGount; i++ {
		inType := this.refFun.Type().In(i - 1)
		rparm, err := this.paramFromLua(i, L.Get(fsky.IfElse(this.isMethod, i+1, i)), inType)
		if err != nil {
			L.RaiseError(fmt.Sprintf("call function %q fail, arguments error, %v", this.luaName, err))
			return funcType.NumOut()
		}
		fixArgs = append(fixArgs, rparm)
	}

	// 不定参数
	var freeArgs = []any{}
	if funcType.IsVariadic() {
		for ; i <= lcount; i++ {
			arg, err := LuaValueToGoValue[any](L.Get(fsky.IfElse(this.isMethod, i+1, i)))
			if err != nil {
				L.RaiseError(fmt.Sprintf("call function %q fail, convert free argument fail, %v", this.luaName, err))
				return funcType.NumOut()
			}
			freeArgs = append(freeArgs, arg)
		}
	}

	defer func() {
		if e := recover(); e != nil {
			L.RaiseError("call go function %q fail, %v", this.luaName, e)
		}
	}()
	var rets []reflect.Value
	if funcType.IsVariadic() {
		rets = this.refFun.CallSlice(append(fixArgs, reflect.ValueOf(freeArgs)))
	} else {
		rets = this.refFun.Call(fixArgs)
	}

	for index, ret := range rets {
		if !ret.IsValid() {
			L.Push(glua.LNil)
			continue
		}
		if ret.Kind() == reflect.Ptr && ret.IsNil() {
			L.Push(glua.LNil)
			continue
		}

		if !ret.CanInterface() {
			if !ret.CanAddr() {
				panic(fmt.Sprintf("the %dth return value is unaccessable", index+1))
			}
			ret = reflect.NewAt(ret.Type(), ret.UnsafePointer()).Elem()
		}

		// 本身就是 lua.LValue
		if lv, ok := ret.Interface().(glua.LValue); ok {
			L.Push(lv)
			continue
		}

		// go 类型
		if ret.CanConvert(rtNumnber) {
			v := ret.Convert(rtNumnber)
			L.Push(glua.LNumber(v.Interface().(float64)))
		} else if ret.CanConvert(rtBool) {
			v := ret.Convert(rtBool)
			L.Push(glua.LBool(v.Interface().(bool)))
		} else if ret.CanConvert(rtString) {
			v := ret.Convert(rtString)
			L.Push(glua.LString(v.Interface().(string)))
		} else {
			tb, err := MarshalTable(ret.Interface())
			if err != nil {
				panic(fmt.Sprintf("marshal the %dth return value object of type %v to lua table fail, %v", index+1, ret.Type(), err))
			}
			L.Push(tb)
		}
	}
	return len(rets)
}

// -----------------------------------------------------------------------------
// private
// -----------------------------------------------------------------------------
type s_Function struct{}

// 调用 lua function
func (s_Function) callFunc(L i_LState, lfunc glua.LValue, nret int, args ...any) error {
	if lfunc.Type() != glua.LTFunction {
		m := L.GetMetatable(lfunc)
		if m == nil || m == glua.LNil {
			return fmt.Errorf(`the function or callable object is not exists`)
		}
		call := L.GetField(m, "__call") // 可调用对象
		if call == nil || call.Type() != glua.LTFunction {
			return fmt.Errorf("lua obj is not a function")
		}
	}
	luaArgs := []glua.LValue{}

F:
	for index, arg := range args {
		// 本身就是 lua Value，则直接放入参数列表
		if larg, ok := arg.(glua.LValue); ok {
			luaArgs = append(luaArgs, larg)
			continue
		}

		// 否则封装为 lua Value
		targ := reflect.TypeOf(arg)
		if targ == nil {
			luaArgs = append(luaArgs, glua.LNil)
			continue
		}

		varg := reflect.ValueOf(arg)
		for targ.Kind() == reflect.Ptr {
			if varg.IsNil() {
				luaArgs = append(luaArgs, glua.LNil)
				continue F
			}
			targ = targ.Elem()
			varg = varg.Elem()
		}

		if !varg.CanInterface() {
			if !varg.CanAddr() {
				return fmt.Errorf("the %dth argument of type %v is unaccessable", index+1, varg.Type())
			}
			varg = reflect.NewAt(targ, varg.UnsafePointer()).Elem()
		}

		if varg.CanConvert(rtNumnber) {
			varg = varg.Convert(rtNumnber)
			luaArgs = append(luaArgs, glua.LNumber(varg.Interface().(float64)))
			continue
		}
		if varg.CanConvert(rtBool) {
			varg = varg.Convert(rtBool)
			luaArgs = append(luaArgs, glua.LBool(varg.Interface().(bool)))
			continue
		}
		if varg.CanConvert(rtString) {
			varg = varg.Convert(rtString)
			luaArgs = append(luaArgs, glua.LString(varg.Interface().(string)))
			continue
		}
		switch targ.Kind() {
		case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map:
			tb, err := s_Table{}.refMarshalTable(varg)
			if err != nil {
				return fmt.Errorf("can't convert argument %d of type %q to lua argument, %v", index+1, targ, err)
			}
			luaArgs = append(luaArgs, tb)
			continue
		}

		ud := L.NewUserData()
		ud.Value = arg
		L.SetMetatable(ud, L.GetTypeMetatable(targ.Name()))
		luaArgs = append(luaArgs, ud)
	}
	return L.CallByParam(glua.P{
		Fn:      lfunc,
		NRet:    nret,
		Protect: true,
	}, luaArgs...)

}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
// 给 lua 脚本设置全局函数
func SetGlobalFunc(L i_LState, name string, fun any) error {
	vf := reflect.ValueOf(fun)
	if !vf.IsValid() {
		return fmt.Errorf("argument fun must be a callable object")
	}
	if vf.Type().Kind() != reflect.Func {
		return fmt.Errorf("argument fun must be a callable object")
	}
	wrap := &S_FuncWrap{
		luaName: name,
		refFun:  vf,
	}
	L.SetGlobal(name, L.NewFunction(wrap.luaFunc))
	return nil
}

// -------------------------------------------------------------------
// 调用 lua 脚本中的全局函数
func CallGlobalFuncN(L *glua.LState, name string, nret int, args ...any) error {
	lfunc := L.GetGlobal(name)
	err := s_Function{}.callFunc(L, lfunc, nret, args...)
	if err != nil {
		return fmt.Errorf("call object function %q fail, %v", name, err)
	}
	return nil
}

// 调用没有返回值的全局函数
func CallGlobalFunc(L *glua.LState, name string, args ...any) error {
	return CallGlobalFuncN(L, name, 0, args...)
}

// 调用有一个返回值的全局函数
func CallGlobalFunc1[R any](L *glua.LState, name string, args ...any) (R, error) {
	var ret R
	err := CallGlobalFuncN(L, name, 1, args...)
	if err != nil {
		return ret, err
	}
	lr := L.Get(-1)
	L.Pop(1)
	v, err := LuaValueToGoValue[R](lr)
	if err != nil {
		return ret, fmt.Errorf("convert return value to go type fail, %v", err)
	}
	return v, nil
}

// 调用有两个返回值的全局函数
func CallGlobalFunc2[R1, R2 any](L *glua.LState, name string, args ...any) (R1, R2, error) {
	var ret1 R1
	var ret2 R2
	err := CallGlobalFuncN(L, name, 2, args...)
	if err != nil {
		return ret1, ret2, err
	}

	lr1 := L.Get(-2)
	lr2 := L.Get(-1)
	L.Pop(2)

	ret1, err = LuaValueToGoValue[R1](lr1)
	if err != nil {
		return ret1, ret2, fmt.Errorf("convert return value 1 to go type fail, %v", err)
	}
	ret2, err = LuaValueToGoValue[R2](lr2)
	if err != nil {
		return ret1, ret2, fmt.Errorf("convert return value 2 to go type fail, %v", err)
	}
	return ret1, ret2, nil
}

// 调用有三个返回值的全局函数
func CallGlobalFunc3[R1, R2, R3 any](L *glua.LState, name string, args ...any) (R1, R2, R3, error) {
	var ret1 R1
	var ret2 R2
	var ret3 R3
	err := CallGlobalFuncN(L, name, 3, args...)
	if err != nil {
		return ret1, ret2, ret3, err
	}

	lr1 := L.Get(-3)
	lr2 := L.Get(-2)
	lr3 := L.Get(-1)
	L.Pop(3)

	ret1, err = LuaValueToGoValue[R1](lr1)
	if err != nil {
		return ret1, ret2, ret3, fmt.Errorf("convert return value 1 to go type fail, %v", err)
	}
	ret2, err = LuaValueToGoValue[R2](lr2)
	if err != nil {
		return ret1, ret2, ret3, fmt.Errorf("convert return value 2 to go type fail, %v", err)
	}
	ret3, err = LuaValueToGoValue[R3](lr3)
	if err != nil {
		return ret1, ret2, ret3, fmt.Errorf("convert return value 3 to go type fail, %v", err)
	}
	return ret1, ret2, ret3, nil
}

// -------------------------------------------------------------------
// 调用 lua 脚本中的对象的成员函数
func CallMethodN(L i_LState, obj glua.LValue, name string, nret int, args ...any) error {
	lfunc := L.GetField(obj, name)
	args = append([]any{obj}, args...)
	err := s_Function{}.callFunc(L, lfunc, nret, args...)
	if err != nil {
		return fmt.Errorf("call object function %q fail, %v", name, err)
	}
	return nil
}

// 调用没有返回值的全局函数
func CallMethod(L i_LState, obj glua.LValue, name string, args ...any) error {
	return CallMethodN(L, obj, name, 0, args...)
}

// 调用有一个返回值的全局函数
func CallMethod1[R any](L i_LState, obj glua.LValue, name string, args ...any) (R, error) {
	var ret R
	err := CallMethodN(L, obj, name, 1, args...)
	if err != nil {
		return ret, err
	}
	lr := L.Get(-1)
	L.Pop(1)
	v, err := LuaValueToGoValue[R](lr)
	if err != nil {
		return ret, fmt.Errorf("convert return value to go type fail, %v", err)
	}
	return v, nil
}

// 调用有两个返回值的全局函数
func CallMethod2[R1, R2 any](L i_LState, obj glua.LValue, name string, args ...any) (R1, R2, error) {
	var ret1 R1
	var ret2 R2
	err := CallMethodN(L, obj, name, 2, args...)
	if err != nil {
		return ret1, ret2, err
	}

	lr1 := L.Get(-2)
	lr2 := L.Get(-1)
	L.Pop(2)

	ret1, err = LuaValueToGoValue[R1](lr1)
	if err != nil {
		return ret1, ret2, fmt.Errorf("convert return value 1 to go type fail, %v", err)
	}
	ret2, err = LuaValueToGoValue[R2](lr2)
	if err != nil {
		return ret1, ret2, fmt.Errorf("convert return value 2 to go type fail, %v", err)
	}
	return ret1, ret2, nil
}

// 调用有三个返回值的全局函数
func CallMethod3[R1, R2, R3 any](L i_LState, obj glua.LValue, name string, args ...any) (R1, R2, R3, error) {
	var ret1 R1
	var ret2 R2
	var ret3 R3
	err := CallMethodN(L, obj, name, 3, args...)
	if err != nil {
		return ret1, ret2, ret3, err
	}

	lr1 := L.Get(-3)
	lr2 := L.Get(-2)
	lr3 := L.Get(-1)
	L.Pop(3)

	ret1, err = LuaValueToGoValue[R1](lr1)
	if err != nil {
		return ret1, ret2, ret3, fmt.Errorf("convert return value 1 to go type fail, %v", err)
	}
	ret2, err = LuaValueToGoValue[R2](lr2)
	if err != nil {
		return ret1, ret2, ret3, fmt.Errorf("convert return value 2 to go type fail, %v", err)
	}
	ret3, err = LuaValueToGoValue[R3](lr3)
	if err != nil {
		return ret1, ret2, ret3, fmt.Errorf("convert return value 3 to go type fail, %v", err)
	}
	return ret1, ret2, ret3, nil
}
