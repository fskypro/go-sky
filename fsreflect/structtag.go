/**
@copyright: fantasysky 2016
@brief: struct tag
@author: fanky
@version: 1.0
@date: 2019-01-25
**/

package fsreflect

import "reflect"

// FieldTagMap 获取结构中指定 tag 类型的：File 到 tag 的映射
func FieldTagsMap(inst interface{}, tagType string) map[string]string {
	var rt reflect.Type
	rinst := reflect.TypeOf(inst)
	if rinst.Kind() == reflect.Ptr {
		rt = rinst.Elem()
	} else {
		rt = rinst
	}

	tags := make(map[string]string)
	for i := 0; i < rt.NumField(); i += 1 {
		f := rt.Field(i)
		tag := f.Tag.Get(tagType)
		if tag != "" {
			tags[f.Name] = tag
		}
	}
	return tags
}

// TagFieldMap 获取结构中指定 tag 类型的：tag 到 field 的映射
func TagFieldsMap(inst interface{}, tagType string) map[string]string {
	var rt reflect.Type
	rinst := reflect.TypeOf(inst)
	if rinst.Kind() == reflect.Ptr {
		rt = rinst.Elem()
	} else {
		rt = rinst
	}

	tags := make(map[string]string)
	for i := 0; i < rt.NumField(); i += 1 {
		f := rt.Field(i)
		tag := f.Tag.Get(tagType)
		if tag != "" {
			tags[tag] = f.Name
		}
	}
	return tags
}
