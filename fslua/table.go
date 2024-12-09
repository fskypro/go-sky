/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: lua table extentions
@author: fanky
@version: 1.0
@date: 2024-12-06
**/

package fslua

import (
	"fmt"
	"reflect"
	"unsafe"

	glua "github.com/yuin/gopher-lua"
	lua "github.com/yuin/gopher-lua"
)

type s_Table struct {
}

// -----------------------------------------------------------------------------
// Unmarshal
// -----------------------------------------------------------------------------
// 数组类型
func (self s_Table) luaValueToRefArray(lv glua.LValue, ra reflect.Value) (reflect.Value, error) {
	if lv.Type() != glua.LTTable {
		return ra, fmt.Errorf("can't convert lua value <%v> to go array %v", lv, ra.Type())
	}
	tb := lv.(*glua.LTable)
	newArray := reflect.New(ra.Type()).Elem()
	for i := 0; i < ra.Len(); i++ {
		value := tb.RawGetInt(i)
		rvNew, err := self.luaValueToRefValue(value, newArray.Index(i))
		if err != nil {
			return ra, fmt.Errorf("can't assign lua value <%v> to go array %v, %v", value, ra.Type(), err)
		}
		newArray.Index(i).Set(rvNew)
	}
	return newArray, nil
}

// slice 类型
func (self s_Table) luaValueToRefSlice(lv glua.LValue, rs reflect.Value) (reflect.Value, error) {
	if lv.Type() != glua.LTTable {
		return rs, fmt.Errorf("can't convert lua value <%v> to go slice %v", lv, rs.Type())
	}
	tb := lv.(*glua.LTable)

	etype := rs.Type().Elem()
	newSlice := reflect.MakeSlice(rs.Type(), 0, 0)
	var err error
	tb.ForEach(func(key, value glua.LValue) {
		if err != nil {
			return
		}
		rvNew, e := self.luaValueToRefValue(value, reflect.New(etype).Elem())
		if e != nil {
			err = fmt.Errorf("can't assign lua value <%v> to go slice %v, %v", value, rs.Type(), e)
			return
		}
		newSlice = reflect.Append(newSlice, rvNew)
	})
	return newSlice, err
}

// map 类型
func (self s_Table) luaValueToRefMap(lv glua.LValue, rm reflect.Value) (reflect.Value, error) {
	if lv.Type() != glua.LTTable {
		return rm, fmt.Errorf("can't convert lua value <%v> to go map %v", lv, rm.Type())
	}
	tb := lv.(*glua.LTable)

	keyType := rm.Type().Key()
	valueType := rm.Type().Elem()
	newMap := reflect.MakeMap(rm.Type())
	var err error
	tb.ForEach(func(key, value glua.LValue) {
		if err != nil {
			return
		}
		refKey, e := self.luaValueToRefValue(key, reflect.New(keyType).Elem())
		if e != nil {
			err = fmt.Errorf("can't assign lua value <%v> to go map %v, map key type error: %v", value, rm.Type(), e)
			return
		}
		refValue, e := self.luaValueToRefValue(value, reflect.New(valueType).Elem())
		if e != nil {
			err = fmt.Errorf("can't assign lua value <%v> to go map %v, map value type error: %v", value, rm.Type(), e)
			return
		}
		newMap.SetMapIndex(refKey, refValue)
	})
	return newMap, err
}

// 结构体类型
func (self s_Table) luaValueToRefObj(lv glua.LValue, robj reflect.Value) (reflect.Value, error) {
	if lv.Type() != glua.LTTable {
		return robj, fmt.Errorf("can't convert lua value <%v> to go struct %v", lv, robj.Type())
	}
	tb := lv.(*glua.LTable)

	setFieldValue := func(tfield reflect.StructField, vfield, value reflect.Value) error {
		if !value.IsValid() {
			return nil
		}
		if vfield.CanSet() {
			vfield.Set(value)
			return nil
		}
		if !vfield.CanAddr() {
			return fmt.Errorf("member %q of struct %v can't be accessable", tfield.Name, robj.Type())
		}
		pv := vfield.Addr()
		if !pv.Elem().CanSet() {
			pv = reflect.NewAt(vfield.Type(), unsafe.Pointer(vfield.UnsafeAddr()))
		}

		if pv.Elem().CanSet() {
			defer func() {
				// 部分值无法设置，提示如下：
				// panic: reflect: reflect.Value.Set using value obtained using unexported field
				// 而，刚好这部分值，不需要重新设置，因为已经在其内部修改了成员值，这里直接忽略掉这部分设置
				recover()
			}()
			pv.Elem().Set(value)
		}
		return nil
	}

	for i := 0; i < robj.Type().NumField(); i++ {
		tfield := robj.Type().Field(i)
		tag := tfield.Tag.Get("lua")
		if tag == "" {
			tag = tfield.Name
		} else if tag == "-" {
			continue
		}
		vfield := robj.Field(i)

		// 基类结构体
		if tfield.Anonymous {
			rvNew, err := self.luaValueToRefValue(lv, vfield)
			if err != nil {
				return robj, fmt.Errorf("can't set lua value <%v> to go struct member %q in %v fail, %v", lv, tfield.Name, robj.Type(), err)
			}
			if err := setFieldValue(tfield, vfield, rvNew); err != nil {
				return robj, err
			}
			continue
		}

		luaValue := tb.RawGetString(tag)
		// lua table 中不存在对应的 key 值
		if luaValue == glua.LNil {
			continue
		}

		// 获取 field 值
		rvNew, err := self.luaValueToRefValue(luaValue, vfield)
		if err != nil {
			return robj, fmt.Errorf("can't set lua value <%v> to go struct member %q in %v fail, %v", luaValue, tfield.Name, robj.Type(), err)
		}

		if err := setFieldValue(tfield, vfield, rvNew); err != nil {
			return robj, err
		}
	}
	return robj, nil
}

// 通过 go 对象的 MumarshalLua 方法反序列化 lua 数据
func (self s_Table) unmarshalByMethod(lv lua.LValue, vobj reflect.Value) (bool, error) {
	const methodName = "UnmarshalLua"

	if vobj.Type().Kind() == reflect.Interface && vobj.IsNil() {
		// 不能对 nil interface 取 method
		return false, nil
	}

	method := vobj.MethodByName(methodName)
	if !method.IsValid() {
		if !vobj.CanAddr() {
			return false, nil
		}
		pobj := reflect.NewAt(vobj.Type(), unsafe.Pointer(vobj.UnsafeAddr()))
		method = pobj.MethodByName(methodName)
		if !method.IsValid() {
			return false, nil
		}
	} else if vobj.Type().Kind() == reflect.Ptr && vobj.IsNil() {
		if vobj.CanSet() {
			vobj.Set(reflect.New(vobj.Type().Elem()))
		} else if !vobj.CanAddr() {
			return false, nil
		} else {
			pobj := reflect.NewAt(vobj.Type(), unsafe.Pointer(vobj.UnsafeAddr()))
			vobj = reflect.New(vobj.Type().Elem())
			pobj.Elem().Set(vobj)
		}
		method = vobj.MethodByName(methodName)
		if !method.IsValid() {
			return false, nil
		}
	}

	// 参数以及返回值个数判断
	if method.Type().NumOut() != 1 || method.Type().NumIn() != 1 {
		return false, nil
	}
	// 参数类型判断
	if method.Type().In(0) != reflect.TypeOf((*glua.LValue)(nil)).Elem() {
		return false, nil
	}
	// 返回值类型判断
	if method.Type().Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return false, nil
	}
	// 调用自定义解释方法
	rets := method.Call([]reflect.Value{reflect.ValueOf(lv)})
	if !rets[0].CanInterface() {
		return false, nil
	}
	if rets[0].Interface() == nil {
		return true, nil
	}
	return true, rets[0].Interface().(error)
}

func (self s_Table) luaValueToRefValue(lv glua.LValue, rv reflect.Value) (reflect.Value, error) {
	if !rv.IsValid() {
		return rv, fmt.Errorf("object must be a non-nil pointer")
	}
	tobj := rv.Type()

	var rvNew reflect.Value
	ok, err := self.unmarshalByMethod(lv, rv)
	if err != nil {
		return rvNew, fmt.Errorf("unmarshal lua value <%v> by method %q of object %v fail, %v", lv, "UnmarshalLua", tobj, err)
	}
	if ok {
		return rv, nil
	}

	// 普通类型
	switch lv.Type() {
	case glua.LTBool:
		rvNew = reflect.ValueOf(bool(lv.(glua.LBool)))
	case glua.LTString:
		rvNew = reflect.ValueOf(string(lv.(glua.LString)))
	case glua.LTNumber:
		rvNew = reflect.ValueOf(float64(lv.(glua.LNumber)))
	case glua.LTNil:
		rvNew = reflect.Zero(tobj)
	}
	if rvNew.IsValid() {
		if rvNew.CanConvert(tobj) {
			return rvNew.Convert(tobj), nil
		}
		return rvNew, fmt.Errorf("can't set lua value <%v> to go type %v", lv, tobj)
	}

	switch tobj.Kind() {
	case reflect.Interface:
		if lv.Type() != glua.LTTable {
			break
		}
		// 将 table 赋值给不确定类型，则认为这个不确定类型为 Map
		newMap := reflect.MakeMap(reflect.TypeOf((map[any]any)(nil)))
		rvNew, err := self.luaValueToRefValue(lv, newMap)
		if err != nil {
			return rv, fmt.Errorf("can't set lua value <%v> to go any type value, %v", lv, err)
		}
		if rv.CanSet() {
			defer func() { recover() }()
			rv.Set(rvNew)
			return rv, nil
		}
	case reflect.Array:
		// 数组
		return self.luaValueToRefArray(lv, rv)
	case reflect.Slice:
		// Slice
		return self.luaValueToRefSlice(lv, rv)
	case reflect.Map:
		// map
		return self.luaValueToRefMap(lv, rv)
	case reflect.Struct:
		// 结构体
		return self.luaValueToRefObj(lv, rv)
	case reflect.Ptr:
		// 指针
		if lv == glua.LNil {
			break
		}
		if rv.IsNil() {
			rv = reflect.New(tobj.Elem())
		}
		rvNew, err := self.luaValueToRefValue(lv, rv.Elem())
		if err != nil {
			return rvNew, fmt.Errorf("can't set lua value <%v> to go type %v, %v", lv, tobj, err)
		}
		if rv.Elem().CanSet() {
			defer func() { recover() }()
			rv.Elem().Set(rvNew)
		}
		return rv, nil
	}
	return rvNew, nil
}

// -----------------------------------------------------------------------------
// Marshal
// -----------------------------------------------------------------------------
// 将 array/slice 转换为 lua table
func (self s_Table) refListToLuaValue(rv reflect.Value) (glua.LValue, error) {
	tb := &glua.LTable{}
	for i := 0; i < rv.Len(); i++ {
		lv, err := self.refValueToLuaValue(rv.Index(i))
		if err != nil {
			return nil, fmt.Errorf("get value of index %d in %v fail, %v", i, rv.Type(), err)
		} else if lv == nil {
			continue
		}
		tb.RawSet(glua.LNumber(i), lv)
	}
	return tb, nil
}

// 将 map 转换为 lua table
func (self s_Table) refMapToLuaValue(rv reflect.Value) (glua.LValue, error) {
	tb := &glua.LTable{}
	for _, key := range rv.MapKeys() {
		lkey, err := self.refValueToLuaValue(key)
		if err != nil {
			return nil, fmt.Errorf(`get key "%v" in map fail, %v`, key, err)
		} else if lkey == nil {
			continue
		}
		lvalue, err := self.refValueToLuaValue(rv.MapIndex(key))
		if err != nil {
			return nil, fmt.Errorf(`get value "%v" in map fail, %v`, key, err)
		} else if lvalue == nil {
			continue
		}
		tb.RawSet(lkey, lvalue)
	}
	return tb, nil

}

// 将 struct 转换为 lua table
func (self s_Table) refObjToLuaValue(rv reflect.Value) (glua.LValue, error) {
	tb := &glua.LTable{}
	for i := 0; i < rv.Type().NumField(); i++ {
		tfield := rv.Type().Field(i)
		vfield := rv.Field(i)
		tag := tfield.Tag.Get("lua")
		if tag == "" {
			tag = tfield.Name
		} else if tag == "-" {
			continue
		}

		// 基类结构
		if tfield.Anonymous {
			tbBase, err := self.refValueToLuaValue(vfield)
			if err != nil {
				return nil, fmt.Errorf("get base object of type %v in %v fail, %v", tfield.Type, rv.Type(), err)
			}
			if tbBase.Type() == glua.LTTable {
				tbBase.(*glua.LTable).ForEach(func(key, value lua.LValue) {
					tb.RawSet(key, value)
				})
			} else {
				tb.RawSetInt(tb.Len()+1, tbBase)
			}
			continue
		}

		lv, err := self.refValueToLuaValue(vfield)
		if err != nil {
			return glua.LNil, fmt.Errorf("get field value which named %q in struct %v fail, %v", tfield.Name, rv.Type(), err)
		} else if lv == nil {
			continue
		}
		tb.RawSet(glua.LString(tfield.Name), lv)
	}
	return tb, nil
}

func (self s_Table) objMethodToLuaValue(rv reflect.Value) (bool, glua.LValue) {
	const methodName = "MarshalLua"
	method := rv.MethodByName(methodName)
	if !method.IsValid() {
		if !rv.CanAddr() {
			return false, glua.LNil
		}
		pv := reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr()))
		method = pv.MethodByName(methodName)
		if !method.IsValid() {
			return false, glua.LNil
		}
	}

	// 参数以及返回值个数判断
	if method.Type().NumOut() != 1 || method.Type().NumIn() != 0 {
		return false, nil
	}
	// 返回值类型判断
	if method.Type().Out(0) != reflect.TypeOf((*glua.LValue)(nil)).Elem() {
		return false, nil
	}
	// 调用自定义解释方法
	rets := method.Call([]reflect.Value{})
	if !rets[0].CanInterface() {
		return false, nil
	}
	if rets[0].Interface() == nil {
		return false, nil
	}
	return true, rets[0].Interface().(glua.LValue)
}

func (self s_Table) refValueToLuaValue(rv reflect.Value) (glua.LValue, error) {
	if !rv.IsValid() {
		return glua.LNil, nil
	}
	tv := rv.Type()

	switch tv.Kind() {
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return glua.LNil, nil
		}
		return self.refValueToLuaValue(rv.Elem())
	case reflect.Array, reflect.Slice:
		return self.refListToLuaValue(rv)
	case reflect.Map:
		return self.refMapToLuaValue(rv)
	case reflect.Struct:
		return self.refObjToLuaValue(rv)
	}

	getValue := func(v reflect.Value) reflect.Value {
		if v.CanInterface() {
			return v
		}
		if v.CanAddr() {
			return reflect.NewAt(tv, unsafe.Pointer(v.UnsafeAddr())).Elem()
		}
		return reflect.ValueOf(fmt.Sprintf("%v", v))
	}

	if rv.CanConvert(reflect.TypeOf(float64(0))) {
		v := getValue(rv).Convert(reflect.TypeOf(float64(0)))
		return glua.LNumber(v.Interface().(float64)), nil
	}
	if rv.CanConvert(reflect.TypeOf(true)) {
		v := getValue(rv).Convert(reflect.TypeOf(true))
		return glua.LBool(v.Interface().(bool)), nil
	}
	if rv.CanConvert(reflect.TypeOf("")) {
		v := getValue(rv).Convert(reflect.TypeOf(""))
		return glua.LString(v.Interface().(string)), nil
	}
	return nil, nil
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
// 将 lua table 解释到 go 结构对象
func UnmarshalTable(tb *glua.LTable, obj any) error {
	tobj := reflect.TypeOf(obj)
	vobj := reflect.ValueOf(obj)
	if tobj == nil || tobj.Kind() != reflect.Ptr || vobj.IsNil() {
		return fmt.Errorf("object must be a non-nil pointer")
	}
	for tobj.Kind() == reflect.Ptr {
		tobj = tobj.Elem()
		if vobj.IsNil() {
			vobj = reflect.New(tobj)
		}
		vobj = vobj.Elem()
	}

	rvalue, err := s_Table{}.luaValueToRefValue(tb, vobj)
	if err != nil {
		return err
	}
	switch tobj.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		if vobj.CanSet() {
			vobj.Set(rvalue)
		}
		if !vobj.CanAddr() {
			return fmt.Errorf("object of type %v can't be accessable", tobj)
		}
		pv := vobj.Addr()
		if !pv.Elem().CanSet() {
			pv = reflect.NewAt(tobj, unsafe.Pointer(vobj.UnsafeAddr()))
		}
		pv.Elem().Set(rvalue)
		return nil
	case reflect.Struct:
		return nil
	}
	return fmt.Errorf("can't unmsrshal lua table to go type %v", tobj)
}

// 将 go 结构对象解释为 lua table
func MarshalTable(obj any) (*glua.LTable, error) {
	tobj := reflect.TypeOf(obj)
	vobj := reflect.ValueOf(obj)
	if tobj == nil || tobj.Kind() != reflect.Ptr || vobj.IsNil() {
		return &glua.LTable{}, nil
	}
	for tobj.Kind() == reflect.Ptr {
		tobj = tobj.Elem()
		if vobj.IsNil() {
			vobj = reflect.New(tobj)
		}
		vobj = vobj.Elem()
	}
	switch tobj.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		lv, err := s_Table{}.refValueToLuaValue(vobj)
		if err != nil {
			return nil, err
		}
		return lv.(*glua.LTable), nil
	}
	return nil, fmt.Errorf("object must be an object pointer or array/slice/map instances")
}
