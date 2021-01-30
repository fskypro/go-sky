/**
@copyright: fantasysky 2016
@brief: base struct of tag and attribute
@author: fanky
@version: 1.0
@date: 2020-02-19
**/

package xmlex

import (
	"bytes"
	"fmt"
	"strings"
)
import (
	"fsky.pro/fsserializer/xmlex/internal/xml"
	"fsky.pro/fsstr/convert"
)

type s_NameText struct {
	name string // 名称（格式为 name 或 spacename:name）
	text string
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
// 验证是否是合法的 tag/attr name
func isValidName(name string) bool {
	bxml := []byte(fmt.Sprintf("<%s>", name))
	r := bytes.NewReader(bxml)
	doc := xml.NewDecoder(r)
	_, err := doc.Token()
	return err == nil
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 获取全名称（格式为 name 或 spacename:name）
func (this *s_NameText) Name() string {
	return this.name
}

// 获取 tag/attr 名称的名字空间部分（即冒号“:” 前的部分）
func (this *s_NameText) Space() string {
	ss := strings.Split(this.name, ":")
	if len(ss) < 2 {
		return ""
	}
	return ss[0]
}

// 获取 tag/attr 名称的元素名称部分（即冒号“:”后面的部分）
func (this *s_NameText) Local() string {
	ss := strings.Split(this.name, ":")
	if len(ss) < 2 {
		return this.name
	}
	return ss[1]
}

func (this *s_NameText) Text() string {
	return this.text
}

// 返回字节形式的节点内容
func (this *s_NameText) TextBytes() []byte {
	return convert.String2Bytes(this.text)
}

// 设置节点内容
func (this *s_NameText) SetText(text string) {
	this.text = text
}

// 设置节点内容
func (this *s_NameText) SetTextBytes(btext []byte) {
	this.text = convert.Bytes2String(btext)
}
