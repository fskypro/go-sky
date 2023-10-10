/**
@copyright: fantasysky 2016
@brief: 以 map 作为参数进行格式化
@author: fanky
@version: 1.0
@date: 2019-01-14
**/

package fsfmt

import (
	"fmt"
)

// Smprintf 通过传入 map[string]interface{} 对指定字符串进行格式化
//	format: 要被格式化的字符串
//	margs：格式化参数
//	譬如：Smprintf("123 %[k1]s 456; %.2[k2]f", map[string]interface{}{"k1": "xxx", "k2": 3.3333})
//	返回：123 xxx 456; 3.33
func Smprintf(format string, args map[string]any) string {
	return SfuncPrintf(format, func(key string) (any, bool) {
		value, ok := args[key]
		return value, ok
	})
}

// Mprintln 通过传入一个 map 参数对指定字符串进行格式化，并将格式化后的内容在标准输出打印
func Mprintf(format string, margs map[string]any) {
	fmt.Printf(Smprintf(format, margs))
}

// Mprintln 通过传入一个 map 参数对指定字符串进行格式化，并将格式化后的内容在标准输出打印一行
func Mprintln(format string, margs map[string]any) {
	fmt.Println(Smprintf(format, margs))
}
