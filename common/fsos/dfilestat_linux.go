/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: file status info
@author: fanky
@version: 1.0
@date: 2022-06-07
**/

// 文件系统使用情况
package fsos

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// -------------------------------------------------------------------
// 挂载点磁盘空间情况
// -------------------------------------------------------------------
type S_DFInfo struct {
	Path        string
	Total       uint64  // 总容量 KB
	Free        uint64  // 剩余容量 KB
	Used        uint64  // 已使用量 KB
	UsedPercent float32 // 已使用率 %
	FreePercent float32 // 剩余使用率 %
}

func GetDFInfo(path string) (*S_DFInfo, error) {
	fs := new(syscall.Statfs_t)
	err := syscall.Statfs(path, fs)
	if err != nil {
		return nil, err
	}
	dfInfo := new(S_DFInfo)
	dfInfo.Path = path
	dfInfo.Total = fs.Blocks * uint64(fs.Bsize) / 1024
	dfInfo.Free = fs.Bfree * uint64(fs.Bsize) / 1024
	dfInfo.Used = dfInfo.Total - dfInfo.Free
	dfInfo.UsedPercent = float32(float64(dfInfo.Used)/float64(dfInfo.Total)) * 100
	dfInfo.FreePercent = float32(float64(dfInfo.Free)/float64(dfInfo.Total)) * 100
	return dfInfo, nil
}

// -------------------------------------------------------------------
// 文件或目录磁盘空间占用情况
// -------------------------------------------------------------------
type S_FileStat struct {
	Path  string
	Total uint64 // 所属磁盘总容量 KB
	Free  uint64 // 所属磁盘剩余容量 KB
	Used  uint64 // 目录或者文件占用的容量 KB
	Files int    // 文件数量
	Dirs  int    // 文件夹数量
}

func GetFileStat(path string) (*S_FileStat, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	finfo := &S_FileStat{
		Path: path,
	}
	abs, err := filepath.Abs(path)
	if err == nil {
		dinfo, err := GetDFInfo(abs)
		if err == nil {
			finfo.Total = dinfo.Total
			finfo.Free = dinfo.Free
		}
	}

	if !fi.IsDir() {
		finfo.Used = uint64(fi.Size() / 1024)
		finfo.Files = 1
		return finfo, nil
	}

	var size uint64
	err = filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("read file %q fail, %v", p, err)
		}
		if fi.IsDir() {
			finfo.Dirs += 1
		} else {
			finfo.Files += 1
		}
		size += uint64(fi.Size())
		return nil
	})
	if err != nil {
		return nil, err
	}

	finfo.Used = size / 1024
	return finfo, nil
}
