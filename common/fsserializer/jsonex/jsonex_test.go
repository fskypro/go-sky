package jsonex

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

const jstr = `
// 测试 json //
{
	/* 测试 json */
	"aa": {
		"xx": 0 0,
		"y\"y": 0x68f,
		"zz": [ 0x123, 234, "ab\"c"],
	},
	/*
	多行
	**********
	注释
	*/
	"bb": {
		"vv": 0X100,       /* 多行注释 */
		"uu": "strin\"g",  // 单行注释
		"ww": 0,
	}
}
`

func TestFilterParse(t *testing.T) {
	fstest.PrintTestBegin("filerParse")
	defer fstest.PrintTestEnd()

	fmt.Println(test(jstr))
}
