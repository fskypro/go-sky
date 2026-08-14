/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: parent & child process
@author: fanky
@version: 1.0
@date: 2025-12-10
**/

package fsps

import (
	"os"
	"os/exec"
	"syscall"
)

// -------------------------------------------------------------------
// 父进程退出管道
// -------------------------------------------------------------------
type S_ExitPipe struct {
	pr    *os.File // 父进程读出文件描述符
	pw    *os.File // 父进程写入文件描述符
	order int      // 子进程中对应的文件描述符编号
}

// 创建退出管道
// cmd 为子进程命令对象
func NewExitPipe(cmd *exec.Cmd) (*S_ExitPipe, error) {
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	cmd.ExtraFiles = append(cmd.ExtraFiles, pr)
	return &S_ExitPipe{
		pr:    pr,
		pw:    pw,
		order: len(cmd.ExtraFiles) + 2,
	}, nil
}

// 子进程中对应的文件描述符编号
func (this *S_ExitPipe) Order() int {
	return this.order
}

// 关闭管道
func (this *S_ExitPipe) Close() {
	this.pr.Close()
	this.pw.Close()
}

// -------------------------------------------------------------------
// 子进程侦测父进程关闭管道侦测
// -------------------------------------------------------------------
// 通过读取父进程退出管道，阻塞线程
// 注意：任何情况下读取管道错误，都会停止阻塞
func BlockByExitPipe(fd int) bool {
	f := os.NewFile(uintptr(fd), "parent-pipe") // 对应上面 S_ExitPipe 中的 order
	_, err := f.Stat()
	if err != nil {
		return false
	}
	buf := make([]byte, 1)
	f.Read(buf)
	return true
}

// 检测到父进程管道错误，则给自己发送 kill 信号
func KillWaitExitPipe(fd int) {
	go func() {
		if BlockByExitPipe(fd) {
			syscall.Kill(os.Getpid(), syscall.SIGTERM)
		}
	}()
}
