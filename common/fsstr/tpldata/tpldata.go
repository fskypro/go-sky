/**
@copyright: fantasysky 2016
@brief: 实现模板数据基类
@author: fanky
@version: 1.0
@date: 2019-03-20
**/

package tpldata

import "strings"

type S_TemplateData struct{}

// -------------------------------------------------------------------
// common
// -------------------------------------------------------------------
// 产生一组整数迭代
func (this *S_TemplateData) RangeInt_(start, end, step int) chan int {
	ch := make(chan int)
	go func() {
		for i := start; i < end; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

// -------------------------------------------------------------------
// math
// -------------------------------------------------------------------
// 参数相加
func (this *S_TemplateData) AddInt_(vs ...int) int {
	total := 0
	for _, v := range vs {
		total += v
	}
	return total
}

// 减法
func (this *S_TemplateData) SubInt_(nu1, nu2 int) int {
	return nu1 - nu2
}

// 是否是单数
func (this *S_TemplateData) IsSingularInt_(v int) bool {
	return (v % 2) == 1
}

// 是否是双数
func (this *S_TemplateData) IsDualInt_(v int) bool {
	return (v % 2) == 0
}

// -------------------------------------------------------------------
// string
// -------------------------------------------------------------------
// 判断 str 是否以 sub 开头
func (this *S_TemplateData) StartsWith_(str string, sub string) bool {
	return strings.HasPrefix(str, sub)
}

// 判断 str 是否以 sub 结尾
func (this *S_TemplateData) EndsWith_(str string, sub string) bool {
	return strings.HasSuffix(str, sub)
}
