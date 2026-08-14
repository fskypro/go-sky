/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: get cpu info in darwin
@author: fanky
@version: 1.0
@date: 2026-02-04
**/

package fshardware

import (
	"os/exec"
	"strconv"
	"strings"

	"fsky.pro/fsky"
)

func GetCpuInfo() (*S_CpuInfo, error) {
	// 获取基本信息
	vendor := getSysctl("machdep.cpu.vendor")
	family := getSysctl("machdep.cpu.family")
	model := getSysctl("machdep.cpu.model")
	modelName := getSysctl("machdep.cpu.brand_string")

	// 核心数
	physicalCores := getSysctlInt("hw.physicalcpu")
	logicalCores := getSysctlInt("hw.logicalcpu")

	// 频率（Hz转MHz）
	freq := getSysctlInt("hw.cpufrequency")
	freqMHz := freq / 1000000

	// 缓存（字节转KB）
	l1i := getSysctlInt("hw.l1icachesize") / 1024
	l1d := getSysctlInt("hw.l1dcachesize") / 1024
	l1 := max(l1i, l1d) // 取指令缓存大小

	l2 := getSysctlInt("hw.l2cachesize") / 1024
	l3 := getSysctlInt("hw.l3cachesize") / 1024

	// 为每个逻辑核心创建记录
	var cores []*S_CpuCore
	for i := range logicalCores {
		core := &S_CpuCore{
			Order:      i + 1,
			VendorID:   fsky.IfElse(vendor == "", "Apple", vendor),
			CpuFamily:  family,
			Model:      model,
			ModelName:  modelName,
			MicroCode:  "", // macOS不提供
			CpuMHz:     freqMHz,
			CacheSize:  l3, // 主缓存用L3
			PhysicalID: 0,  // macOS通常单路
			CpuCores:   physicalCores,
			Flags:      []string{}, // macOS不提供
			Bugs:       "",         // macOS不提供
			AddrSize:   "",         // macOS不提供
			L1Cache:    l1,
			L2Cache:    l2,
			L3Cache:    l3,
		}
		cores = append(cores, core)
	}

	return &S_CpuInfo{
		Cores: cores,
	}, nil
}

func getSysctl(key string) string {
	cmd := exec.Command("sysctl", "-n", key)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func getSysctlInt(key string) int {
	val := getSysctl(key)
	i, _ := strconv.Atoi(val)
	return i
}
