/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: ping tool base on ping command
@author: fanky
@version: 1.0
@date: 2023-12-04
**/

package fsping

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var ErrUnsupported = errors.New("unsupported in this operation system")
var ErrNoPingCmd = errors.New("ping command is not exists")

func getArgs(timeout int) []string {
	switch runtime.GOOS {
	case "linux":
		return strings.Fields(fmt.Sprintf("-w %d -c 1", timeout))
	case "darwin":
		return strings.Fields(fmt.Sprintf("-t %d -c 1", timeout))
	case "windows":
		return strings.Fields(fmt.Sprintf("-w %d -n 1", timeout*1000))
	}
	return nil
}

func Ping4(addr string, timeout int) (bool, error) {
	args := getArgs(timeout)
	if args == nil {
		return false, ErrUnsupported
	}
	args = append(args, addr)
	cmd := exec.Command("ping", args...)
	_, err := cmd.CombinedOutput()
	if err == nil {
		return true, nil
	}
	if strings.Contains(err.Error(), "no command") {
		return false, ErrNoPingCmd
	}
	return false, err
}
