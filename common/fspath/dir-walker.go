/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 高性能文件夹遍历对象，对于巨量文件夹的遍历，减少内存损耗
@author: fanky
@version: 1.0
@date: 2026-04-17
**/

package fspath

import (
	"io"
	"os"
	"path/filepath"
)

type S_FileInfo struct {
	os.FileInfo
	skip     bool // 不要访问其子文件和子文件夹
	FullPath string
}

func (this *S_FileInfo) Escape() {
	this.skip = true
}

func (this *S_FileInfo) HasChildren() (bool, error) {
	if !this.IsDir() {
		return false, nil
	}
	f, err := os.Open(this.FullPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// -------------------------------------------------------------------
// DirWalker
// -------------------------------------------------------------------
type S_DirWalker struct {
	Root      string // 根文目录
	BatchSize int    // 每批次读取子文件的数量
}

func NewDirWalker(root string) *S_DirWalker {
	return &S_DirWalker{
		Root:      root,
		BatchSize: 100,
	}
}

func NewDirWalkerWithBatchSize(root string, batchSize int) *S_DirWalker {
	return &S_DirWalker{
		Root:      root,
		BatchSize: max(1, batchSize),
	}
}

func (this *S_DirWalker) walk(dir string, fun func(string, *S_FileInfo, error) bool) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()

	for {
		names, err := f.Readdirnames(this.BatchSize)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for _, name := range names {
			fullPath := filepath.Join(dir, name)
			info, err := os.Lstat(fullPath)
			if err != nil {
				fun(fullPath, nil, err)
				continue
			}
			winfo := &S_FileInfo{
				FileInfo: info,
				FullPath: fullPath,
			}

			if info.IsDir() {
				if !fun(fullPath, winfo, nil) {
					return nil
				}
				if !winfo.skip {
					this.walk(fullPath, fun)
				}
			} else {
				if !fun(fullPath, winfo, nil) {
					return nil
				}
			}
		}
	}
	return nil
}

// 参数 fun 为遍历文件回调函数：
// fullPath: 子文件/文件夹的全目录
// info    : 文件/文件夹信息
// err     : 在遍历到当前文件时，是否有错误
// 如果 fun 返回 false，则终止正个遍历
// 如果要停止遍历某个文件夹，则可以调用 info.Escape() 方法
func (this *S_DirWalker) Walk(fun func(fullPath string, info *S_FileInfo, err error) bool) error {
	return this.walk(this.Root, fun)
}
