/**
@copyright: fantasysky 2016
@brief: udp argument informations
@author: fanky
@version: 1.0
@date: 2019-12-21
**/

package udp

import "net"

type S_UDPInfo struct {
	net.UDPAddr
	BuffSize int
}

func NewUDPInfoAddr(addr string, buffSize int) (info *S_UDPInfo, err error) {
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

func NewUDPInfo(ip string, port int, buffSize int) *S_UDPInfo {
	info := &S_UDPInfo{}
	info.UDPAddr.IP = net.ParseIP(ip)
	info.UDPAddr.Port = port
	info.BuffSize = buffSize
	if info.BuffSize < 64 {
		info.BuffSize = 64
	}
	return info
}
