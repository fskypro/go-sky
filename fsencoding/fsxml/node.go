/**
@copyright: fantasysky 2016
@brief: xml node
@author: fanky
@version: 1.0
@date: 2020-01-31
**/

package fsxml

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

import (
	"fsky.pro/fserror"
	"fsky.pro/fsstr/convert"
)

// 检索子节点路径中的 tag
type s_TagSearch struct {
	tag       string
	index     int
	attrName  string
	attrValue string
}

// 支持：
// xx/yy/zz
// xx/yy[1]/zz
// xx/yy[a="v"]/zz
// xx/yy[a='v']/zz
func _firstSearchPath(path string) (tag *s_TagSearch, tail string, err error) {
	path = strings.TrimLeft(path, " ")
	path = strings.TrimLeft(path, "\t")
	path = strings.TrimLeft(path, "/")

	tag = new(s_TagSearch)
	tmp := make([]byte, 0)
	inSubscript := false
	inAttrValue := false
	attrValueChar := byte(' ')

	var idx int
	var ch byte
L:
	for idx, ch = range convert.String2Bytes(path) {
		switch ch {
		case '"', '\'', '`':
			// 当前在双引号/单引号属性值区段内
			if attrValueChar == ch {
				// 当前双引号是转义字符后的，因此不是结束双引号
				if idx > 1 && path[idx-1] == '\\' { // 转义字符
					tmp = append(tmp, ch)
				} else {
					// 结束属性值区段
					attrValueChar = ' '
					inAttrValue = false
					tag.attrValue = string(tmp)
					tmp = []byte{}
				}
				continue L
			}
			// 当前已经在单引号/双引号属性区段内
			if attrValueChar == '\'' ||
				attrValueChar == '"' ||
				attrValueChar == '`' {
				tmp = append(tmp, ch)
				continue L
			}

			// 当前在属性值区段，但不在双引号/单引号区段或单引号区段内
			if inAttrValue {
				tmp = bytes.TrimSpace(tmp)
				// 如果 “=” 后面不紧接双引号或单引号，则认为该属性值不用双引号或单引号括回来
				if len(tmp) > 0 {
					tmp = append(tmp, ch)
				} else {
					// 进入双引号属性值区段
					attrValueChar = ch
				}
			} else {
				// 当前不在下标区段内，不允许出现双引号或单引号
				err = fmt.Errorf("xml path is not allow to be contain char '%c'", ch)
				return
			}
		case '[':
			// 起始中括在属性值区段内，则该起始中括属于属性值的一部分
			if inAttrValue {
				tmp = append(tmp, ch)
				continue L
			}
			// 已经在下标区段内，仍然出现前中括号，则，这是不合法的路径
			if inSubscript {
				err = fmt.Errorf("char '%c' is not allow to be in xml path subscript.", ch)
				return
			}
			// tag 部分结束，标记进入下标区段
			tag.tag = string(bytes.TrimSpace(tmp))
			tmp = []byte{}
			inSubscript = true
		case ']':
			// 遇到后中括，但是，当前如果当前在属性值区段内，则认为后中括属于属性值的一部分
			if attrValueChar != ' ' {
				tmp = append(tmp, ch)
				continue L
			}
			// 没有下标起始括符就出现了结束符号，这是错误的
			if !inSubscript {
				err = fmt.Errorf("char '%c' is not allow to be in xml path.", ch)
				return
			}
			// 结束下标区段
			inSubscript = false

			// 如果当前在属性值区段，但是又没有双引号或单引号起始，
			// 则该属性值是没双引号或单引号括回的
			if inAttrValue {
				// 如果在双引号或单引号结束的时候已经设置，就不需要再这里设置了
				tag.attrValue = string(tmp)
				tmp = []byte{}
				inAttrValue = false
				continue L
			}

			tmp = bytes.TrimSpace(tmp)
			// 如果前面没出现过属性键值下标，则为索引下标
			if tag.attrName == "" {
				index, e := strconv.Atoi(string(tmp))
				if e != nil { // 错误的索引下标
					err = fmt.Errorf("index subscript in xml path must be an intager, but not %q.", string(tmp))
					return
				} else {
					tag.index = index
				}
				tmp = []byte{}
			}

			// 如果前面已经设置过属性键值下标，则 tmp 一定是空，
			// 否则就是类似于这种：tagxx[k="v" 44]，这是错误的
			if len(tmp) > 0 {
				err = fmt.Errorf("error string %q in xml path subscript.", string(tmp))
				return
			}
		case '=':
			// 如果当前处于属性值区，则等号属于属性值的一部分
			if inAttrValue {
				tmp = append(tmp, ch)
				continue L
			}
			// 不在下标区域中出现等号，这是错误的
			if !inSubscript {
				err = fmt.Errorf("char '%c' is not allow to be in xml path.", ch)
				return
			}

			// 下标区段内出现等号，则认为属性名称结束，进入属性值区段
			inAttrValue = true
			tag.attrName = string(bytes.TrimSpace(tmp))
			tmp = []byte{}
			// 此时，属性名称不能为空
			if tag.attrName == "" {
				err = fmt.Errorf("attribute name in subscript of xml path is not allow to be an empty string.")
				return
			}
		case '/':
			// 如果当前在属性值区段内，则认为斜杠为属性值的一部分
			if inAttrValue {
				tmp = append(tmp, ch)
				continue L
			}
			// 在下标区段，但又不在属性值区段出现斜杠，这是错误的
			if inSubscript {
				err = fmt.Errorf("char '%c' is not allow to be in xml path subscript.", ch)
				return
			}
			// 第一个 tag 结束
			// 注意：这里必须先判断 tag.tag 是否为空，如果为空才对其赋值
			// 因为有可能在下标开始之前，已经对 tag.tag 赋过值，这里再赋值的话，会冲掉前面的
			if tag.tag == "" {
				tag.tag = string(bytes.TrimSpace(tmp))
				tmp = []byte{}
			} else if len(bytes.TrimSpace(tmp)) > 0 {
				// 如果跑到这里来，说明 tag.tag 已经在下标起始之前已经赋过值
				// 如果这里 tmp 不为空，则可能出现类似这样的路径：aa/bb[2]xx/cc
				// 这是错误的
				if tag.attrName != "" {
					err = fmt.Errorf(`string %q after xml path sub tag "%s[%s=%s] is invalid.`,
						string(tmp), tag.tag, tag.attrName, tag.attrValue)
				} else {
					err = fmt.Errorf(`string %q after xml path sub tag "%s[%s=%d] is invalid"`,
						string(tmp), tag.tag, tag.attrName, tag.index)
				}
				return
			}
			break L
		default:
			tmp = append(tmp, ch)
		}
	}

	// 未离开下标区就结束，这是错误的
	if attrValueChar != ' ' || inSubscript {
		err = fmt.Errorf("error xml path, subscript is not complete.")
		return
	}
	if len(tmp) > 0 {
		tag.tag = string(bytes.TrimSpace(tmp))
	}
	if idx < len(path) {
		tail = path[idx+1:]
	}
	return
}

// -------------------------------------------------------------------
// Node
// -------------------------------------------------------------------
type S_Node struct {
	s_NameText

	docPtr    *S_Doc
	parentPtr *S_Node
	attrPtrs  []*S_Attr
	childPtrs []*S_Node
	isCData   bool
}

// 创建一个 xml 节点
// 注意：如果 name 是不合法的 xml tag 名称，则返回 nil
func CreateNode(doc *S_Doc, name, text string) *S_Node {
	if !isValidName(name) {
		return nil
	}
	return &S_Node{
		s_NameText: s_NameText{name, text},
		docPtr:     doc,
		parentPtr:  nil,
		attrPtrs:   []*S_Attr{},
		childPtrs:  []*S_Node{},
		isCData:    false,
	}
}

// ---------------------------------------------------------
// private
// ---------------------------------------------------------
// 获取指定的名字空间
func (this *S_Node) getNameSpace(nsName string) string {
	if nsName == "xml" {
		return xmlNameSpace
	}

	var space string
	count := len(this.childPtrs)
	for i := count - 1; i >= 0; i -= 1 {
		space = this.childPtrs[i].getNameSpace(nsName)
		if space != "" {
			return space
		}
	}
	if space == "" {
		for _, attr := range this.attrPtrs {
			if attr.Space() != "xmlns" {
				continue
			}
			if attr.Local() == nsName {
				return attr.Text()
			}
		}
	}
	return ""
}

// ---------------------------------------------------------
// public
// ---------------------------------------------------------
// 返回节点所属文档
func (this *S_Node) Doc() *S_Doc {
	return this.docPtr
}

// 返回节点所在文档的根节点
func (this *S_Node) Root() *S_Node {
	return this.docPtr.rootPtr
}

// 返回节点的父节点
func (this *S_Node) Parent() *S_Node {
	return this.parentPtr
}

// 获取指点指定名称的属性
func (this *S_Node) Attr(attrName string) *S_Attr {
	for _, attr := range this.attrPtrs {
		if attr.name == attrName {
			return attr
		}
	}
	return nil
}

// 获取节点的 ID 属性（没有的话，则返回空字符串）
// 提示：相当于 Attr("id").Text()
func (this *S_Node) ID() string {
	return this.Attr("id").text
}

// 返回节点的属性列表
func (this *S_Node) Attrs() []*S_Attr {
	return this.attrPtrs
}

// 返回属性个数
func (this *S_Node) AttrCount() int {
	return len(this.attrPtrs)
}

// 获取指定路径下的子孙节点，如：
// <root> <aa>
//	<bb k='v1'> xx </bb>
//	<bb k='v2'> yy </bb>
// </aa> </root>
// root.Child("aa/bb") 返回第一个 bb 节点
// root.Child("aa/bb[0]") 返回第一个 bb 节点
// root.Child("aa/bb[-1]") 返回第二个 bb 节点
// root.Child("aa/bb[k=v2]") 返回第二个 bb 节点
// 注意：node.Child("") 返回 node 本身
func (this *S_Node) GetChild(path string) (*S_Node, error) {
	tag, tail, err := _firstSearchPath(path)
	if err != nil {
		return nil, fserror.Wrapf(err, "get xml child node fail: %q.", path)
	}
	if tag.tag == "" {
		return this, nil
	}

	var child *S_Node
	if tag.attrName != "" {
		for _, ch := range this.childPtrs {
			if ch.Name() != tag.tag {
				continue
			}
			if ch.Attr(tag.attrName).Text() == tag.attrValue {
				child = ch
				break
			}
		}
	} else if tag.index >= 0 {
		idx := 0
		for _, ch := range this.childPtrs {
			if ch.Name() != tag.tag {
				continue
			}
			if idx == tag.index {
				child = ch
				break
			}
			idx += 1
		}
	} else {
		idx := -1
		for i := this.ChildCount() - 1; i >= 0; i -= 1 {
			ch := this.childPtrs[i]
			if ch.Name() != tag.tag {
				continue
			}
			if idx == tag.index {
				child = ch
				break
			}
			idx -= 1
		}
	}
	if child == nil {
		return nil, fmt.Errorf("xml child node is not exist: %q.", path)
	}
	if len(tail) > 0 {
		node, e := child.GetChild(tail)
		if e != nil {
			return node, fserror.Wrapf(e, "get xml child node fail: %q.", path)
		}
		return node, nil
	}
	return child, nil
}

func (this *S_Node) Child(path string) *S_Node {
	node, _ := this.GetChild(path)
	return node
}

// 获取指定索引的子节点
func (this *S_Node) ChildByIndex(index int) *S_Node {
	if index < 0 {
		return nil
	}
	if index > len(this.childPtrs) {
		return nil
	}
	return this.childPtrs[index]
}

// 返回所有子节点
func (this *S_Node) Children() []*S_Node {
	return this.childPtrs
}

// 获取子节点个数
func (this *S_Node) ChildCount() int {
	return len(this.childPtrs)
}

// 是否是 cdata 数据
func (this *S_Node) IsCData() bool {
	return this.isCData
}

// 将节点内容设置为 CDATA
func (this *S_Node) SetIsCData(b bool) {
	this.isCData = true
}

// ---------------------------------------------------------
// 如果属性已经存在，则修改属性值；如果属性不存在，则增加属性
func (this *S_Node) SetAttr(attr *S_Attr) bool {
	if attr == nil {
		return false
	}
	for _, a := range this.attrPtrs {
		if a.name == attr.name {
			a.text = attr.text
			return true
		}
	}
	this.attrPtrs = append(this.attrPtrs, attr)
	return true
}

// 删除属性
// 如果要删除的属性存在，则返回 true
func (this *S_Node) RemoveAttr(name string) bool {
	index := -1
	var attr *S_Attr
	for index, attr = range this.attrPtrs {
		if attr.name == name {
			break
		}
	}
	if index > 0 {
		if index == len(this.attrPtrs)-1 {
			this.attrPtrs = this.attrPtrs[:index]
		} else {
			this.attrPtrs = append(this.attrPtrs[:index], this.attrPtrs[index+1:]...)
		}
		return true
	}
	return false
}

// 清除所有属性
func (this *S_Node) ClearAttrs() {
	this.attrPtrs = []*S_Attr{}
}

// ---------------------------------------------------------
// 添加一个子节点
func (this *S_Node) AddChild(child *S_Node) bool {
	if child == nil {
		return false
	}
	child.docPtr = this.docPtr
	child.parentPtr = this
	this.childPtrs = append(this.childPtrs, child)
	return true
}

// 添加一组子节点
func (this *S_Node) AddChildren(children []*S_Node) {
	for _, child := range children {
		this.AddChild(child)
	}
}

// 删除指定名称的子节点，如果指定子节点存在，则返回 true，如果要移除的子节点不存在，则返回 false
func (this *S_Node) RemoveChild(name string) bool {
	removed := false
	chs := make([]*S_Node, 0)
	for _, ch := range this.childPtrs {
		if ch.name != name {
			chs = append(chs, ch)
		} else {
			removed = true
		}
	}
	this.childPtrs = chs
	return removed
}

// 删除指定条件的子节点
// 参数 f 返回 true，则删除
func (this *S_Node) RemoveChildOf(f func(*S_Node) bool) bool {
	removed := false
	chs := make([]*S_Node, 0)
	for _, ch := range this.childPtrs {
		if f(ch) {
			removed = true
		} else {
			chs = append(chs, ch)
		}
	}
	this.childPtrs = chs
	return removed

}

// 清空所有子节点
func (this *S_Node) ClearChildren() {
	this.childPtrs = []*S_Node{}
}

// ---------------------------------------------------------
// 获取第一个子节点
func (this *S_Node) FirstChild() *S_Node {
	if len(this.childPtrs) > 0 {
		return this.childPtrs[0]
	}
	return nil
}

// 获取最后一个子节点
func (this *S_Node) LastChild() *S_Node {
	if len(this.childPtrs) > 0 {
		return this.childPtrs[len(this.childPtrs)-1]
	}
	return nil
}

// 获取所有名称为 name 的所有子节点
func (this *S_Node) ChildrenOfName(name string) []*S_Node {
	chs := make([]*S_Node, 0)
	for _, n := range this.childPtrs {
		if n.name == name {
			chs = append(chs, n)
		}
	}
	return chs
}

// 获取符合条件的所有子节点
func (this *S_Node) ChildrenOfFunc(f func(*S_Node) bool) []*S_Node {
	chs := make([]*S_Node, 0)
	for _, n := range this.childPtrs {
		if f(n) {
			chs = append(chs, n)
		}
	}
	return chs
}

// 搜索指定名称的第一个子孙节点
func (this *S_Node) FindNode(name string) *S_Node {
	if this.name == name {
		return this
	}

	for _, ch := range this.childPtrs {
		node := ch.FindNode(name)
		if node != nil {
			return node
		}
	}
	return nil
}

// 搜索具有属性 ID，并且 ID 值为指定值的第一个节点
func (this *S_Node) FindChildByID(id string) *S_Node {
	if this.ID() == id {
		return this
	}

	for _, ch := range this.childPtrs {
		node := ch.FindChildByID(id)
		if node != nil {
			return node
		}
	}
	return nil
}

// ---------------------------------------------------------
// 获取前一个兄弟节点
func (this *S_Node) PreSibling() *S_Node {
	if this.parentPtr == nil {
		return nil
	}
	var node *S_Node = nil
	for _, n := range this.parentPtr.childPtrs {
		if n == this {
			return node
		} else {
			node = n
		}
	}
	return nil
}

// 获取后一个兄弟节点
func (this *S_Node) NextSibling() *S_Node {
	if this.parentPtr == nil {
		return nil
	}
	var node *S_Node = nil
	maxIndex := len(this.parentPtr.childPtrs) - 1
	for index := maxIndex; index >= 0; index -= 1 {
		if this.parentPtr.childPtrs[index] == this {
			return node
		} else {
			node = this.parentPtr.childPtrs[index]
		}
	}
	return nil
}

// ---------------------------------------------------------
// 遍历所有节点
// 如果 f 返回 false，则结束遍历
func (this *S_Node) Travel(f func(*S_Node) bool) bool {
	if !f(this) {
		return false
	}
	for _, node := range this.childPtrs {
		if !node.Travel(f) {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------
// 克隆节点（深拷贝）
func (this *S_Node) Clone(doc *S_Doc) *S_Node {
	node := &S_Node{
		s_NameText: s_NameText{this.name, this.text},
		docPtr:     doc,
		parentPtr:  nil,
		attrPtrs:   []*S_Attr{},
		childPtrs:  []*S_Node{},
		isCData:    this.isCData,
	}
	for _, attr := range this.attrPtrs {
		node.attrPtrs = append(node.attrPtrs, attr.clone())
	}
	for _, child := range this.childPtrs {
		node.AddChild(child.Clone(doc))
	}
	return node
}
