/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: md5 generator
@author: fanky
@version: 1.0
@date: 2022-02-24
**/

package fsmd5

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// 生成全是大写字母的 MD5 码
func UpperMD5(s string) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(s)))
}

// 生成全是小写的 MD5 码
func LowerMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// 对文件计算 MD5
func FileMD5(filename string, upper bool) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	for buf, reader := make([]byte, 65536), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		hash.Write(buf[:n])
	}

	if upper {
		return fmt.Sprintf("%X", hash.Sum(nil)), nil
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
