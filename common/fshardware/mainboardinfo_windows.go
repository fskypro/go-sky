//go:build windows

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

	"github.com/StackExchange/wmi"
)

type s_Win32BaseBoard struct {
	Manufacturer string
	Product      string
	Version      string
	SerialNumber string
}

func GetMotherboardInfo() (*S_MainBaseBoard, error) {
	var boards []s_Win32BaseBoard
	query := "SELECT Manufacturer, Product, Version, SerialNumber FROM Win32_BaseBoard"
	err := wmi.Query(query, &boards)
	if err != nil || len(boards) == 0 {
		return nil, fmt.Errorf("exec query(%q) fail, %v", query, err)
	}
	return &S_MainBaseBoard{
		Manufacturer: boards[0].Manufacturer,
		Product:      boards[0].Product,
		Version:      boards[0].Version,
		SerialNumber: boards[0].SerialNumber,
	}, nil
}
