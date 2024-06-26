/**
@copyright: fantasysky 2016
@brief: 以 object 作为参数进行格式化
@author: fanky
@version: 1.0
@date: 2019-01-14
**/

package fsfmt

import (
	"fmt"
	"reflect"
	"unsafe"

	"fsky.pro/fsreflect"
)

// SobjPrintf
//  通过传入一个 object 参数，并以 object 的成员值作为字符串格式化参数
//	format: 要被格式化的字符串
//	args：格式化参数所属的对象
//  type Arg struct {
//		Value1 string  `fm:"k1"`
//		value2 float64 `fm:"k2"`
//  }
//	譬如：Smprintf("123 %[k1]s 456; %.2[k2]f", &Arg{"xxx", 3.33}, "fm")
//	返回：123 xxx 456; 3.33
func SobjPrintf(format string, obj any, tag string) string {
	return SfuncPrintf(format, func(key string) (any, bool) {
		ok := false
		var member reflect.Value
		fsreflect.TrivalStructMembers(obj, false, func(info *fsreflect.S_TrivalStructInfo) bool {
			if info.IsBase { return true }
			if tag == "" {
				if info.Field.Name == key {
					ok = true
					member = info.FieldValue
					return false
				}
			} else if info.Field.Tag.Get(tag) == key {
				ok = true
				member = info.FieldValue
				return false
			}
			return true
		})
		if !ok { return nil, false }

		if !member.IsValid() { return nil, false }

		if member.CanInterface() {
			return member.Interface(), true
		}
		if member.Type().Kind() == reflect.Ptr {
			if member.IsNil() { return 0, true }
			return "<can't access ptr>", true
		} else if member.CanAddr() {
			vv := reflect.NewAt(member.Type(), unsafe.Pointer(member.UnsafeAddr()))
			return vv.Elem().Interface(), true
		}
		return nil, false
	})
}

// ObjPrintln
// 通过传入一个 object 参数，并以 object 的成员值作为字符串格式化参数，将格式化后的内容在标准输出打印
func ObjPrintf(format string, args any) {
	fmt.Printf(SobjPrintf(format, args, ""))
}

// ObjPrintln
// 通过传入一个 object 参数，并以 object 的成员值作为字符串格式化参数，将格式化后的内容在标准输出打印，并将格式化后的内容在标准输出打印一行
func ObjPrintln(format string, args any) {
	fmt.Println(SobjPrintf(format, args, ""))
}
