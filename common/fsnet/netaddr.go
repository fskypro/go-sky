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
type S_Addr struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// 新建 S_Addr
// addrPort 格式为：地址或域名:端口号
// 如果传入的地址格式不正确，则返回 nil
func NewAddr(addr string) *S_Addr {
	hp := strings.Split(addr, ":")
	if len(hp) != 2 {
		return nil
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return nil
	}
	return &S_Addr{
		Host: hp[0],
		Port: port,
	}
}

func (this *S_Addr) Addr() string {
	return fmt.Sprintf("%s:%d", this.Host, this.Port)
}

func (this *S_Addr) String() string {
	return fmt.Sprintf(`"%s:%d"`, this.Host, this.Port)
}

func (this *S_Addr) Clone() *S_Addr {
	return &S_Addr{
		Host: this.Host,
		Port: this.Port,
	}
}
