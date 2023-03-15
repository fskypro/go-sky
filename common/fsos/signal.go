/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: signal
@author: fanky
@version: 1.0
@date: 2021-07-10
**/

package fsos

import (
	"os"
	"os/signal"
	"syscall"
)

// 创建退出信号量通道
func MakeExitSignalChan() chan os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c,
		syscall.SIGKILL, // 无条件结束程序(不能被捕获、阻塞或忽略)，“kill -9” 时触发
		syscall.SIGHUP,  // 终端控制进程结束(终端连接断开)
		syscall.SIGINT,  // 用户发送 INTR 字符，Ctrl+C 时触发
		syscall.SIGTERM, // 结束程序(可以被捕获、阻塞或忽略)，“kill” 时触发
		syscall.SIGQUIT, // 用户发送 QUIT 字符(Ctrl+/)触发
		syscall.SIGABRT, // 调用 abort 函数触发
		syscall.SIGSTOP, // 停止进程(不能被捕获、阻塞或忽略)
	)
	return c
}
