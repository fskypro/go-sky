package fsos

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestGetCpuInfo(t *testing.T) {
	fstest.PrintTestBegin("GetCpuInfo")
	defer fstest.PrintTestEnd()

	info, err := GetCpuInfo()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	for _, c := range info.CpuCores {
		fmt.Printf("%#v\n", c)
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
