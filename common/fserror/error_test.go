package fserror

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"fsky.pro/fstest"
)

type E1 struct {
	error
	S string
}

type E2 struct {
	error
}

type E3 struct {
	error
}

func test1() error {
	e1 := E1{fmt.Errorf("ADFASDFASDFAS"), "xxx"}
	e2 := E2{fmt.Errorf("ADFASDFASDFAS")}
	return errors.Join(e1, e2)
}

func test2() error {
	return E3{fmt.Errorf("xxxx: %d", 100)}
}

func TestIsErrType(t *testing.T) {
	fstest.PrintTestBegin("fserror.IsError")
	defer fstest.PrintTestEnd()
	err := test1()
	fmt.Println(1111, IsError[E1](err))
	fmt.Println(2222, IsError[E2](err))

	err = test2()
	fmt.Println(3333, IsError[E3](err))
	fstest.PrintTestEnd()
}

func TestErrorf(t *testing.T) {
	fstest.PrintTestBegin("Errorf")
	defer fstest.PrintTestEnd()

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	jerror := func(format string, args ...any) *S_JError {
		return JCFLErrorf(" ", 1, dir, format, args...)
	}

	f1 := func() error {
		return jerror("error from function1, no=%d", 1)
	}

	f2 := func() error {
		return jerror("error from function2, no=%d", 2).Join(f1())
	}

	f3 := func() error {
		return jerror("error from function3, no=%d", 3).Join(f2())
	}

	fmt.Println(f3())
}
