/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: unix dialer creator
@author: fanky
@version: 1.0
@date: 2025-05-28
**/

package unixrpc

import (
	"net"
	"time"
)

func createDialier(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, timeout)
}
