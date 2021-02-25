package fmtex

import (
	"testing"

	"fsky.pro/fstest"
)

func TestMprintln(t *testing.T) {
	fstest.PrintTestBegin("Smprintf")

	args := make(map[string]interface{})
	args["a"] = 100
	args["b"] = "xxx"
	args["c"] = 3.455
	Mprintln("http://aaa |%-10[a]d| bbb |%10[b]s| ccc |%[c]f| ddd", args)

	fstest.PrintTestEnd()
}
