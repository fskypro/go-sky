/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: mainboard info
@author: fanky
@version: 1.0
@date: 2026-02-04
**/

package fshardware

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetMotherboardInfo() (*S_MainBaseBoard, error) {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("exec command %q fail, %v", "system_pro", err)
	}

	info := new(S_MainBaseBoard)

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Model Name:") {
			info.Manufacturer = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Model Identifier:") {
			info.Product = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Model Number:") {
			info.Version = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Serial Number") {
			segs := strings.Split(line, ":")
			if len(segs) > 1 {
				info.SerialNumber = strings.TrimSpace(segs[len(segs)-1])
			}
		}
	}
	return info, nil
}
