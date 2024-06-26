package fserror

import (
	"errors"
	"fmt"
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

type E3 struct{
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
	fstest.PrintTestBegin("fserror.IsErrType")
	err := test1()
	fmt.Println(1111, IsErrType[E1](err))
	fmt.Println(2222, IsErrType[E2](err))

	err = test2()
	fmt.Println(3333, IsErrType[E3](err))
	fstest.PrintTestEnd()
}
