/**
@copyright: exht 2015
@website: http://www.exht.com
@brief:
@author: fanky
@version: 1.0
@date: 2023-09-18
**/

package fsexcel

import (
	"math"
	"reflect"
	"unsafe"

	"fsky.pro/fsreflect"
)

var strType = reflect.TypeOf("")

// 将 excel 表格列的 26 进制转换为十进制索引
// A == 0; Z == 25; AA == 26
func ColOrderToIndex(order string) int {
	const base float64 = 26
	var result int = 0
	letters := []byte(order)
	exp := len(letters) - 1
	if exp < 0 {
		return 0
	}
	for _, c := range letters {
		n := int(c - 'A' + 1)
		result += n * int(math.Pow(base, float64(exp)))
		exp--
	}
	return result - 1
}

// 将索引转换为 excel 表格列序号
func IndexToColOrder(index int) string {
	result := ""
	for index >= 0 {
		remainder := index % 26
		c := byte(remainder) + 'A'
		result = string(c) + result
		index = index/26 - 1
	}
	return result
}

// 将指定的结构对象转换为表头字段定义
// 例如：
//
//	type Headers struct {
//	   XXX string
//	   yyy string
//	}
//
// 执行 ToExcelHeaders 后：
// var headers = ToExcelHeaders(&Headers{})
// -> headers.XXX == XXX
// -> headers.yyy == yyy
// 也就是说，ToExcelHeaders 将传入的结构体对象的成员值，转换为成员名称
// 注意：
// -> 只转换 string 的成员
// -> 忽略 tag 为 `excel:"-"` 的成员
func ToExcelHeaders[T any](h *T) *T {
	fsreflect.TrivalStructMembers(h, true, func(info *fsreflect.S_TrivalStructInfo) bool {
		tag := info.Field.Tag.Get("excel")
		if tag == "-" {
			return true
		}
		if info.Field.Type != strType {
			return true
		}

		fv := info.FieldValue
		if fv.CanSet() {
			fv.Set(reflect.ValueOf(info.Field.Name))
		} else if fv.CanAddr() {
			reflect.NewAt(strType, unsafe.Pointer(fv.UnsafeAddr())).Elem().SetString(info.Field.Name)
		}
		return true
	})
	return h
}

// 检查表中是否存在指定的全部表头
// 如果返回空字符串，则表示所有表头都存在，否则返回不存在的表头名称
func (this *s_Table[T]) Checkheaders(headers any) string {
	unexists := ""
	fsreflect.TrivalStructMembers(headers, true, func(info *fsreflect.S_TrivalStructInfo) bool {
		tag := info.Field.Tag.Get("excel")
		if tag == "-" {
			return true
		}
		if info.Field.Type != strType {
			return true
		}
		fv := info.FieldValue
		hName := ""
		if fv.CanInterface() {
			hName = info.FieldValue.Interface().(string)
		}
		if !fv.CanAddr() {
			return true
		}
		hName = reflect.NewAt(strType, unsafe.Pointer(fv.UnsafeAddr())).Elem().Interface().(string)
		if this.GetHeaderByName(hName) == nil {
			unexists = hName
			return false
		}
		return true
	})
	return unexists
}
