//go:build darwin || linux
// +build darwin linux

/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: daemon
@author: fanky
@version: 1.0
@date: 2021-07-10
**/

package fsos

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"syscall"
)

type S_DaemonFile struct {
	FileName string // pid file name
	Pid      int    // current process id
	LastPid  int    // old process id
}

// 杀死旧进程
// 如果进程不存在，则返回：err == syscall.ESRCH
func (this *S_DaemonFile) KillLast(signal syscall.Signal) error {
	if this.LastPid == 0 {
		return nil
	}
	return syscall.Kill(this.LastPid, signal)
}

// 创建进程ID文件
// 如果进程ID文件已经存在：
//   cover == true ：则覆盖掉原来的进程文件
//   cover == false：则不覆盖直接返回(err = nil; df = new-daemon-file-info)
func CreateDaemonFile(file string, cover bool) (df *S_DaemonFile, err error) {
	df = &S_DaemonFile{
		FileName: file,
		Pid:      os.Getpid(),
	}

	// 创建文件路径
	root := path.Dir(file)
	if e := os.MkdirAll(root, 0440); e != nil {
		err = fmt.Errorf("create pid file path fail, %v", e)
		return
	}

	// 读取旧文件
	data, err := os.ReadFile(file)
	if err == nil {
		df.LastPid, _ = strconv.Atoi(string(data))
		if !cover {
			return
		}
	}

	if err = ioutil.WriteFile(file, []byte(strconv.Itoa(df.Pid)), 0440); err != nil {
		err = fmt.Errorf("create pid file fail, %v", err)
	}
	return
}
