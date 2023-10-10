/**
@copyright: exht 2015
@website: http://www.exht.com
@brief:
@author: fanky
@version: 1.0
@date: 2023-09-18
**/

package fsexcel

import "math"

// 将 excel 表格列的 26 进制转换为十进制索引
// A == 0; Z == 25; AA == 26
func ColOrderToIndex(order string) int {
	const base float64 = 26
	var result int = 0
	letters := []byte(order)
	exp := len(letters) - 1
	if exp < 0 { return 0 }
	for _, c := range letters {
		n := int(c - 'A' + 1)
		result += n * int(math.Pow(base, float64(exp)))
		exp--
	}
	return result - 1
}

// 将索引转换为 excel 表格列序号
func IndexToColOrder(index int) string {
	result := ""
	for index >= 0 {
		remainder := index % 26
		c := byte(remainder) + 'A'
		result = string(c) + result
		index = index/26 - 1
	}
	return result
}
