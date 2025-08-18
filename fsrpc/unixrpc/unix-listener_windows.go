/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: listener creator in windows
@author: fanky
@version: 1.0
@date: 2025-05-28
**/

package unixrpc

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/natefinch/npipe"
)

func createListener(addr string) (net.Listener, error) {
	addr = filepath.Base(addr)
	ln, err := npipe.Listen(`\\.\pipe\` + addr)
	if err != nil {
		return nil, fmt.Errorf("listen unix address %q fail, %v", addr, err)
	}
	return ln, nil
}
