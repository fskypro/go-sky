package jsonex

import (
	"fmt"
	"testing"

	"fsky.pro/fsstr/fmtex"
	"fsky.pro/fstest"
)

type s_Sub struct {
	Abc int `json:"abc"`
}

type s_Root struct {
	Aa int    `json:"aa"`
	Bb string `json:"bb"`
	Cc int    `json:"cc"`

	Dd *s_Sub `json:"dd"`
}

func TestLoad(t *testing.T) {
	fstest.PrintTestBegin("Load")
	js := new(s_Root)
	err := Load("./test.js", js)
	if err != nil {
		fmt.Println("load jsonex file fail: ", err.Error())
		return
	}

	fmt.Println("load jsonex file success:")
	fmt.Println(fmtex.SprintStruct(js, "\t", "  "))
	fstest.PrintTestEnd()
}
