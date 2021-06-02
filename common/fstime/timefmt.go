/**
@copyright: fantasysky 2016
@brief: 时间格式化
@author: fanky
@version: 1.0
@date: 2019-02-21
**/

package fstime

import "time"

// Str2DateTime 将格式为 “2006-01-02 15:04:05” 的时间字符串转换为 UTC 日期时间
func Str2DateTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

// Str2LocDateTime 将格式为 “2006-01-02 15:04:05” 的时间字符串转换为 Local 日期时间
func Str2LocDateTime(str string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", str, time.Now().Location())
}

// DateTime2Str 将指定时间转换为格式为“2006-01-02 15:04:05”的字符串
func DateTime2Str(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
