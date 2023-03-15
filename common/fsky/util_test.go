package fsky

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

type A struct {
	Name  string
	Value int
}

type B struct {
	As []*A
}

func TestDeepCopy(t *testing.T) {
	fstest.PrintTestBegin("DeepCopy")

	a := &A{
		Name:  "xxxx",
		Value: 100,
	}

	b := &B{
		As: []*A{a},
	}

	bb := new(B)
	err := DeepCopy(bb, b)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	b.As[0].Name = "yyyy"
	b.As[0].Value = 200
	fmt.Println(bb.As[0])

	fstest.PrintTestEnd()
}
