/**
@copyright: fantasysky 2016
@brief: 实现一些 byte 数组相关功能函数
@author: fanky
@version: 1.0
@date: 2019-12-25
**/

package fsbytes

import "fmt"

// 从第 n 个字符开始查找指定字符串，如果存在，则返回第一个匹配的子字符串的位置，不存在，则返回 -1
func IndexN(bs []byte, n int, sub []byte) int {
	count := len(bs)
	findCount := count - len(sub)
	if findCount < 0 || n > findCount {
		return -1
	}
	fmt.Println(findCount, n)
	for idx := n; idx < count; idx++ {
		found := true
		inc := idx
		for _, ch := range sub {
			if ch == bs[inc] {
				inc += 1
			} else {
				found = false
				break
			}
		}

		if found {
			return idx
		}
	}
	return -1
}
