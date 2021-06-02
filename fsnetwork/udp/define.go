/**
@copyright: fantasysky 2016
@brief: udp server definations
@author: fanky
@version: 1.0
@date: 2019-12-21
**/

package udp

import "net"

// 接收消息函数
type F_Receiver func(*net.UDPAddr, []byte)

// 服务关闭通知函数
type F_Closer func()
