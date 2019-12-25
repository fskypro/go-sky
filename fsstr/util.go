/**
@copyright: fantasysky 2016
@brief: 实现一些字符串相关功能函数
@author: fanky
@version: 1.0
@date: 2019-12-25
**/

package fsstr

import "fsky.pro/fsbytes"
import "fsky.pro/fsstr/convert"

// 从第 n 个字符开始查找指定字符串，如果存在，则返回第一个匹配的子字符串的位置，不存在，则返回 -1
func IndexN(str string, n int, sub string) int {
	return fsbytes.IndexN(convert.String2Bytes(str), n, convert.String2Bytes(sub))
}
