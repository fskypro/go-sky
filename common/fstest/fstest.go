/**
* @brief: testlogger.go
* @copyright: 2016 fantasysky
* @author: fanky
* @version: 1.0
* @date: 2018-12-29
 */

package fstest

import (
	"fmt"
	"strings"
)

const _spcount = 69

// 打印测试起始分割线
func PrintTestBegin(name string) func() {
	splitter := "|" + strings.Repeat("-", _spcount)
	fmt.Println(splitter)
	fmt.Println("| test ", name)
	fmt.Println(splitter)
	return PrintTestEnd
}

// 打印测试结束分割线
func PrintTestEnd() {
	splitter := "|" + strings.Repeat("-", _spcount)
	fmt.Println(splitter + "\n")
}
