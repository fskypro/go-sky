/**
@copyright: fantasysky 2016
@brief: 时间格式化测试
@author: fanky
@version: 1.0
@date: 2019-05-21
**/

package fstime

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestFormatTime(*testing.T) {
	fstest.PrintTestBegin("FormatTime")
	t, _ := Str2DateTime("2019-05-20 20:34:56")
	fmt.Println(t)
	t, _ = Str2LocDateTime("2019-05-20 20:34:56")
	fmt.Println(t)

	fstest.PrintTestEnd()
}
