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
	"os"
	"strings"
)

func GetMotherboardInfo() (*S_MainBaseBoard, error) {
	info := new(S_MainBaseBoard)

	// /sys/class/dmi/id/ 包含主板信息
	files := []string{
		"/sys/class/dmi/id/board_vendor",
		"/sys/class/dmi/id/board_name",
		"/sys/class/dmi/id/board_version",
		"/sys/class/dmi/id/board_serial",
	}

	strInfos := []string{}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read mainboard info file %q fail, %v", path, err)
		}
		strInfos = append(strInfos, string(data))
	}

	info.Manufacturer = strings.TrimSpace(strInfos[0])
	info.Product = strings.TrimSpace(strInfos[1])
	info.Version = strings.TrimSpace(strInfos[2])
	info.SerialNumber = strings.TrimSpace(strInfos[3])

	return info, nil
}
