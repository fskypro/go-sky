package fsos

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestGetDFInfo(t *testing.T) {
	fstest.PrintTestBegin("GetDFInfo")
	defer fstest.PrintTestEnd()

	info, _ := GetDFInfo("/")
	fmt.Printf("%#v\n", info)
	info, _ = GetDFInfo("/opt/freeswitch-1.10.7")
	fmt.Printf("%#v\n", info)
	info, _ = GetDFInfo("/home")
	fmt.Printf("%#v\n", info)
}

func TestGetFileStat(t *testing.T) {
	fstest.PrintTestBegin("GetFileStat")
	defer fstest.PrintTestEnd()

	finfo, err := GetFileStat("/opt/gcc-11.2")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Printf("%#v\n", finfo)
}
