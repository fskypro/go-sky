package fsruntime

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func test3() {
}

func test2() {
	test3()
}

func test1() {
	test2()
}

func TestGetFuncName(t *testing.T) {
	fstest.PrintTestBegin("GetFuncName")()
	defer fstest.PrintTestEnd()
	fmt.Println(11111, GetFuncName(test1))
}

func TestFindFuncWithName(t *testing.T) {
	fstest.PrintTestBegin("FindFuncWithName")()
	defer fstest.PrintTestEnd()

	fmt.Println(FindFuncWithName("fsky.pro/fsruntime.test1"))
	fmt.Println(FindFuncWithName("fsky.pro/fsruntime.test2"))
}
