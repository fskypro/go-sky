package fsos

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestGetCpuInfo(t *testing.T) {
	fstest.PrintTestBegin("GetCpuInfo")
	defer fstest.PrintTestEnd()
	info , err := GetCpuInfo()
	if err != nil {
		fmt.Printf("get cpu info fail, %v\n", err)
		return
	}
	for _, core := range info.Cores {
		fmt.Printf("core %d: %#v\n", core.Order, core)
	}
}

func TestGetCpuStat(t *testing.T) {
	fstest.PrintTestBegin("GetCpuStat")
	defer fstest.PrintTestEnd()

	stat, err := GetCpuStat()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	fmt.Printf("total: used=%f%%, free=%f%%\n", stat.UsedPercent, stat.FreePercent)
	for _, coreStat := range stat.CoreStats {
		fmt.Printf("%#v\n", coreStat)
	}
}

