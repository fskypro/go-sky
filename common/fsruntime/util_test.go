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
	defer fstest.PrintTestBegin("GetFuncName")()
	fmt.Println(11111, GetFuncName(test1))
}

func TestFindFuncWithName(t *testing.T) {
	defer fstest.PrintTestBegin("FindFuncWithName")()

	fmt.Println(FindFuncWithName("fsky.pro/fsruntime.test1"))
	fmt.Println(FindFuncWithName("fsky.pro/fsruntime.test2"))
}
