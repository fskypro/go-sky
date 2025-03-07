/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: lua userdata
@author: fanky
@version: 1.0
@date: 2025-01-13
**/

package fslua

import (
	"fmt"
	"reflect"
	"slices"

	"fsky.pro/fsreflect"
	glua "github.com/yuin/gopher-lua"
)

type s_UserData struct{}

func (self s_UserData) toLValue(L i_LState, rv reflect.Value) glua.LValue {
	if !rv.IsValid() {
		return glua.LNil
	}
	tv := rv.Type()
	t, v := tv, rv
	for t.Kind() == reflect.Ptr {
		v = rv.Elem()
		t = tv.Elem()
	}
	if t.Kind() == reflect.Struct {
		// 注意：这里第二个参数不能用 v，必须用 rv
		// 否则将搜索不到它的成员方法
		return self.wrapUserData(L, rv, nil)
	}
	lv, err := s_Util{}.refGoValueToLuaValue(v)
	if err != nil {
		return glua.LNil
	}
	return lv
}

func (self s_UserData) wrapUserData(L i_LState, vobj reflect.Value, tbIndex *glua.LTable) glua.LValue {
	if !vobj.IsValid() {
		return glua.LNil
	}
	v := vobj
	tobj := vobj.Type()
	t := tobj
	for t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return glua.LNil
		}
		t = t.Elem()
		v = v.Elem()
	}

	// 值本身已经是 lua value
	isLValue := false
	glvalue := reflect.TypeOf((*glua.LValue)(nil)).Elem()
	if v.CanConvert(glvalue) {
		v = v.Convert(glvalue)
		isLValue = true
	}

	var obj any
	if v.CanInterface() {
		obj = v.Interface()
	} else if v.CanAddr() {
		if t.Kind() == reflect.Struct {
			obj = reflect.New(t).Elem().Interface()
		} else {
			obj = reflect.NewAt(t, v.UnsafePointer()).Elem().Interface()
		}
	} else {
		return glua.LNil
	}
	if isLValue {
		return obj.(glua.LValue)
	}

	ud := L.NewUserData()

	if tbIndex == nil {
		// 如果是基类对象不会走这里来
		ud.Value = obj
		tbMeta := L.NewTypeMetatable(t.Name())
		L.SetMetatable(ud, tbMeta)
		tbIndex = L.NewTable()
		L.SetField(tbMeta, "__index", tbIndex)

		// 导出方法
		methods := map[string]glua.LGFunction{}
		for i := 0; i < vobj.NumMethod(); i++ {
			mname := tobj.Method(i).Name
			funcWrap := &S_FuncWrap{luaName: mname, refFun: vobj.Method(i), isMethod: true}
			methods[mname] = glua.LGFunction(funcWrap.luaFunc)
		}
		L.SetFuncs(tbIndex, methods)
	}

	// 导出成员变量
	fsreflect.TrivalStructMembers(obj, true, func(info *fsreflect.S_TrivalStructInfo) bool {
		tag := info.Field.Tag.Get("lua")
		if tag == "" {
			tag = info.Field.Name
		}
		if info.IsBase {
			self.wrapUserData(L, info.FieldValue, tbIndex)
		} else {
			L.SetField(tbIndex, tag, self.toLValue(L, info.FieldValue))
		}
		return true
	})
	return ud
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 注册 go 对象
// attrs 为要导出的成员变量名称列表，attrs 等于以下值时：
//
//	nil           ：导出全部成员
//	[]string{}    ：所有成员都不导出
//	[]string{...} ：导出指定成员
//
// methods 为要导出的方法名称列表，methods 等于以下值时：
//
//	nil           ：导出全部公有方法
//	[]string{}    ：不导出任何方法
//	[]string{...} ：导出指定方法
//
// type A struct {}
func WrapUserData(L i_LState, obj any, attrs []string, methods []string) (*glua.LUserData, error) {
	tobj := reflect.TypeOf(obj)
	vobj := reflect.ValueOf(obj)
	if tobj == nil {
		return nil, fmt.Errorf("argument obj must not be a not nil value")
	}

	t := tobj
	v := vobj
	for t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("argument obj must not be a not nil value")
		}
		t = t.Elem()
		v = v.Elem()
	}

	ud := L.NewUserData()
	ud.Value = obj
	tbMeta := L.NewTypeMetatable(t.Name())
	L.SetMetatable(ud, tbMeta)
	tbIndex := L.NewTable()
	L.SetField(tbMeta, "__index", tbIndex)

	// 导出方法
	ms := map[string]glua.LGFunction{}
	if len(methods) > 0 {
		for _, member := range methods {
			m := vobj.MethodByName(member)
			if !m.IsValid() {
				return nil, fmt.Errorf("%q is not a valid method", member)
			} else {
				funcWrap := &S_FuncWrap{luaName: member, refFun: m, isMethod: true}
				ms[member] = glua.LGFunction(funcWrap.luaFunc)
			}
		}
	} else {
		for i := 0; i < vobj.NumMethod(); i++ {
			mname := tobj.Method(i).Name
			funcWrap := &S_FuncWrap{luaName: mname, refFun: vobj.Method(i), isMethod: true}
			ms[mname] = glua.LGFunction(funcWrap.luaFunc)
		}
	}
	L.SetFuncs(tbIndex, ms)

	// 导出成员变量
	if len(attrs) == 0 {
		fsreflect.TrivalStructMembers(obj, true, func(info *fsreflect.S_TrivalStructInfo) bool {
			tag := info.Field.Tag.Get("lua")
			if tag == "" {
				tag = info.Field.Name
			}
			if info.IsBase {
				s_UserData{}.wrapUserData(L, info.FieldValue, tbIndex)
			} else {
				L.SetField(tbIndex, tag, s_UserData{}.toLValue(L, info.FieldValue))
			}
			return true
		})
	} else {
		fsreflect.TrivalStructMembers(obj, true, func(info *fsreflect.S_TrivalStructInfo) bool {
			if !slices.Contains(attrs, info.Field.Name) {
				return true
			}
			attrs = slices.DeleteFunc(attrs, func(e string) bool { return e == info.Field.Name })

			tag := info.Field.Tag.Get("lua")
			if tag == "" {
				tag = info.Field.Name
			}

			if info.IsBase {
				s_UserData{}.wrapUserData(L, info.FieldValue, tbIndex)
			} else {
				L.SetField(tbIndex, tag, s_UserData{}.toLValue(L, info.FieldValue))
			}
			return len(attrs) > 0
		})
	}
	return ud, nil
}
