/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: debug utils
@author: fanky
@version: 1.0
@date: 2021-09-04
**/

package fsdebug

import (
	"bytes"
	"runtime/debug"
	"strings"

	"fsky.pro/fsdef"
)

// 打印调用栈
// start 表示从第几层开始
// prefix 表示每一行内容前插入的前缀字符串
func CallStack(start int, prefix string) string {
	if start < 0 {
		start = 0
	}
	start += 2
	lines := bytes.Split(debug.Stack(), []byte(fsdef.Endline))
	if len(lines) == 0 {
		return ""
	}
	train := strings.Builder{}
	train.WriteString(prefix)
	train.Write(lines[0])
	train.WriteString(fsdef.Endline)
	index := 1 + start*2
	if index >= len(lines) {
		return train.String()
	}
	lines = lines[index:]
	count := len(lines)
	if count == 0 {
		return train.String()
	}
	for _, line := range lines[:count-1] {
		train.WriteString(prefix)
		train.Write(line)
		train.WriteString(fsdef.Endline)
	}
	train.WriteString(prefix)
	train.Write(lines[count-1])
	return train.String()
}
