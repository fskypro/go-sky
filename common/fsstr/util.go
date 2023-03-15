/**
@copyright: fantasysky 2016
@brief: 实现一些字符串相关功能函数
@author: fanky
@version: 1.0
@date: 2019-12-25
**/

package fsstr

import (
	"fmt"
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
func JoinFunc[T any](s []T, sep string, f func(e T) string) string {
	var sb strings.Builder
	if len(s) > 0 {
		sb.WriteString(f(s[0]))
	}
	for i := 1; i < len(s); i++ {
		sb.WriteString(sep)
		sb.WriteString(f(s[i]))
	}
	return sb.String()
}

// 将 s 中的元素以 v% 的形式与 sep 逐个拼接
func JoinAny[T any](s []T, sep string) string {
	if len(s) == 0 {
		return ""
	}
	str := fmt.Sprintf("%v", s[0])
	for i := 1; i < len(s); i++ {
		str += fmt.Sprintf(",%v", s[i])
	}
	return str
}

// 将以固定字符分割元素的字符串，分割成指定的 slice
func SplitFunc[T any](str string, sep string, f func(string) T) []T {
	if str == "" {
		return []T{}
	}
	elems := []T{}
	for _, s := range strings.Split(str, sep) {
		elems = append(elems, f(s))
	}
	return elems
}
