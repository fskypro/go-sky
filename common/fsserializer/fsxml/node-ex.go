/**
@copyright: fantasysky 2016
@brief: extension of node
@author: fanky
@version: 1.0
@date: 2020-02-01
**/

package fsxml

// 获取指定名称一组子节点的 text，如：
// <root> <item> aa </item> <item> bb </item>  </root>
// 调用 ChildTextList("item") 返回 []string{"aa", "bb"}
func (this *S_Node) ReadItem(name string) []string {
	items := make([]string, 0)
	nodes := this.ChildrenOfName(name)
	for _, node := range nodes {
		items = append(items, node.text)
	}
	return items
}

// 添加一组同名称的子节点
func (this *S_Node) WriteItem(name string, texts []string) {
	for _, text := range texts {
		node := CreateNode(this.docPtr, name, text)
		this.AddChild(node)
	}
}

// 获取指定名称一组子节点的 text，如：
// <root> <item> aa </item> <item> bb </item>  </root>
// 调用 ChildBytesList("item") 返回 [][]byte{[]byte("aa"), []byte("bb")}
func (this *S_Node) ReadItemBytes(name string) [][]byte {
	items := make([][]byte, 0)
	nodes := this.ChildrenOfName(name)
	for _, node := range nodes {
		items = append(items, []byte(node.text))
	}
	return items
}

// 添加一组同名的子节点
func (this *S_Node) WriteItemBytes(name string, btexts [][]byte) {
	for _, btext := range btexts {
		node := CreateNode(this.docPtr, name, string(btext))
		this.AddChild(node)
	}
}
