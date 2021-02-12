/**
@copyright: fantasysky 2016
@brief: 实现格式化一个结构体
@author: fanky
@version: 1.0
@date: 2019-01-08
**/

package fmtex

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"fsky.pro/fsenv"
)

// 获取指针类型成员
func _getPMember(st interface{}, mname string) interface{} {
	vst := reflect.ValueOf(st)
	if vst.Type().Kind() == reflect.Ptr {
		vst = vst.Elem()
	}
	field := vst.FieldByName(strings.Trim(mname, " "))
	if !field.IsValid() {
		return nil
	}
	if !field.CanInterface() {
		return nil
	}
	return field.Interface()
}

func _skipPMember(pidx *int, s string) string {
	var ret string
	count := 0
	for ; *pidx < len(s); *pidx = *pidx + 1 {
		ch := s[*pidx]
		if ch == ')' {
			if count == 0 {
				count += 1
			} else {
				ret += string(ch)
				break
			}
		}
		ret += string(ch)
	}
	return ret
}

// SprintStruct 以初始化结构的格式，将一个结构体格式化为字符串
// 参数：
//	st: 要格式化的结构体
//	prefix: 整个输出结构体的每一行的前缀
//	ident: 缩进字符串
func SprintStruct(st interface{}, prefix, ident string) string {
	s := fmt.Sprintf("%#v", st)
	out := bytes.NewBufferString(prefix)

	layer := 0       // 嵌套层数
	mname := ""      // 成员名称
	inmname := false // 是否解释成员名称
	qu := false      // 进入双引号
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
			inmname = true

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
			mname = ""
			inmname = true

			out.WriteByte(ch)
			out.WriteString(fsenv.Endline + prefix + strings.Repeat(ident, layer))
			last := i + 1
			if last < len(s) && s[last] == ' ' {
				i += 1
			}

		// 成员名称与值分隔符
		case ':':
			inmname = false
			out.WriteByte(ch)
			out.WriteByte(' ')

		// 指针型成员
		case '(':
			member := _getPMember(st, mname)
			if member == nil { // 不可导出成员
				out.WriteString(_skipPMember(&i, s))
			} else {
				elem := reflect.ValueOf(member).Elem()
				if elem.Type().Kind() == reflect.Struct {
					subPrefix := prefix + strings.Repeat(ident, layer)
					mstr := SprintStruct(member, subPrefix, ident)
					out.WriteString(mstr[len(subPrefix):])
				} else {
					out.WriteString(fmt.Sprintf("(&%s)(%#v)", elem.Type().Name(), elem.Interface()))
				}
				_skipPMember(&i, s)
			}
		default:
			if inmname {
				mname = mname + string(ch)
			}
			out.WriteByte(ch)
		}
	}
	return out.String()
}
