/**
* @brief: net.go
* @copyright: 2016 fantasysky
* @author: fanky
* @version: 1.0
* @date: 2018-12-27
 */

package fsnet

import (
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// GetFreeTCPPort 获取一个空闲的端口号，找不到则返回 0
func GetFreeTCPPort() int {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0
	}
	defer ln.Close()

	addr := ln.Addr().String()
	_, strPort, err := net.SplitHostPort(addr)
	if err != nil {
		return 0
	}
	p, err := strconv.Atoi(strPort)
	if err != nil {
		return 0
	}
	return p
}

// RandomFreeTCPPort
// 在指定范围内随机一个空闲端口号(包括 max)
func RandomFreeTCPPort(min, max int) int {
	if max < min {
		min, max = max, min
	}
	if min > 65535 {
		return 0
	}
	if max <= 0 {
		return 0
	}
	if max > 65535 {
		max = 65535
	}
	if min <= 0 {
		min = 1
	}

	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	var halfChoice func(int, int) int
	halfChoice = func(start, end int) int {
		if start > end {
			return 0
		}
		p := start + seed.Intn(end-start+1)
		if ln, e := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(p)); e == nil {
			ln.Close()
			return p
		}
		// 从左半部分取
		if seed.Intn(2) == 0 {
			if port := halfChoice(start, p-1); port > 0 {
				return port
			}
			// 从右半部分取
			return halfChoice(p+1, max)
		} else {
			if port := halfChoice(p+1, max); port > 0 {
				return port
			}
			return halfChoice(start, p-1)
		}
	}
	return halfChoice(min, max)
}

// -------------------------------------------------------------------
// IP 地址转换函数
// -------------------------------------------------------------------
// 将四字节整数转换为字符串形式 IP 地址
// 类似 C 语言的 inet_ntoa
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// 将 IP 地址转换为四字节整数
// 类似 C 语言的 inet_aton
// 注意：如果是不合法的 ip 地址，返回 -1
func InetAtoN(ip string) int64 {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return -1
	}
	ret := big.NewInt(0)
	ret.SetBytes(netIP.To4())
	return ret.Int64()
}

// --------------------------------------------------------
// 将 IP 地址和端转换为 UINT64 值
// 参数 addr 格式可以为 aa.bb.cc.dd:port 或者 aa.bb.cc.dd
// 如果是前者，则参数 port 不起作用
// 返回：(aa << 24) + (bb << 16) + (cc << 8) + dd + (port << 32)
//	注意：
//		如果地址部分为空字符串，则地址将被置为 0.0.0.0
//		如果解释失败，则返回 -1
func AddrToInt64(addr string, port uint16) int64 {
	var a string
	p := int64(port)
	ap := strings.Split(addr, ":")
	if len(ap) == 0 {
		a = "0.0.0.0"
	} else if len(ap) == 1 {
		a = ap[0]
	} else if len(ap) == 2 {
		a = ap[0]
		up, err := strconv.ParseUint(ap[1], 10, 16)
		if err != nil {
			return -1
		}
		p = int64(up)
	}
	ip := InetAtoN(a)
	if ip == -1 {
		return -1
	}
	return ip + (p << 32)
}

// 将数值形式的 IP 地址，转换为字符串形式，格式为 aa.bb.cc.dd:port
// 如果数值小于
func Int64ToAddr(iaddr int64) string {
	const minPort int64 = int64(1 << 32)

	if iaddr < minPort { // 没有端口
		return fmt.Sprintf("%d.%d.%d.%d",
			byte(iaddr>>24),
			byte(iaddr>>16),
			byte(iaddr>>8),
			uint16(iaddr))
	}
	return fmt.Sprintf("%d.%d.%d.%d:%d",
		byte(iaddr>>24),
		byte(iaddr>>16),
		byte(iaddr>>8),
		byte(iaddr),
		uint16(iaddr>>32))
}
