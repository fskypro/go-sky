/**
@copyright: fantasysky 2016
@brief: udp argument informations
@author: fanky
@version: 1.0
@date: 2019-12-21
**/

package udp

import (
	"fmt"
	"net"
)

// -------------------------------------------------------------------
// S_UDPInfo
// -------------------------------------------------------------------
type S_UDPInfo struct {
	net.UDPAddr
	BuffSize int
}

// 新建 udp 地址
// add 格式为：ip:port
func NewUDPInfoWithAddr(addr string, buffSize int) (info *S_UDPInfo, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}
	info = &S_UDPInfo{*udpAddr, buffSize}
	if info.BuffSize < 64 {
		info.BuffSize = 64
	}
	return
}

// 新建服务器监听信息
func NewUDPInfo(ip string, port int, buffSize int) (*S_UDPInfo, error) {
	return NewUDPInfoWithAddr(fmt.Sprintf("%s:%d", ip, port), buffSize)
}
