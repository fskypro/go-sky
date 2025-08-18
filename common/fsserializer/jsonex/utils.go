/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: json utils
@author: fanky
@version: 1.0
@date: 2025-03-07
**/

package jsonex

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
)

// 将一个对象序列化为 json 后，转换成 map
func ObjectJsonMap(obj any) (map[string]any, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("marchal object to json fail, %v", err)
	}
	m := map[string]any{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal json data to map fail, %v", err)
	}
	return m, nil
}

// -------------------------------------------------------------------
// 将一个对象序列化为 json 字节，并且可以指定 json Tag
// 注意：如果 对象成员的jtag 不存在，则会找名为 “json” 的 Tag
// 示例请看 utils_test.go
func MarshalJsonWithTag(obj any, jtag string) ([]byte, error) {
	var toMap func(any, bool) any
	toMap = func(v any, first bool) any {
		rt := reflect.TypeOf(v)
		rv := reflect.ValueOf(v)
		if rt == nil || !rv.IsValid() {
			return nil
		}
		for rt.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return nil
			}

			rt = rt.Elem()
			rv = rv.Elem()
		}
		if rt.Kind() != reflect.Struct {
			return v
		}

		if !first {
			// 加 first 判断的目的是，不对最外层结构体进行 MarshalJSON 成员方法判断
			// 因为假设最外层结构体定义了 MarshalJSON 方法，并且在 MarshalJSON 方法中调用了 MarshalJsonWithTag 函数，则会造成死循环
			// 另外，这里原则上还需要对 MarshalJSON 的参数和返回值进行判断的，但是出现与该方法重名却不同意义的可能性极小，所以这里不对参数和返回值进行判断
			_, ok := rt.MethodByName("MarshalJSON")
			if ok {
				return v
			}
		}

		m := map[string]any{}
	L:
		for i := range rt.NumField() {
			tf := rt.Field(i)
			vf := rv.Field(i)
			if !tf.Anonymous {
				// 非匿名成员
				tag := tf.Tag.Get(jtag)
				if tag == "" {
					tag = tf.Tag.Get("json")
				}
				if tag == "-" {
					continue
				}
				if tag == "" {
					// 如果限定的 json tag 不存在，则使用成员名称作为 key
					tag = tf.Name
					if len(tag) == 0 {
						continue
					}
					if !unicode.IsUpper(rune(tag[0])) {
						continue
					}
				}
				if vf.CanInterface() {
					m[tag] = toMap(vf.Interface(), false)
				}
				continue
			}

			// 匿名成员
			t := tf.Type
			for t.Kind() == reflect.Ptr {
				if vf.IsNil() {
					continue L
				}
				t = t.Elem()
				vf = vf.Elem()
			}
			if !vf.IsValid() || t.Kind() != reflect.Struct {
				continue
			}
			if !vf.CanInterface() {
				continue
			}

			// 父结构体
			mm := toMap(vf.Interface(), first)
			if _, ok := mm.(map[string]any); !ok {
				continue
			}
			for k, v := range mm.(map[string]any) {
				_, ok := m[k]
				if !ok {
					m[k] = v
				}
			}
		}
		return m
	}
	return json.Marshal(toMap(obj, true))
}
