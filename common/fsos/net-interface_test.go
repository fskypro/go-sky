package fsos

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestGetNetInterfaceInfos(t *testing.T) {
	fstest.PrintTestBegin("GetNetInterfaceInfos")
	defer fstest.PrintTestEnd()
	infos, _ := GetNetInterfaceInfos()
	for _, info := range infos {
		fmt.Printf("%+v\n", info)
	}
}
