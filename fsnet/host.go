/**
@copyright: fantasysky 2016
@brief: host
@author: fanky
@version: 1.0
@date: 2020-06-10
**/

package fsnet

import (
	"fmt"
	"strconv"
	"strings"
)

// -------------------------------------------------------------------
// 主机结构
// -------------------------------------------------------------------
type S_Host struct {
	Addr string `json:"addr"`
	Port uint16 `json:"port"`
}

// 新建 S_Host
// addrPort 格式为：地址或域名:端口号
// 如果传入的地址格式不正确，则返回 nil
func NewHost(addrPort string) *S_Host {
	ap := strings.Split(addrPort, ":")
	if len(ap) != 2 {
		return nil
	}
	port := ap[1]
	iport, err := strconv.ParseUint(port, 0, 16)
	if err != nil {
		return nil
	}
	return &S_Host{
		Addr: ap[0],
		Port: uint16(iport),
	}
}

func (this *S_Host) String() string {
	return fmt.Sprintf("%s:%d", this.Addr, this.Port)
}

func (this *S_Host) Clone() *S_Host {
	return &S_Host{
		Addr: this.Addr,
		Port: this.Port,
	}
}
