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
	"crypto/md5"
	"fmt"
)

// 生成全是大写字母的 MD5 码
func UpperMD5(s string) string {
	return fmt.Sprintf("%X", md5.Sum([]byte(s)))
}

// 生成全是小写的 MD5 码
func LowerMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
