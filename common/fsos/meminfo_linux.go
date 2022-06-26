/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: memory info
@author: fanky
@version: 1.0
@date: 2022-06-06
**/

package fsos

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type S_MemInfo struct {
	Total uint64 // KB
	Free  uint64 // KB
	Used  uint64 // KB

	SwapTotal uint64 // KB
	SwapFree  uint64 // KB
	SwapUsed  uint64 // KB

	UsedPercent     float32 // %
	SwapUsedPercent float32 // %
}

func GetMemInfo() (*S_MemInfo, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("read memory file fail, %v", err)
	}
	memInfo := new(S_MemInfo)
	prefixs := []string{
		"MemTotal",
		"MemFree",
		"SwapTotal",
		"SwapFree",
	}
	count := len(prefixs)

	reader := bufio.NewReader(file)
	for {
		data, _, err := reader.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("read memory file fail, %v", err)
		}
		kv := strings.Split(string(data), ":")
		if len(kv) != 2 {
			continue
		}
		key, value := kv[0], kv[1]

		for _, prefix := range prefixs {
			if prefix == key {
				goto L
			}
		}
		continue

	L:
		vv := strings.Split(strings.TrimSpace(value), " ")
		if len(vv) < 1 {
			continue
		}
		v, _ := strconv.ParseUint(vv[0], 10, 64)
		switch key {
		case "MemTotal":
			memInfo.Total = v
		case "MemFree":
			memInfo.Free = v
		case "SwapTotal":
			memInfo.SwapTotal = v
		case "SwapFree":
			memInfo.SwapFree = v
		}
		count--
		if count <= 0 {
			break
		}
	}
	if memInfo.Total <= 0 {
		return nil, errors.New("meminfo file is invalid")
	}
	memInfo.UsedPercent = float32(float64(memInfo.Total-memInfo.Free)/float64(memInfo.Total)) * 100
	if memInfo.SwapTotal > 0 {
		memInfo.SwapUsedPercent = float32(float64(memInfo.SwapTotal-memInfo.SwapFree)/float64(memInfo.SwapTotal)) * 100
	}
	memInfo.Used = memInfo.Total - memInfo.Free
	memInfo.SwapUsed = memInfo.SwapTotal - memInfo.SwapFree
	return memInfo, nil
}

// -------------------------------------------------------------------
// 用下面的方法也可以获得内存信息，但是发现并不是很准确
// -------------------------------------------------------------------
/*
func GetMemInfo() (*S_MemInfo, error) {
	stat := new(runtime.MemStats)
	runtime.ReadMemStats(stat)
	memInfo := new(S_MemInfo)

	sysInfo := new(syscall.Sysinfo_t)
	err := syscall.Sysinfo(sysInfo)
	if err != nil {
		return nil, err
	}
	memInfo.Total = sysInfo.Totalram / 1024
	memInfo.Free = sysInfo.Freeram / 1024
	memInfo.Used = memInfo.Total - memInfo.Free
	memInfo.UsedPercent = float32(float64(memInfo.Used)/float64(memInfo.Total)) * 100
	return memInfo, nil
}*/
