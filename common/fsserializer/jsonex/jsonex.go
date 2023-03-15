/**
@copyright: fantasysky 2016
@brief: jsonex
@author: fanky
@version: 1.0
@date: 2019-06-08
**/

// jsonex 用的是 go 标准库的 json 解释器，只是允许 json 文件插入单行和多行注释
// 注释与 javascript 一致
// 示例请看 jsonex_test.go
package jsonex

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// -------------------------------------------------------------------
// json file reader
// -------------------------------------------------------------------
type s_Reader struct {
	r    *bufio.Reader
	buff bytes.Buffer

	lastComma  bool // 是否暂存有逗号
	lastSlash  bool // 是否暂存有 /
	lastZero   bool // 是否暂存 '0'
	firstValue bool // 是否是值开始
}

func newReader(r io.Reader) *s_Reader {
	return &s_Reader{
		r:    bufio.NewReader(r),
		buff: *bytes.NewBuffer([]byte{}),
	}
}

// ---------------------------------------------------------
// 跳过单行注释
func (this *s_Reader) skipComment() {
	this.r.ReadLine()
}

// 跳过多行注释
func (this *s_Reader) skipMulComment() {
	_, err := this.r.ReadBytes('*')
	if err != nil {
		return
	}
	char, err := this.r.ReadByte()
	if err != nil {
		return
	}
	if char == '/' {
		return
	}
	this.skipMulComment()
}

// 提取字符串内部
func (this *s_Reader) pickString() {
	bs, err := this.r.ReadBytes('"')
	if err != nil {
		return
	}
	this.buff.Write(bs)
	if bytes.HasSuffix(bs, []byte{'\\', '"'}) {
		this.pickString()
	}
}

// 转换十六进制
func (this *s_Reader) transformHexValue(x byte) {
	chars := []byte{}
	for {
		char, err := this.r.ReadByte()
		if err != nil {
			break
		}
		if strings.Contains(" ,\r\n\t]}", string(char)) {
			this.r.UnreadByte()
			break
		} else {
			chars = append(chars, char)
		}
	}
	value, err := strconv.ParseUint(string(chars), 16, 64)
	if err != nil {
		this.buff.Write([]byte{'0', x})
		this.buff.Write(chars)
	} else {
		this.buff.WriteString(fmt.Sprintf("%v", value))
	}
}

// ---------------------------------------------------------
func (this *s_Reader) flushLastChar() {
	if this.lastComma {
		this.buff.WriteByte(',')
	}
	if this.lastZero {
		this.buff.WriteByte('0')
	}
	this.lastComma = false
	this.lastZero = false
}

func (this *s_Reader) filterParse() error {
	for {
		char, err := this.r.ReadByte()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		switch char {
		case '/':
			if this.lastSlash {
				this.skipComment()
			}
			this.lastSlash = !this.lastSlash
			continue
		case '*':
			if this.lastSlash {
				this.skipMulComment()
			}
			this.lastSlash = !this.lastSlash
			continue
		case '"':
			this.flushLastChar()
			this.firstValue = false

			this.buff.WriteByte('"')
			this.pickString()
			continue
		case ':':
			this.buff.WriteByte(char)
			this.firstValue = true
			continue
		case ',':
			this.flushLastChar()
			this.lastComma = true
			this.firstValue = true
			continue
		case '[':
			this.flushLastChar()
			this.firstValue = true
			this.buff.WriteByte('[')
			continue
		case '}', ']':
			if this.lastComma {
				this.lastComma = false // 忽略最后一个逗号
			}
			this.firstValue = false
			this.flushLastChar()
			this.buff.WriteByte(char)
			continue
		case ' ', '\t', '\r', '\n':
			if this.lastZero {
				this.buff.WriteByte('0')
				this.lastZero = false
				this.buff.WriteByte(char)
			}
			continue
		case '0': // 十六进制
			if this.lastComma {
				this.buff.WriteByte(',')
				this.lastComma = false
			}
			if this.lastZero {
				this.buff.WriteString("00")
				this.lastZero = false
			} else {
				this.lastZero = true
			}
			this.firstValue = false
			continue
		case 'x', 'X':
			if this.lastComma {
				this.buff.WriteByte(',')
				this.lastComma = false
			}
			if this.lastZero {
				this.transformHexValue(char)
				this.lastZero = false
			} else {
				this.buff.WriteByte(char)
			}
			this.firstValue = false
			continue
		}
		this.flushLastChar()
		this.firstValue = false
		this.buff.WriteByte(char)
	}
}

// -------------------------------------------------------------------
// package public
// -------------------------------------------------------------------
func Load(path string, inst any) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	reader := newReader(fi)
	err = reader.filterParse()
	if err != nil {
		return err
	}

	return json.Unmarshal(reader.buff.Bytes(), inst)
}

func test(jstr string) (string, error) {
	r := newReader(strings.NewReader(jstr))
	err := r.filterParse()
	if err != nil {
		return "", err
	}
	return r.buff.String(), nil
}
