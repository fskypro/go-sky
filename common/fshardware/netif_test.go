package fshardware

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestGetNetInterfaceInfos(t *testing.T) {
	fstest.PrintTestBegin("GetNetInterfaceInfos")
	defer fstest.PrintTestEnd()
	nis, _ := GetPhysicalNetInterfaces()
	for _, ni := range nis {
		fmt.Printf("%#v\n", ni)
	}
}
