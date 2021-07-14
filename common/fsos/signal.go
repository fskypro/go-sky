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
		syscall.SIGHUP,  // 终端控制进程结束(终端连接断开)
		syscall.SIGINT,  // 用户发送INTR字符(Ctrl+C)触发
		syscall.SIGTERM, // kill
		syscall.SIGQUIT)
	return c
}
