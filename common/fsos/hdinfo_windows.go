/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: disk information for windows
@author: fanky
@version: 1.0
@date: 2022-06-07
**/

package fsos

import (
	"github.com/shirou/gopsutil/disk"
)

type S_HDInfo struct {
	Path        string
	Total       uint64  // 总容量 KB
	Free        uint64  // 剩余容量 KB
	Used        uint64  // 已使用量 KB
	UsedPercent float32 // 已使用率 %
	FreePercent float32 // 剩余使用率 %
}

func GetHDInfo(path string) (*S_HDInfo, error) {
	usage, err := disk.UsageWithContext(nil, path)
	if err != nil { return nil, err }
	return &S_HDInfo{
		Path: path,
		Total: usage.Total,
		Free: usage.Free,
		Used: usage.Used,
		UsedPercent: float32(usage.UsedPercent),
	}, nil
}
