/**
@copyright: fantasysky 2016
@brief: name-text extention
@author: fanky
@version: 1.0
@date: 2020-02-20
**/

package fsxml

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// -------------------------------------------------------------------
// inner functions
// -------------------------------------------------------------------
var _spaces = [4]rune{' ', '\t', '\r', '\n'}

func _split(text string) []string {
	isSpace := func(ch rune) bool {
		for _, space := range _spaces {
			if space == ch {
				return true
			}
		}
		return false
	}

	item := make([]rune, 0)
	items := make([]string, 0)
	for _, ch := range text {
		if isSpace(ch) {
			if len(item) > 0 {
				items = append(items, string(item))
				item = []rune{}
			}
		} else {
			item = append(item, ch)
		}
	}
	if len(item) > 0 {
		items = append(items, string(item))
	}
	return items
}

// -------------------------------------------------------------------
// node extensions
// -------------------------------------------------------------------
// 将 node 的 text 内容转换为 int8 返回
// 如果转换失败，则 err 为：
// 	RangeError：数值超出范围
// 	TypeError：不是数值类型
func (this *s_NameText) Int8() (int8, error) {
	v, err := strconv.ParseInt(this.text, 0, 8)
	if err != nil {
		return 0, convertError(err, this.text, "int8")
	}
	return int8(v), nil
}

// 将 node 的 text 内容转换为 int8 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsInt8(def int8) int8 {
	v, err := strconv.ParseInt(this.text, 0, 8)
	if err != nil {
		return def
	}
	return int8(v)
}

// 将 node 的 text 内容转换为一组 int8
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) Int8s() ([]int8, error) {
	values := make([]int8, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 8)
		if err != nil {
			return nil, convertError(err, item, "int8")
		} else {
			values = append(values, int8(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 int8 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsInt8s(def []int8) []int8 {
	values := make([]int8, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 8)
		if err != nil {
			return def
		} else {
			values = append(values, int8(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 int8 值
func (this *s_NameText) SetInt8(value int8) {
	this.text = strconv.Itoa(int(value))
}

// 将 node 的 text 设置为以空格分开的一组 int8 值
func (this *s_NameText) SetInt8s(values []int8) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = strconv.Itoa(int(value))
		} else {
			this.text += " " + strconv.Itoa(int(value))
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 uint 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) UInt() (uint, error) {
	v, err := strconv.ParseUint(this.text, 0, 64)
	if err != nil {
		return 0, convertError(err, this.text, "uint8")
	}
	return uint(v), nil
}

// 将 node 的 text 内容转换为 uint 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsUInt(def uint) uint {
	v, err := strconv.ParseUint(this.text, 0, 64)
	if err != nil {
		return def
	}
	return uint(v)
}

// 将 node 的 text 内容转换为一组 uint
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) UInts() ([]uint, error) {
	values := make([]uint, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 64)
		if err != nil {
			return nil, convertError(err, item, "uint")
		} else {
			values = append(values, uint(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 uint 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsUInts(def []uint) []uint {
	values := make([]uint, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 64)
		if err != nil {
			return def
		} else {
			values = append(values, uint(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 uint 值
func (this *s_NameText) SetUInt(value uint) {
	this.text = fmt.Sprintf("%d", value)
}

// 将 node 的 text 设置为以空格分开的一组 uint8 值
func (this *s_NameText) SetUInts(values []uint) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%d", value)
		} else {
			this.text = fmt.Sprintf("%s %d", this.text, value)
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 uint8 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) UInt8() (uint8, error) {
	v, err := strconv.ParseUint(this.text, 0, 8)
	if err != nil {
		return 0, convertError(err, this.text, "uint8")
	}
	return uint8(v), nil
}

// 将 node 的 text 内容转换为 uint8 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsUInt8(def uint8) uint8 {
	v, err := strconv.ParseUint(this.text, 0, 8)
	if err != nil {
		return def
	}
	return uint8(v)
}

// 将 node 的 text 内容转换为一组 uint8
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) UInt8s() ([]uint8, error) {
	values := make([]uint8, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 8)
		if err != nil {
			return nil, convertError(err, item, "uint8")
		} else {
			values = append(values, uint8(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 uint8 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsUInt8s(def []uint8) []uint8 {
	values := make([]uint8, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 8)
		if err != nil {
			return def
		} else {
			values = append(values, uint8(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 uint8 值
func (this *s_NameText) SetUInt8(value uint8) {
	this.text = strconv.Itoa(int(value))
}

// 将 node 的 text 设置为以空格分开的一组 uint8 值
func (this *s_NameText) SetUInt8s(values []uint8) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = strconv.Itoa(int(value))
		} else {
			this.text += " " + strconv.Itoa(int(value))
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 int16 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) Int16() (int16, error) {
	v, err := strconv.ParseInt(this.text, 0, 16)
	if err != nil {
		return 0, convertError(err, this.text, "int16")
	}
	return int16(v), nil
}

// 将 node 的 text 内容转换为 int16 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsInt16(def int16) int16 {
	v, err := strconv.ParseInt(this.text, 0, 16)
	if err != nil {
		return def
	}
	return int16(v)
}

// 将 node 的 text 内容转换为一组 int16
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) Int16s() ([]int16, error) {
	values := make([]int16, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 16)
		if err != nil {
			return nil, convertError(err, item, "int16")
		} else {
			values = append(values, int16(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 int16 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsInt16s(def []int16) []int16 {
	values := make([]int16, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 16)
		if err != nil {
			return def
		} else {
			values = append(values, int16(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 int16 值
func (this *s_NameText) SetInt16(value int16) {
	this.text = strconv.Itoa(int(value))
}

// 将 node 的 text 设置为以空格分开的一组 int16 值
func (this *s_NameText) SetInt16s(values []int16) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = strconv.Itoa(int(value))
		} else {
			this.text += " " + strconv.Itoa(int(value))
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 uint16 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) UInt16() (uint16, error) {
	v, err := strconv.ParseUint(this.text, 0, 16)
	if err != nil {
		return 0, convertError(err, this.text, "uint16")
	}
	return uint16(v), nil
}

// 将 node 的 text 内容转换为 uint16 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsUInt16(def uint16) uint16 {
	v, err := strconv.ParseUint(this.text, 0, 16)
	if err != nil {
		return def
	}
	return uint16(v)
}

// 将 node 的 text 内容转换为一组 uint16
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) UInt16s() ([]uint16, error) {
	values := make([]uint16, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 16)
		if err != nil {
			return nil, convertError(err, item, "uint16")
		} else {
			values = append(values, uint16(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 uint16 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsUInt16s(def []uint16) []uint16 {
	values := make([]uint16, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 16)
		if err != nil {
			return def
		} else {
			values = append(values, uint16(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 uint16 值
func (this *s_NameText) SetUInt16(value uint16) {
	this.text = strconv.Itoa(int(value))
}

// 将 node 的 text 设置为以空格分开的一组 uint16 值
func (this *s_NameText) SetUInt16s(values []uint16) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = strconv.Itoa(int(value))
		} else {
			this.text += " " + strconv.Itoa(int(value))
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 int 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) Int() (int, error) {
	v, err := strconv.ParseInt(this.text, 0, 32)
	if err != nil {
		return 0, convertError(err, this.text, "int32")
	}
	return int(v), nil
}

// 将 node 的 text 内容转换为 int32 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsInt(def int) int {
	v, err := strconv.ParseInt(this.text, 0, 32)
	if err != nil {
		return def
	}
	return int(v)
}

// 将 node 的 text 内容转换为一组 int32
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) Ints() ([]int, error) {
	values := make([]int, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 32)
		if err != nil {
			return nil, convertError(err, item, "int32")
		} else {
			values = append(values, int(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 int32 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsInts(def []int) []int {
	values := make([]int, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 32)
		if err != nil {
			return def
		} else {
			values = append(values, int(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 int32 值
func (this *s_NameText) SetInt(value int) {
	this.text = fmt.Sprintf("%d", value)
}

// 将 node 的 text 设置为以空格分开的一组 int32 值
func (this *s_NameText) SetInts(values []int) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%d", value)
		} else {
			this.text += fmt.Sprintf(" %d", value)
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 int32 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) Int32() (int32, error) {
	v, err := strconv.ParseInt(this.text, 0, 32)
	if err != nil {
		return 0, convertError(err, this.text, "int32")
	}
	return int32(v), nil
}

// 将 node 的 text 内容转换为 int32 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsInt32(def int32) int32 {
	v, err := strconv.ParseInt(this.text, 0, 32)
	if err != nil {
		return def
	}
	return int32(v)
}

// 将 node 的 text 内容转换为一组 int32
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) Int32s() ([]int32, error) {
	values := make([]int32, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 32)
		if err != nil {
			return nil, convertError(err, item, "int32")
		} else {
			values = append(values, int32(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 int32 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsInt32s(def []int32) []int32 {
	values := make([]int32, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 32)
		if err != nil {
			return def
		} else {
			values = append(values, int32(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 int32 值
func (this *s_NameText) SetInt32(value int32) {
	this.text = fmt.Sprintf("%d", value)
}

// 将 node 的 text 设置为以空格分开的一组 int32 值
func (this *s_NameText) SetInt32s(values []int32) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%d", value)
		} else {
			this.text += fmt.Sprintf(" %d", value)
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 uint32 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) UInt32() (uint32, error) {
	v, err := strconv.ParseUint(this.text, 0, 32)
	if err != nil {
		return 0, convertError(err, this.text, "int32")
	}
	return uint32(v), nil
}

// 将 node 的 text 内容转换为 uint32 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsUInt32(def uint32) uint32 {
	v, err := strconv.ParseUint(this.text, 0, 32)
	if err != nil {
		return def
	}
	return uint32(v)
}

// 将 node 的 text 内容转换为一组 uint32
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) UInt32s() ([]uint32, error) {
	values := make([]uint32, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 32)
		if err != nil {
			return nil, convertError(err, item, "uint32")
		} else {
			values = append(values, uint32(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 uint32 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsUInt32s(def []uint32) []uint32 {
	values := make([]uint32, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 32)
		if err != nil {
			return def
		} else {
			values = append(values, uint32(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 uint32 值
func (this *s_NameText) SetUInt32(value uint32) {
	this.text = fmt.Sprintf("%d", value)
}

// 将 node 的 text 设置为以空格分开的一组 int32 值
func (this *s_NameText) SetUInt32s(values []uint32) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%d", value)
		} else {
			this.text += fmt.Sprintf(" %d", value)
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 int64 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) Int64() (int64, error) {
	v, err := strconv.ParseInt(this.text, 0, 64)
	if err != nil {
		return 0, convertError(err, this.text, "int64")
	}
	return int64(v), nil
}

// 将 node 的 text 内容转换为 int64 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsInt64(def int64) int64 {
	v, err := strconv.ParseInt(this.text, 0, 64)
	if err != nil {
		return def
	}
	return int64(v)
}

// 将 node 的 text 内容转换为一组 int64
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) Int64s() ([]int64, error) {
	values := make([]int64, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 64)
		if err != nil {
			return nil, convertError(err, item, "int64")
		} else {
			values = append(values, int64(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 int64 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsInt64s(def []int64) []int64 {
	values := make([]int64, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseInt(item, 0, 64)
		if err != nil {
			return def
		} else {
			values = append(values, int64(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 uint64 值
func (this *s_NameText) SetInt64(value int64) {
	this.text = fmt.Sprintf("%d", value)
}

// 将 node 的 text 设置为以空格分开的一组 int64 值
func (this *s_NameText) SetInt64s(values []int64) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%d", value)
		} else {
			this.text += fmt.Sprintf(" %d", value)
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 uint64 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) UInt64() (uint64, error) {
	v, err := strconv.ParseUint(this.text, 0, 64)
	if err != nil {
		return 0, convertError(err, this.text, "uint64")
	}
	return uint64(v), nil
}

// 将 node 的 text 内容转换为 uint64 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsUInt64(def uint64) uint64 {
	v, err := strconv.ParseUint(this.text, 0, 64)
	if err != nil {
		return def
	}
	return uint64(v)
}

// 将 node 的 text 内容转换为一组 uint64
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) UInt64s() ([]uint64, error) {
	values := make([]uint64, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 64)
		if err != nil {
			return nil, convertError(err, item, "uint64")
		} else {
			values = append(values, uint64(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 uint64 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsUInt64s(def []uint64) []uint64 {
	values := make([]uint64, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseUint(item, 0, 64)
		if err != nil {
			return def
		} else {
			values = append(values, uint64(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 uint64 值
func (this *s_NameText) SetUInt64(value uint64) {
	this.text = fmt.Sprintf("%d", value)
}

// 将 node 的 text 设置为以空格分开的一组 uint64 值
func (this *s_NameText) SetUInt64s(values []uint64) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%d", value)
		} else {
			this.text += fmt.Sprintf(" %d", value)
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 float32 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) Float32() (float32, error) {
	v, err := strconv.ParseFloat(this.text, 32)
	if err != nil {
		return 0, convertError(err, this.text, "float32")
	}
	return float32(v), nil
}

// 将 node 的 text 内容转换为 float32 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsFloat32(def float32) float32 {
	v, err := strconv.ParseFloat(this.text, 32)
	if err != nil {
		return def
	}
	return float32(v)
}

// 将 node 的 text 内容转换为一组 float32
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) Float32s() ([]float32, error) {
	values := make([]float32, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseFloat(item, 32)
		if err != nil {
			return nil, convertError(err, item, "float32")
		} else {
			values = append(values, float32(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 float32 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsFloat32s(def []float32) []float32 {
	values := make([]float32, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseFloat(item, 32)
		if err != nil {
			return def
		} else {
			values = append(values, float32(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 float32 值
func (this *s_NameText) SetFloat32(value float32) {
	this.text = fmt.Sprintf("%f", value)
}

// 将 node 的 text 设置为以空格分开的一组 uint64 值
func (this *s_NameText) SetFloat32s(values []float32) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%f", value)
		} else {
			this.text += fmt.Sprintf(" %f", value)
		}
	}
}

// ---------------------------------------------------------
// 将 node 的 text 内容转换为 float64 返回
// 如果转换失败，则 err 为：
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
func (this *s_NameText) Float64() (float64, error) {
	v, err := strconv.ParseFloat(this.text, 64)
	if err != nil {
		return 0, convertError(err, this.text, "float64")
	}
	return float64(v), nil
}

// 将 node 的 text 内容转换为 float64 返回
// 如果转换失败，则返回默认值 def
func (this *s_NameText) AsFloat64(def float64) float64 {
	v, err := strconv.ParseFloat(this.text, 64)
	if err != nil {
		return def
	}
	return float64(v)
}

// 将 node 的 text 内容转换为一组 float64
// 	ErrRange：数值超出范围
// 	ErrSyntax：不是数值类型
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) Float64s() ([]float64, error) {
	values := make([]float64, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseFloat(item, 64)
		if err != nil {
			return nil, convertError(err, item, "float64")
		} else {
			values = append(values, float64(v))
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 float64 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsFloat64s(def []float64) []float64 {
	values := make([]float64, 0)
	items := _split(this.text)
	for _, item := range items {
		v, err := strconv.ParseFloat(item, 64)
		if err != nil {
			return def
		} else {
			values = append(values, float64(v))
		}
	}
	return values
}

// 将 node 的 text 设置为一个 float64 值
func (this *s_NameText) SetFloat64(value float64) {
	this.text = fmt.Sprintf("%f", value)
}

// 将 node 的 text 设置为以空格分开的一组 uint64 值
func (this *s_NameText) SetFloat64s(values []float64) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%f", value)
		} else {
			this.text += fmt.Sprintf(" %f", value)
		}
	}
}

// ---------------------------------------------------------
func (this *s_NameText) Bool() (bool, error) {
	text := strings.ToLower(this.text)
	if text == "true" || text == "1" {
		return true, nil
	} else if text == "false" || text == "0" {
		return false, nil
	}
	err := fmt.Errorf("can't convert %q to a bool value", this.text)
	return false, convertError(err, this.text, "bool")
}

// 将 node 的 text 根据特征换换为布尔型返回
// 返回：
//	如果 text 为 true/1 则返回 true（不区分大小写）
//	如果 text 为 false/0 则返回 false（不区分大小写）
//	其他情况返回 def
func (this *s_NameText) AsBool(def bool) bool {
	text := strings.ToLower(this.text)
	if text == "true" || text == "1" {
		return true
	} else if text == "false" || text == "0" {
		return false
	}
	return def
}

// 将 node 的 text 内容转换为一组 bool
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) Bools() ([]bool, error) {
	values := make([]bool, 0)
	items := _split(this.text)
	for _, item := range items {
		item = strings.ToLower(item)
		if item == "true" || item == "1" {
			values = append(values, true)
		} else if item == "false" || item == "0" {
			values = append(values, false)
		} else {
			err := fmt.Errorf("can't convert %q to a bool value", item)
			return nil, convertError(err, item, "uint")
		}
	}
	return values, nil
}

// 将 node 的 text 内容转换为一组 bool 返回
// 如果转换失败，则返回默认值 def
// 注意：text 中每个元素必须用空白字符(空格、\t、换行)分开
func (this *s_NameText) AsBools(def []bool) []bool {
	values := make([]bool, 0)
	items := _split(this.text)
	for _, item := range items {
		item = strings.ToLower(item)
		if item == "true" || item == "1" {
			values = append(values, true)
		} else if item == "false" || item == "0" {
			values = append(values, false)
		} else {
			return def
		}
	}
	return values
}

// 将 node 的 text 设置为一个 bool 值
func (this *s_NameText) SetBool(value bool) {
	this.text = fmt.Sprintf("%v", value)
}

// 将 node 的 text 设置为以空格分开的一组 uint8 值
func (this *s_NameText) SetBools(values []bool) {
	this.text = ""
	for index, value := range values {
		if index == 0 {
			this.text = fmt.Sprintf("%v", value)
		} else {
			this.text = fmt.Sprintf("%s %v", this.text, value)
		}
	}
}

// ---------------------------------------------------------
// 返回以空格作为分隔符的字符串列表
func (this *s_NameText) AsTexts() []string {
	if this.text == "" {
		return []string{}
	}
	re := regexp.MustCompile(`\s+`)
	return re.Split(this.text, -1)
}

// 返回以指定正则模式分隔字符串为分隔符的一组文本内容
func (this *s_NameText) AsReTexts(sp string) ([]string, error) {
	if this.text == "" {
		return []string{}, nil
	}
	re, err := regexp.Compile(sp)
	if err != nil {
		return []string{}, fmt.Errorf("error splitter pattern %q", sp)
	}
	return re.Split(this.text, -1), nil
}

// 将 node 的 text 设置为以指定分隔符合并的一组字符串
func (this *s_NameText) SetTexts(values []string, sp string) {
	this.text = strings.Join(values, sp)
}
