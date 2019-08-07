/**
@copyright: fantasysky 2016
@brief: 时间格式化测试
@author: fanky
@version: 1.0
@date: 2019-05-21
**/

package fstime

import "fmt"
import "testing"
import "fsky.pro/fstest"

func TestStr2LocDateTime(t *testing.T) {
	fstest.PrintTestBegin("Dawn")
	time, _ := Str2DateTime("2019-05-20 20:34:56")
	fmt.Println(time)
	time, _ = Str2LocDateTime("2019-05-20 20:34:56")
	fmt.Println(time)
	fstest.PrintTestEnd()
}
