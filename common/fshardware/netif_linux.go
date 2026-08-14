/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: net interface
@author: fanky
@version: 1.0
@date: 2022-08-25
**/

package fshardware

import (
	"fmt"
	"net"
	"os"
)

// 判断是否为虚拟网卡
func isPhysicalInterface(iface net.Interface) bool {
	if iface.Name == "lo" || iface.HardwareAddr.String() == "00:00:00:00:00:00" {
		return false
	}
	_, err := os.Stat(fmt.Sprintf("/sys/class/net/%s/device", iface.Name))
	return err == nil
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
