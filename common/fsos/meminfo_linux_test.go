package fsos

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestGetMemInfo(t *testing.T) {
	fstest.PrintTestBegin("GetMemInfo")
	defer fstest.PrintTestEnd()

	memInfo, err := GetMemInfo()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	fmt.Printf("%#v\n", memInfo)
}
