package fsos

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

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
