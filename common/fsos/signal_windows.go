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
)

// 创建退出信号量通道
func MakeExitSignalChan() chan os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	return c
}
