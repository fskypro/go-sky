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
func (this *S_DaemonFile) KillLast(signal syscall.Signal) error {
	if this.LastPid == 0 {
		return nil
	}
	return syscall.Kill(this.LastPid, signal)
}

// 创建进程ID文件
// existsFail 为 true 时，如果有同进程存在，则放弃创建，并返回错误
func CreateDaemonFile(file string) (df *S_DaemonFile, err error) {
	df = &S_DaemonFile{
		FileName: file,
		Pid:      os.Getpid(),
	}
	root := path.Dir(file)
	if e := os.MkdirAll(root, 0440); e != nil {
		err = fmt.Errorf("create pid file fail, %v", e)
		return
	}
	data, err := ioutil.ReadFile(file)
	if err == nil {
		pid, _ := strconv.Atoi(string(data))
		df.LastPid = pid
	}

	if err = ioutil.WriteFile(file, []byte(strconv.Itoa(df.Pid)), 0440); err != nil {
		df = nil
		err = fmt.Errorf("create pid file fail, %v", err)
		return
	}
	return
}
