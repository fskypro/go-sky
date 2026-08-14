/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: utils
@author: fanky
@version: 1.0
@date: 2025-05-30
**/

package fshttp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

// 根据参数中的 get 函数，反序列化请求参数中的每个字段
func UnmarshalRequestArgs(obj any, get func(string) (string, bool)) error {
	tobj := reflect.TypeOf(obj)
	if tobj == nil || tobj.Kind() != reflect.Pointer {
		return fmt.Errorf("output obj argument must be a un-nil struct object pointer")
	}
	vobj := reflect.ValueOf(obj)
	if vobj.IsNil() {
		return fmt.Errorf("output obj argument must be a un-nil struct object pointer")
	}
	tobj = tobj.Elem()
	vobj = vobj.Elem()
L:
	for i := 0; i < tobj.NumField(); i++ {
		sfield := tobj.Field(i)
		if sfield.Anonymous {
			continue
		}
		if strings.HasSuffix(sfield.Name, "__") {
			// __ 结尾的为默认值成员
			continue
		}

		vfield := vobj.Field(i)
		tfield := vfield.Type()

		tag := sfield.Tag.Get("urlkey")
		if tag == "-" {
			continue
		}
		if tag == "" {
			tag = sfield.Tag.Get("json")
			if tag == "-" {
				continue
			}
			if tag == "" {
				tag = sfield.Name
			}
		}
		var value any
		svalue, ok := get(tag)
		if !ok {
			// 没有传入值
			defValue := vobj.FieldByName(sfield.Name + "__") // 原成员名称后面加 “__” 表示默认值
			if !defValue.IsValid() {                         // 默认值也不设置，则表示该参数是必传参数
				return newNoReqArgError(tag)
			}
			if defValue.Type() != tfield {
				return fmt.Errorf("the type of default value for member is not match, name=%q, default-name=%q", sfield.Name, sfield.Name+"__")
			}
			reflect.NewAt(tfield, unsafe.Add(unsafe.Pointer(vobj.UnsafeAddr()), sfield.Offset)).Elem().Set(defValue)
			continue L
		} else if svalue == "" {
			continue L
		}
		switch ft := tfield.Kind(); ft {
		case reflect.String:
			value = svalue
		case reflect.Int:
			if v, err := strconv.Atoi(svalue); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = v
			}
		case reflect.Int8:
			if v, err := strconv.Atoi(svalue); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = int8(v)
			}
		case reflect.Int16:
			if v, err := strconv.Atoi(svalue); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = int16(v)
			}
		case reflect.Int32:
			if v, err := strconv.ParseInt(svalue, 10, 32); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = int32(v)
			}
		case reflect.Int64:
			if v, err := strconv.ParseInt(svalue, 10, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = int64(v)
			}
		case reflect.Uint:
			if v, err := strconv.ParseUint(svalue, 10, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint(v)
			}
		case reflect.Uint8:
			if v, err := strconv.ParseUint(svalue, 10, 8); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint8(v)
			}
		case reflect.Uint16:
			if v, err := strconv.ParseUint(svalue, 10, 16); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint16(v)
			}
		case reflect.Uint32:
			if v, err := strconv.ParseUint(svalue, 10, 32); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint32(v)
			}
		case reflect.Uint64:
			if v, err := strconv.ParseUint(svalue, 10, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint64(v)
			}
		case reflect.Bool:
			svalue = strings.ToLower(svalue)
			if svalue == "1" || svalue == "true" {
				value = true
			} else if svalue == "0" || svalue == "false" {
				value = false
			} else {
				return newReqArgTypeError(tag, tfield)
			}
		case reflect.Float32:
			if v, err := strconv.ParseFloat(svalue, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = float32(v)
			}
		case reflect.Float64:
			if v, err := strconv.ParseFloat(svalue, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = v
			}
		case _timeKind:
			v, err := time.ParseInLocation(time.DateTime, svalue, time.Local)
			if err != nil {
				v, err = time.ParseInLocation(time.DateOnly, svalue, time.Local)
			}
			if err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = v
			}
		default:
			return fmt.Errorf("unsupport value type %v", tfield)
		}
		vvalue := reflect.ValueOf(value).Convert(tfield)
		reflect.NewAt(tfield, unsafe.Add(unsafe.Pointer(vobj.UnsafeAddr()), sfield.Offset)).Elem().Set(vvalue)
	}
	return nil
}

// 解释 POST json 参数到结构体对象，如果传入参数中不存在结构体中的成员则返回错误，定义了默认值的例外
// 如：
//
//	args := &struct {
//	   Value1 string `json:"value1"`
//	   Value2 string `json:"value2"`
//	   Value1__ strring
//	}{
//	   Value1__: "default value",
//	}
//
// 用以下 json 字符串解码时，会提示缺少参数 value1：{"value2": "xxxx"}
// 用以下 json 字符串解码时，不会有问题，并且 args.Value2 == "default value"：{"value": "xxxx"}
func UnmarshalPostJsonBody(data []byte, obj any) error {
	// 对象类型判断
	var tobj = reflect.TypeOf(obj)
	if tobj == nil {
		return fmt.Errorf("output object mustn't be nil value")
	}
	if tobj.Kind() != reflect.Pointer {
		return fmt.Errorf("output object must be a pointer")
	}
	var vobj = reflect.ValueOf(obj)

	// 定位到终极类型
	for {
		if tobj.Kind() == reflect.Pointer {
			if vobj.IsNil() {
				return fmt.Errorf("output object's pointer mustn't be nil")
			}
			tobj = tobj.Elem()
			vobj = vobj.Elem()
		} else {
			break
		}
	}

	// 如果是纯 map 或 slice，则直接反序列
	if tobj.Kind() == reflect.Map || tobj.Kind() == reflect.Slice {
		err := json.Unmarshal(data, obj)
		if err != nil {
			return fmt.Errorf("unmarshal json data to %v fail", tobj)
		}
		return nil
	}

	// 查找出传入 json 的所有顶层 key
	jmap := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &jmap); err != nil {
		return fmt.Errorf("parse json data fail, %v", err)
	}
	if err := json.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("unmarshal json to object fail, %v", err)
	}

	// 遍历结构体所有成员与成员值
	type temp struct {
		rv     reflect.Value
		field  reflect.StructField
		vfield reflect.Value
	}
	var trivalStruct func(reflect.Value, map[string]*temp)
	trivalStruct = func(rv reflect.Value, members map[string]*temp) {
		if !rv.IsValid() {
			return
		}
		rt := rv.Type()
		if rt.Kind() != reflect.Struct {
			return
		}
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			if field.Tag.Get("json") == "-" && !strings.HasSuffix(field.Name, "__") {
				continue
			}
			vfield := rv.Field(i)
			if !field.Anonymous {
				tag := field.Name
				if unicode.IsLower(rune(tag[0])) {
					continue
				}
				if strings.HasSuffix(tag, "__") {
					members[tag] = &temp{rv, field, vfield}
					continue
				} else if v, ok := field.Tag.Lookup("json"); ok {
					tag = v
				}
				members[tag] = &temp{rv, field, vfield}
				continue
			}
			// 匿名结构体
			tfield := field.Type
			for tfield.Kind() == reflect.Pointer {
				if !vfield.IsValid() {
					break
				}
				if vfield.IsNil() {
					break
				}
				tfield = tfield.Elem()
				vfield = vfield.Elem()
			}
			// 继承结构体，继续往上层遍历
			if tfield.Kind() == reflect.Struct {
				trivalStruct(vfield, members)
			}
		}
	}
	members := map[string]*temp{}
	trivalStruct(vobj, members)

	// 检查传入 json 数据中，是否有缺少的成员，
	// 如果有，则查找同名并以 __ 结尾的成员值，作为其默认值
	// 如果找不到默认值，则提示传入 json 中缺少字段
	for k, m := range members {
		if strings.HasSuffix(m.field.Name, "__") {
			// 排除默认值成员
			continue
		}
		if _, ok := jmap[k]; ok {
			// json 提供了该成员的值
			continue
		}
		// json 中没有传入值，则查找默认值
		var def = members[m.field.Name+"__"]
		if def == nil || !def.vfield.IsValid() {
			// 传入 json 中没有对应的字段，并且缺少默认值
			return newNoReqArgError(k)
		}
		if def.field.Type != m.field.Type {
			// 默认值与 json 映射字段值类型不一致
			return fmt.Errorf("the type of default member %[1]q is not match the type of member %[2]q, "+
				"typeof(%[1]s)=%[3]v, typeof(%[2]s)=%[4]v",
				def.field.Name, m.field.Name, def.field.Type, m.field.Type)
		}
		// 将默认值拷贝给对应字段成员
		reflect.NewAt(m.field.Type, unsafe.Add(unsafe.Pointer(m.rv.UnsafeAddr()), m.field.Offset)).Elem().Set(def.vfield)
	}
	return nil
}
