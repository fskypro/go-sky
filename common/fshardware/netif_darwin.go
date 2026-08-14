/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: net interface
@author: fanky
@version: 1.0
@date: 2026-02-05
**/

package fshardware

import "net"

func isPhysicalInterface(iface net.Interface) bool {
	if iface.Flags&net.FlagLoopback != 0 || len(iface.HardwareAddr) == 0 {
		return false
	}

	// MAC 地址本地管理位检查
	// 虚拟网卡的 MAC 地址第一个字节的第 2 位通常为 1（如 4e, 36, 42 的二进制第 2 位都是 1）
	return iface.HardwareAddr[0]&0x02 == 0
}

// 获取所有物理网卡信息
func GetPhysicalNetInterfaces() ([]T_NetInterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ifs := []T_NetInterfaceInfo{}
	for _, iface := range interfaces {
		if isPhysicalInterface(iface) {
			ifs = append(ifs, iface)
		}
	}
	return ifs, nil
}
