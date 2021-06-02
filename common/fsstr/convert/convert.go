/**
@copyright: fantasysky 2016
@brief: 字符串转换
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package convert

import "unsafe"
import "reflect"

// SliceByteToString 将 []byte 字符数组直接转换为 string
func Bytes2String(bs []byte) string {
	bsHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	sHeader := reflect.StringHeader{
		Data: bsHeader.Data,
		Len:  bsHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&sHeader))
}

// StringToSliceByte 将 string 内部字符数组直接转换为 []byte
func String2Bytes(s string) []byte {
	sHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bsHeader := reflect.SliceHeader{
		Data: sHeader.Data,
		Len:  sHeader.Len,
		Cap:  sHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bsHeader))
}
