/**
@copyright: fantasysky 2016
@brief: 实现时间相关函数
@author: fanky
@version: 1.0
@date: 2019-02-21
**/

package fstime

import (
	"strconv"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// 时间段，最大单位 日，最新单位 秒
// -----------------------------------------------------------------------------
type S_Hms struct {
	hs int
	ms int
	ss int
}

func NewHmsFromDuration(du time.Duration) *S_Hms {
	return &S_Hms{
		hs: int(du.Hours()),
		ms: int(du.Minutes()) % 60,
		ss: int(du.Seconds()) % 60,
	}
}

func (this *S_Hms) Days() int {
	return this.hs / 24
}

func (this *S_Hms) Hours() int {
	return this.hs
}

func (this *S_Hms) Minutes() int {
	return this.hs*60 + this.ms
}

func (this *S_Hms) Seconds() int64 {
	return int64(this.Minutes())*60 + int64(this.ss)
}

// -------------------------------------------------------------------
func (this *S_Hms) Hour() int {
	return this.hs % 24
}

func (this *S_Hms) Minute() int {
	return this.ms
}

func (this *S_Hms) Second() int {
	return this.ss
}

// -------------------------------------------------------------------
func (this *S_Hms) Format(fmt string) string {
	fmt = strings.ReplaceAll(fmt, "{D}", strconv.Itoa(this.Days()))
	fmt = strings.ReplaceAll(fmt, "{d}", strconv.Itoa(this.Days()))
	fmt = strings.ReplaceAll(fmt, "{H}", strconv.Itoa(this.Hours()))
	fmt = strings.ReplaceAll(fmt, "{M}", strconv.Itoa(this.Minutes()))
	fmt = strings.ReplaceAll(fmt, "{S}", strconv.FormatInt(this.Seconds(), 10))
	fmt = strings.ReplaceAll(fmt, "{h}", strconv.Itoa(this.Hour()))
	fmt = strings.ReplaceAll(fmt, "{m}", strconv.Itoa(this.ms))
	fmt = strings.ReplaceAll(fmt, "{s}", strconv.Itoa(this.ss))
	return fmt
}
