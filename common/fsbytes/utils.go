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

// 将字节数转换为合理单位的大小
// 第二个参数为保留几位小数
func BytesTo(size int64, decs int) string {
	const (
		_          = iota             // 忽略第一个值 (0)
		KB float64 = 1 << (10 * iota) // 1 KB = 1024
		MB                            // 1 MB = 1024 * 1024
		GB                            // 1 GB = 1024 * 1024 * 1024
		TB                            // 1 TB
	)

	switch {
	case size >= int64(TB):
		return fmt.Sprintf(fmt.Sprintf("%%.%dfTB", decs), float64(size)/TB)
	case size >= int64(GB):
		return fmt.Sprintf(fmt.Sprintf("%%.%dfGB", decs), float64(size)/GB)
	case size >= int64(MB):
		return fmt.Sprintf(fmt.Sprintf("%%.%dfMB", decs), float64(size)/MB)
	case size >= int64(KB):
		return fmt.Sprintf(fmt.Sprintf("%%.%dfKB", decs), float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}
