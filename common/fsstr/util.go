/**
@copyright: fantasysky 2016
@brief: 实现一些字符串相关功能函数
@author: fanky
@version: 1.0
@date: 2019-12-25
**/

package fsstr

import (
	"reflect"
	"strings"

	"fsky.pro/fsbytes"
	"fsky.pro/fsstr/convert"
)

// 从第 n 个字符开始查找指定字符串，如果存在，则返回第一个匹配的子字符串的位置，不存在，则返回 -1
func IndexN(str string, n int, sub string) int {
	return fsbytes.IndexN(convert.String2Bytes(str), n, convert.String2Bytes(sub))
}

// -------------------------------------------------------------------
// 去掉字符串两边空白字符
// -------------------------------------------------------------------
// 去掉字符串左边的空白字符
// 空白字符包括 ' '、'\t'、'\r'、'\n'
func TrimLeftEmpty(s string) string {
	i := -1
	var c rune
	for i, c = range s {
		switch c {
		case ' ', '\t', '\r', '\n':
			continue
		}
		i--
		break
	}
	return s[i+1:]
}

// 去掉字符串右边的空白字符
// 空白字符包括 ' '、'\t'、'\r'、'\n'
func TrimRightEmpty(s string) string {
	i := len(s) - 1
	for i >= 0 {
		switch s[i] {
		case ' ', '\t', '\r', '\n':
			i--
			continue
		}
		break
	}
	return s[0 : i+1]
}

// 去掉字符串左右两边的空白字符
// 空白字符包括 ' '、'\t'、'\r'、'\n'
func TrimEmpty(s string) string {
	s = TrimLeftEmpty(s)
	return TrimRightEmpty(s)
}

// 将 slice 中所有元素用 sep 分割
// a 必须是一个 slice 或者 array
// sep 为分隔字符串
// f 传入 a 的每一个元素， 返回元素的字符串表现形式
func JoinFunc(a interface{}, sep string, f func(e interface{}) string) string {
	t := reflect.TypeOf(a)
	if t == nil {
		return ""
	}
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return ""
	}
	v := reflect.ValueOf(a)
	if v.IsNil() {
		return ""
	}
	var sb strings.Builder
	if v.Len() > 0 {
		sb.WriteString(f(v.Index(0).Interface()))
	}
	for i := 1; i < v.Len(); i++ {
		sb.WriteString(sep)
		sb.WriteString(f(v.Index(i).Interface()))
	}
	return sb.String()
}
