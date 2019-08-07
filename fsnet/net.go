/**
* @brief: net.go
* @copyright: 2016 fantasysky
* @author: fanky
* @version: 1.0
* @date: 2018-12-27
 */

package fsnet

import "net"
import "strconv"

// GetFreePort 获取一个空闲的端口号
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
