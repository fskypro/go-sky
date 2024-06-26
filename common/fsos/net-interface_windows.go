/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: net interface
@author: fanky
@version: 1.0
@date: 2024-03-11
**/

package fsos

import (
	"fmt"
	"strings"
	"github.com/shirou/gopsutil/net"
)

// -------------------------------------------------------------------
// mac address
// -------------------------------------------------------------------
// interface name map mac address
type S_NetInterface struct {
	Name      string
	IfIndex   int
	MacAddr   string
	Mtu       int
	OperState string
}

func GetNetInterfaceInfos() ([]*S_NetInterface, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("get interfaces fail, %v", err)
	}
	nis := []*S_NetInterface{}
	for _, netInterface := range netInterfaces {
		ni := &S_NetInterface{
			Name: netInterface.Name,
			IfIndex: netInterface.Index,
			MacAddr: netInterface.HardwareAddr,
			Mtu: netInterface.MTU,
			OperState: strings.Join(netInterface.Flags, ","),
		}
		nis = append(nis, ni)
	}
	return nis, nil
}
