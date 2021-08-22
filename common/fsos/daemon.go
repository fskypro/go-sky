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
)

type S_DaemonFile struct {
	FileName string // pid file name
	Pid      int    // current process id
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

	if err = ioutil.WriteFile(file, []byte(strconv.Itoa(df.Pid)), 0440); err != nil {
		df = nil
		err = fmt.Errorf("create pid file fail, %v", err)
		return
	}
	return
}
