/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 通过指定关键字格式化
@author: fanky
@version: 1.0
@date: 2023-09-28
**/

package fsfmt

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func SfuncPrintf(format string, fun func(string) (any, bool)) string {
	buf := bytes.NewBuffer([]byte{})
	args := make([]interface{}, 0)
	count := len(format) // 格式化串长度
	inFmt := false       // 是否进入格式化
	inKey := false       // 是否进入 key 解释中
	var key string       // 当前 key
	site := 1            // 新的格式化位置
	for i := 0; i < count; i += 1 {
		ch := format[i]

		// 提取 key
		if inKey {
			if ch == ']' { // 获取一个键结束
				trimKey := strings.Trim(key, " ")
				value, ok := fun(trimKey)
				if !ok { // key 不存在
					buf.WriteString(key)
				} else {
					buf.WriteString(strconv.FormatInt(int64(site), 10))
					args = append(args, value)
					site += 1
				}
				inFmt = false
				inKey = false
				key = ""
				buf.WriteByte(']')
			} else {
				key = key + string(ch)
			}
			continue
		}

		// 判断是否进入格式化
		if ch == '%' {
			next := i + 1
			if next < count && format[next] == '%' { // 双 %% 号，只取一个
				inFmt = false
				i = next
				buf.WriteByte('%')
			} else {
				inFmt = true
			}
			buf.WriteByte('%')
			continue
		}

		// 进入格式化
		if inFmt {
			if ch == '[' { // 进入 key
				inKey = true
			}
		}
		buf.WriteByte(ch)
	}
	return fmt.Sprintf(buf.String(), args...)
}
