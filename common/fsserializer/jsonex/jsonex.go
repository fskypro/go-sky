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

import "io"
import "os"
import "encoding/json"

// -------------------------------------------------------------------
// json file reader
// -------------------------------------------------------------------
type s_Reader struct {
	r    io.Reader
	data []byte
}

func newReader(r io.Reader) *s_Reader {
	return &s_Reader{
		r:    r,
		data: make([]byte, 0),
	}
}

// 读取一个字符
func (this *s_Reader) readByte() (byte, error) {
	var bs = make([]byte, 1)
	n, err := this.r.Read(bs)
	if err == io.EOF || n < 1 {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return bs[0], nil
}

// 跳过单行注释
func (this *s_Reader) skipSingleComment() error {
	for {
		char, err := this.readByte()
		if err != nil {
			return err
		}

		if char == '\n' || char == '\r' || char == 0 {
			return nil
		}
	}
	return nil
}

// 跳过多行注释
func (this *s_Reader) skipMultiComment() error {
	star := false
	for {
		char, err := this.readByte()
		if err != nil {
			return err
		}
		if char == 0 {
			return nil
		}

		// 注释结束
		if star && char == '/' {
			return nil
		}

		if char == '*' {
			star = true
		}
	}
	return nil
}

func (this *s_Reader) read() error {
	quote := false // 进入字符串内部

	lastBSlash := false // 上一个字符是否是反斜杠
	lastSlash := false  // 上一个字符是否是斜杠
	lastComma := false  // 上一个有效字符是否是逗号

	// 清除斜杠标记
	var clearSlash = func() {
		if lastSlash {
			this.data = append(this.data, '/')
			lastSlash = false
		}
	}

	// 清除逗号标记
	var clearComma = func() {
		if lastComma {
			this.data = append(this.data, ',')
			lastComma = false
		}
	}

	for {
		char, err := this.readByte()
		if err != nil {
			return err
		}
		if char == 0 {
			return nil
		}

		// 字符串内部
		if quote {
			if char == '"' {
				if !lastBSlash { // 离开字符串内部
					quote = false
				}
			} else if char == '\\' {
				if lastBSlash { // 连续两个反斜杠，表示输出一个反斜杠
					lastBSlash = false
				} else {
					lastBSlash = true
				}
			} else {
				lastBSlash = false
			}

			lastSlash = false
			this.data = append(this.data, char)
			continue
		}

		// 进入字符串内部
		if char == '"' {
			quote = true

			clearSlash()
			clearComma()
			this.data = append(this.data, char)
			continue
		}

		// 记录下斜杠
		if char == '/' {
			if lastSlash { // 连续两个斜杠
				err = this.skipSingleComment() // 跳过单行注释
				if err != nil {
					return err
				}
				lastSlash = false
			} else {
				lastSlash = true
			}

			continue
		}

		// 多行注释
		if char == '*' && lastSlash {
			err = this.skipMultiComment()
			if err != nil {
				return err
			}
			lastSlash = false
			continue
		}

		// 忽略掉所有不在字符串内部的空白字符
		if char == ' ' || char == '\t' || char == '\r' || char == '\n' {
			clearSlash()
			continue
		}

		// 去掉最后一个元素的逗号
		if char == '}' || char == ']' {
			lastComma = false
			clearSlash()
			this.data = append(this.data, char)
			continue
		}
		if lastComma {
			lastComma = false
			this.data = append(this.data, ',', char)
			continue
		}
		if char == ',' { // 逗号暂时忽略（等下一个字符如果不是 ]、} 再取，如果是 ]、} 则忽略掉）
			lastComma = true
			clearSlash()
			continue
		}

		clearSlash()
		clearComma()
		this.data = append(this.data, char)
	}
}

// -------------------------------------------------------------------
// package public
// -------------------------------------------------------------------
func Load(path string, inst interface{}) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()

	reader := newReader(fi)
	err = reader.read()
	if err != nil {
		return err
	}

	return json.Unmarshal(reader.data, inst)
}
