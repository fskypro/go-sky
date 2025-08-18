/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: unix listener creator
@author: fanky
@version: 1.0
@date: 2025-05-28
**/

package unixrpc

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func createListener(addr string) (net.Listener, error) {
	dir := filepath.Dir(addr)
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		if err := os.Mkdir(dir, 0700); err != nil {
			return nil, fmt.Errorf("create run path %q fail, %v", dir, err)
		}
	}
	os.Remove(addr)

	ln, err := net.Listen("unix", addr)
	if err != nil {
		return nil, fmt.Errorf("listen unix address %q fail, %v", addr, err)
	}
	return ln, err
}
