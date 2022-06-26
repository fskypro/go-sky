/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief:
@author: fanky
@version: 1.0
@date: 2022-06-07
**/

package fsos

import (
	"fmt"
	"syscall"
	"unsafe"
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
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return nil, fmt.Errorf("load kernel library fail, %v", err)
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")
	if err != nil {
		return nil, fmt.Errorf("execute kernel library function fail, %v", err)
	}

	totalBytes := int64(0)
	freeBytesAvailable := int64(0)
	freeBytes := int64(0)
	r, a, b := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:"))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&freeBytes)),
		0, 0)

	hdInfo := new(S_HDInfo)
	hdInfo.Path = path
	hdInfo.Total = totalBytes / 1024
	hdInfo.Free = freeBytes / 1024
	hdInfo.Used = hdInfo.Total - hdInfo.Free
	dfInfo.UsedPercent = float32(float64(dfInfo.Used)/float64(dfInfo.Total)) * 100
	dfInfo.FreePercent = float32(float64(dfInfo.Free)/float64(dfInfo.Total)) * 100
	return hdInfo, nil
}
