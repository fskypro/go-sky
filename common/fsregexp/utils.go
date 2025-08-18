/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 常用正则表达式校验
@author: fanky
@version: 1.0
@date: 2025-07-17
**/

package fsregexp

import "regexp"

// -------------------------------------------------------------------
// 是否是 IPv4
// -------------------------------------------------------------------
var reIPv4 = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])$`)

func IsIPv4[IP ~string](ip IP) bool {
	return reIPv4.MatchString(string(ip))
}
