/**
@copyright: fantasysky 2016
@brief: xml node
@author: fanky
@version: 1.0
@date: 2020-01-31
**/

package fsxml

import (
	"strconv"
	"strings"
)

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
//	<bb k=v1> xx </bb>
//	<bb k=v2> yy </bb>
// </aa> </root>
// root.Child("aa/bb") 返回第一个 bb 节点
// root.Child("aa/bb[0]") 返回第一个 bb 节点
// root.Child("aa/bb[-1]") 返回第二个 bb 节点
// root.Child("aa/bb[k=v2]") 返回第二个 bb 节点
// 注意：node.Child("") 返回 node 本身
func (this *S_Node) Child(path string) *S_Node {
	if path == "" {
		return this
	}

	// 解释字键和索引或属性键值
	getTagIndex := func(tag string) (key string, index int, ak, av string) {
		subscript := ""
		hasSubscript := false
		for _, c := range tag {
			if c == '[' {
				hasSubscript = true
			} else if c == ']' {
				break
			} else if hasSubscript {
				subscript += string([]rune{c})
			} else {
				key += string([]rune{c})
			}
		}
		akv := strings.Split(subscript, "=")
		if len(akv) == 2 {
			ak, av = strings.TrimSpace(akv[0]), strings.TrimSpace(akv[1])
		} else if len(subscript) > 0 {
			index, _ = strconv.Atoi(subscript)
		}
		return
	}

	tags := strings.Split(path, "/")
	var node *S_Node = this
	var child *S_Node
	for _, tag := range tags {
		child = nil
		tag, idx, ak, av := getTagIndex(tag)
		calc := 0

		if len(ak) > 0 { // 查找与指定属性值匹配的节点
			for _, child = range node.childPtrs {
				if child.name != tag {
					continue
				}
				attr := child.Attr(ak)
				if attr == nil {
					continue
				}
				if attr.Text() == av {
					node = child
					break
				}
			}
		} else if idx >= 0 { // 顺序索引
			for _, child = range node.childPtrs {
				if child.name != tag {
					continue
				}
				if calc == idx {
					node = child
					break
				} else {
					calc += 1
				}
			}
		} else { // 逆序索引
			calc = -1
			for i := len(node.childPtrs) - 1; i >= 0; i -= 1 {
				child = node.childPtrs[i]
				if child.name != tag {
					continue
				}
				if calc == idx {
					node = child
					break
				} else {
					calc -= 1
				}
			}
		}
		if node != child {
			return nil
		}
	}
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
