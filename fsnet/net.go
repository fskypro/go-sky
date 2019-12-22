/**
* @brief: net.go
* @copyright: 2016 fantasysky
* @author: fanky
* @version: 1.0
* @date: 2018-12-27
 */

package fsnet

import "fmt"
import "net"
import "strconv"
import "strings"

// -------------------------------------------------------------------
// 主机结构
// -------------------------------------------------------------------
type S_Host struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

func (this *S_Host) GetAddr() string {
	return fmt.Sprintf("%s:%d", this.IP, this.Port)
}

// -------------------------------------------------------------------
// GetFreePort 获取一个空闲的端口号
// -------------------------------------------------------------------
func GetFreePort() (port int, err error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()

	addr := ln.Addr().String()
	_, strPort, err := net.SplitHostPort(addr)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strPort)
}

// -------------------------------------------------------------------
// IP 地址转换函数
// -------------------------------------------------------------------
// 将 IP 地址和端转换为 UINT64 值
// 参数 addr 格式可以为 aa.bb.cc.dd:port 或者 aa.bb.cc.dd
// 如果是前者，则参数 port 不起作用
// 返回：(aa << 24) + (bb << 16) + (cc << 8) + dd + (port << 32)
//	注意：如果解释失败，则返回 -1
func AddrToInt64(addr string, port uint16) int64 {
	ap := strings.Split(addr, ":")
	if len(ap) == 0 {
		return -1
	}
	segs := strings.Split(ap[0], ".")
	if len(segs) != 4 {
		return -1
	}
	aa, err := strconv.ParseUint(segs[0], 10, 8)
	if err != nil {
		return -1
	}
	bb, err := strconv.ParseUint(segs[1], 10, 8)
	if err != nil {
		return -1
	}
	cc, err := strconv.ParseUint(segs[2], 10, 8)
	if err != nil {
		return -1
	}
	dd, err := strconv.ParseUint(segs[3], 10, 8)
	if err != nil {
		return -1
	}
	vaddr := (aa << 24) + (bb << 16) + (cc << 8) + dd

	// 地址中带有端口
	if len(ap) > 1 {
		vport, err := strconv.ParseUint(ap[1], 10, 16)
		if err != nil {
			return -1
		}
		return int64(vaddr + (vport << 32))
	}

	// 端口号放在 port 参数中
	if port > 0 {
		return int64(vaddr + uint64(port)<<32)
	}

	return int64(vaddr)
}

// 将数值形式的 IP 地址，转换为字符串形式，格式为 aa.bb.cc.dd:port
// 如果数值小于
func Int64ToAddr(vaddr int64) string {
	const minPort int64 = int64(1 << 32)

	aa := (vaddr >> 24) & int64(255)
	bb := (vaddr >> 16) & int64(255)
	cc := (vaddr >> 8) & int64(255)
	dd := vaddr & int64(255)
	if vaddr < minPort { // 没有端口
		return fmt.Sprintf("%d.%d.%d.%d", aa, bb, cc, dd)
	}
	port := vaddr >> 32
	return fmt.Sprintf("%d.%d.%d.%d:%d", aa, bb, cc, dd, port)
}
