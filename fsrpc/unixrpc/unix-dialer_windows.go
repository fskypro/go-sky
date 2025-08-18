/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: dialer creator in windows
@author: fanky
@version: 1.0
@date: 2025-05-28
**/

package unixrpc

import (
	"net"
	"path/filepath"
	"time"

	"github.com/natefinch/npipe"
)

func createDialier(addr string, timeout time.Duration) (net.Conn, error) {
	addr = filepath.Base(addr)
	return npipe.DialTimeout(`\\.\pipe\`+addr, timeout)
}
