/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: net interface
@author: fanky
@version: 1.0
@date: 2024-03-11
**/

package fshardware

import (
	"net"
	"strings"

	"github.com/StackExchange/wmi"
)

type s_Win32NetworkAdapter struct {
	MACAddress  string
	PNPDeviceID string
}

var physicalMacs map[string]bool

func init() {
	physicalMacs = make(map[string]bool)

	var adapters []s_Win32NetworkAdapter
	query := "SELECT MACAddress, PNPDeviceID FROM Win32_NetworkAdapter WHERE MACAddress IS NOT NULL"
	err := wmi.Query(query, &adapters)
	if err == nil {
		for _, adapter := range adapters {
			// PCI设备才是真正的物理网卡
			if strings.HasPrefix(strings.ToUpper(adapter.PNPDeviceID), "PCI\\") {
				mac := strings.ToLower(strings.ReplaceAll(adapter.MACAddress, ":", "-"))
				physicalMacs[mac] = true
			}
		}
	}
}

func isPhysicalInterface(iface net.Interface) bool {
	if iface.HardwareAddr.String() == "00:00:00:00:00:00" {
		return false
	}
	if iface.Flags&net.FlagLoopback != 0 || len(iface.HardwareAddr) == 0 {
		return false
	}

	mac := strings.ToLower(strings.ReplaceAll(iface.HardwareAddr.String(), ":", "-"))
	return physicalMacs[mac]
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
