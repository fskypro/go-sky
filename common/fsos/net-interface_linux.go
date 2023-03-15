/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: net interface
@author: fanky
@version: 1.0
@date: 2022-08-25
**/

package fsos

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
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
	root := "/sys/class/net"
	infos, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read net interface dir fail, %v", err)
	}
	nis := []*S_NetInterface{}
	for _, info := range infos {
		if info.Name() == "lo" {
			continue
		}
		ni := &S_NetInterface{Name: info.Name()}

		bs, err := ioutil.ReadFile(path.Join(root, info.Name(), "ifindex"))
		if err != nil {
			bs = []byte{}
		}
		ni.IfIndex, _ = strconv.Atoi(strings.TrimSpace(string(bs)))

		bs, err = ioutil.ReadFile(path.Join(root, info.Name(), "address"))
		if err != nil {
			bs = []byte{}
		}
		ni.MacAddr = strings.TrimSpace(string(bs))

		bs, err = ioutil.ReadFile(path.Join(root, info.Name(), "mtu"))
		if err != nil {
			bs = []byte{}
		}
		ni.Mtu, err = strconv.Atoi(strings.TrimSpace(string(bs)))

		bs, err = ioutil.ReadFile(path.Join(root, info.Name(), "operstate"))
		if err != nil {
			bs = []byte{}
		}
		ni.OperState = strings.TrimSpace(string(bs))
		nis = append(nis, ni)
	}

	return nis, nil
}
