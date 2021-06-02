/**
@copyright: fantasysky 2016
@brief: 序列化 json 对象
@author: fanky
@version: 1.0
@date: 2019-06-06
**/

package freejson

import (
	"fmt"
	"io"
	"strings"

	"fsky.pro/fsenv"
	"fsky.pro/fsstr/convert"
)

// json 转移符
var _trans2str = map[byte][]byte{
	'\\': []byte{'\\', '\\'},
	'\b': []byte{'\\', 'b'},
	'\r': []byte{'\\', 'r'},
	'\n': []byte{'\\', 'n'},
	'\f': []byte{'\\', 'f'},
	'"':  []byte{'\\', '"'},
}

// ----------------------------------------------------------------------------
// S_FmtInfo
// ----------------------------------------------------------------------------
type S_FmtInfo struct {
	Newline      string // 换行符
	Indent       string // 缩进符号
	IndentLayers int    // 缩进几层，为 0 则表示所有 object 都缩进
	IndentList   bool   // 列表是否缩进（每个元素一行）
}

func NewFmtInfo() *S_FmtInfo {
	return &S_FmtInfo{
		Newline:      fsenv.Endline,
		Indent:       "    ",
		IndentLayers: 0,
		IndentList:   false,
	}
}

// ----------------------------------------------------------------------------
// s_Writer
// ----------------------------------------------------------------------------
type s_Writer struct {
	w          io.Writer
	subWriters map[JType]func(*s_Writer, I_Value) error
	layers     int // 当前嵌套层数

	isFormat     bool   // 是否格式化
	newline      string // 换行符
	indent       string // 缩进符
	indentLayers int    // 缩进的层数
	indentList   bool   // 列表是否缩进
}

func newWriter(w io.Writer, fmtInfo *S_FmtInfo) *s_Writer {
	writer := &s_Writer{
		w:      w,
		layers: 0,

		isFormat:     false,
		newline:      "",
		indent:       "",
		indentLayers: -1,
		indentList:   false,
	}
	writer.subWriters = map[JType]func(*s_Writer, I_Value) error{
		TObject:  _writeObject,
		TList:    _writeList,
		TString:  _writeString,
		TInt64:   _writeOrigin,
		TUInt64:  _writeOrigin,
		TFloat64: _writeOrigin,
		TBool:    _writeOrigin,
		TNull:    _writeOrigin,
	}

	// 转移格式化参数
	writer.isFormat = fmtInfo != nil &&
		fmtInfo.Newline != "" &&
		fmtInfo.Indent != "" &&
		fmtInfo.IndentLayers >= 0
	if writer.isFormat {
		writer.newline = fmtInfo.Newline
		writer.indent = fmtInfo.Indent
		writer.indentLayers = fmtInfo.IndentLayers
		writer.indentList = fmtInfo.IndentList
	}
	return writer
}

// -------------------------------------------------------------------
// module private
// -------------------------------------------------------------------
// 写入一个字符
func (this *s_Writer) _writeByte(b byte) error {
	_, err := this.w.Write([]byte{b})
	return err
}

// 写入一个字符串
func (this *s_Writer) _writeString(str string) error {
	_, err := this.w.Write(convert.String2Bytes(str))
	return err
}

// 写入字符数组
func (this *s_Writer) _writeBytes(bs []byte) error {
	_, err := this.w.Write(bs)
	return err
}

// 写入 object 的换行
func (this *s_Writer) _newlineObject(layers int) error {
	if this.isFormat && (this.indentLayers == 0 || this.layers <= this.indentLayers) {
		err := this._writeString(this.newline)
		err = this._writeString(strings.Repeat(this.indent, layers))
		return err
	}
	return nil
}

// 写入 list 的换行
func (this *s_Writer) _newlineList(layers int) error {
	if this.indentList && (this.indentLayers == 0 || this.layers <= this.indentLayers) {
		err := this._writeString(this.newline)
		err = this._writeString(strings.Repeat(this.indent, layers))
		return err
	}
	return nil
}

// --------------------------------------------------------
// 写对象
func _writeObject(this *s_Writer, value I_Value) error {
	this._writeByte('{')
	this.layers += 1

	first := true
	var err error
	value.AsObject().For(func(k string, v I_Value) bool {
		if !first {
			err = this._writeByte(',')
		}
		first = false

		// 写键
		err = this._newlineObject(this.layers)
		if this.isFormat {
			err = this._writeString(fmt.Sprintf("%q: ", k))
		} else {
			err = this._writeString(fmt.Sprintf("%q:", k))
		}

		// 写值
		err = this.subWriters[v.Type()](this, v)
		return err == nil
	})
	if err == nil {
		err = this._newlineObject(this.layers - 1)
		err = this._writeByte('}')
	}
	this.layers -= 1
	return err
}

// 写列表
func _writeList(this *s_Writer, value I_Value) error {
	this._writeByte('[')
	this.layers += 1

	first := true
	var err error
	value.AsList().For(func(index int, elem I_Value) bool {
		if !first {
			err = this._writeByte(',')
		}
		first = false

		err = this._newlineList(this.layers)
		this.subWriters[elem.Type()](this, elem)
		return err == nil
	})

	err = this._newlineList(this.layers - 1)
	err = this._writeByte(']')
	this.layers -= 1
	return err
}

// 写数值
func _writeString(this *s_Writer, value I_Value) error {
	err := this._writeByte('"')
	for _, char := range value.AsString().value {
		bs, ok := _trans2str[char]
		if ok {
			err = this._writeBytes(bs)
		} else {
			err = this._writeByte(char)
		}
		if err != nil {
			return err
		}
	}
	err = this._writeByte('"')
	return err
}

// 写其他输出与源格式相同的值
func _writeOrigin(this *s_Writer, value I_Value) error {
	return this._writeString(fmt.Sprintf("%v", value))
}

// -------------------------------------------------------------------
// package public
// -------------------------------------------------------------------
func (this *s_Writer) Write(value I_Value) error {
	return this.subWriters[value.Type()](this, value)
}
