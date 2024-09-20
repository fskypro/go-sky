/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: os utils
@author: fanky
@version: 1.0
@date: 2022-07-07
**/

package fsos

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 按提示获取一组输入字符串
func ScanInputs(asks [][2]string) []string {
	out := make([]string, len(asks))
	input := bufio.NewScanner(os.Stdin)
	for index, ask := range asks {
	L:
		fmt.Printf("%s: ", ask[0])
		input.Scan()
		if input.Text() == "" {
			if ask[1] == "" {
				goto L
			} else {
				out[index] = ask[1]
			}
		} else {
			out[index] = input.Text()
		}
	}
	return out
}

// 判断是否 go run 启动的运行
func IsGoRun() (bool, error) {
	tempDir := os.TempDir()
	execDir, err := filepath.Abs(os.Args[0])
	if err != nil {
		return false, err
	}
	return strings.HasPrefix(execDir, tempDir), nil
}
