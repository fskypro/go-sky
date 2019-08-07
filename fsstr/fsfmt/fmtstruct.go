/**
@copyright: fantasysky 2016
@brief: 实现格式化一个结构体
@author: fanky
@version: 1.0
@date: 2019-01-08
**/

package fsfmt

import "fmt"
import "bytes"
import "strings"
import "fsky.pro/fsenv"

// SprintStruct 以初始化结构的格式，将一个结构体格式化为字符串
// 参数：
//	st: 要格式化的结构体
//	prefix: 整个输出结构体的每一行的前缀
//	ident: 缩进字符串
func SprintStruct(st interface{}, prefix, ident string) string {
	s := fmt.Sprintf("%#v", st)
	out := bytes.NewBufferString(prefix)

	layer := 0  // 嵌套层数
	qu := false // 进入双引号
	for i := 0; i < len(s); i += 1 {
		ch := s[i]

		// 字符串值
		if qu {
			out.WriteByte(ch)
			if ch == '"' {
				qu = false
			}
			continue
		}

		switch ch {
		// 进入字符串成员值
		case '"':
			out.WriteByte(ch)
			qu = true

		// 进入子域
		case '{':
			out.WriteByte(ch)
			last := i + 1
			if s[last] == '}' { // 空结构
				out.WriteByte('}')
				i += 1
			} else {
				layer += 1
				out.WriteString(fsenv.Endline + prefix + strings.Repeat(ident, layer))
			}

		// 离开子域
		case '}':
			out.WriteByte(',')
			layer -= 1
			out.WriteString(fsenv.Endline + prefix + strings.Repeat(ident, layer))
			out.WriteByte(ch)

		// 成员间隔
		case ',':
			out.WriteByte(ch)
			out.WriteString(fsenv.Endline + prefix + strings.Repeat(ident, layer))
			last := i + 1
			if last < len(s) && s[last] == ' ' {
				i += 1
			}

		// 变量与值分隔符
		case ':':
			out.WriteByte(ch)
			out.WriteByte(' ')

		default:
			out.WriteByte(ch)
		}
	}
	return out.String()
}
