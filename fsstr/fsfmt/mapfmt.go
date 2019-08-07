/**
@copyright: fantasysky 2016
@brief: 以 map 作为参数进行格式化
@author: fanky
@version: 1.0
@date: 2019-01-14
**/

package fsfmt

import "fmt"
import "bytes"
import "strings"
import "strconv"

// Smprintf 通过传入 map[string]interface{} 对指定字符串进行格式化
//	format: 要被格式化的字符串
//	margs：格式化参数
//	譬如：Smprintf("123 %[name]s 456", map[string]interface{}{"name": "xxx"})
//	返回：123 xxx 456
func Smprintf(format string, margs map[string]interface{}) string {
	buf := bytes.NewBuffer(make([]byte, 0))
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
				value, ok := margs[trimKey]
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
			if next < count && format[next] == '%' {
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

// Mprintln 通过传入一个 map 参数对指定字符串进行格式化，并将格式化后的内容在标准输出打印
func Mprintf(format string, margs map[string]interface{}) {
	fmt.Printf(Smprintf(format, margs))
}

// Mprintln 通过传入一个 map 参数对指定字符串进行格式化，并将格式化后的内容在标准输出打印一行
func Mprintln(format string, margs map[string]interface{}) {
	fmt.Println(Smprintf(format, margs))
}
