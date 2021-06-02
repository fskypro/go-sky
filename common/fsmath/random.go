/**
@copyright: fantasysky 2016
@brief: 实现一些随机功能
@author: fanky
@version: 1.0
@date: 2019-01-19
**/

package fsmath

import "math/rand"
import "time"

// RandSrcond 从指定秒数区间中，随机生成一个秒值
//	min, max 不能小于 0，否则 panic
//	随机数取值范围是：[min, max)
func RandSecond(min, max int) time.Duration {
	space := rand.Intn(max - min)
	return time.Second * time.Duration(min+space)
}

// RandMillisecond 从指定的毫秒区间中，随机生成一个毫秒值
//	min，max 不能小于 0，否则 panic
//	随机数取值范围是：[min, max)
func RandMillisecond(min, max int) time.Duration {
	space := rand.Intn(max - min)
	return time.Millisecond * time.Duration(min+space)
}
