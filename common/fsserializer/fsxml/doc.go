/**
@copyright: fantasysky 2016
@brief: xml document
@author: fanky
@version: 1.0
@date: 2020-01-31
**/

package fsxml

import (
	"bytes"
	"io"
	"os"
)
import (
	"fsky.pro/fsstr/convert"
)

const (
	Header = `<?xml version="1.0" encoding="UTF-8"?>`
)

type S_Doc struct {
	Header     string   // 类似于：<?xml version="1.0" encoding="ISO-8859-1"?>
	ProcInsts  []string // 类似于：<?target instruction>
	Directives []string // 类似于：<!directive>

	rootPtr *S_Node // 根节点
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func newEmptyDoc() *S_Doc {
	return &S_Doc{
		Header:     "",
		ProcInsts:  []string{},
		Directives: []string{},
		rootPtr:    nil,
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 新建一个带空根节点的文档
// 注意：如果 rootName 是不合法的 xml tag 名称，则返回 nil
func NewDoc(rootName string) *S_Doc {
	doc := &S_Doc{
		Header:     Header,
		ProcInsts:  []string{},
		Directives: []string{},
	}
	doc.rootPtr = CreateNode(doc, rootName, "")
	if doc.rootPtr == nil {
		return nil
	}
	return doc
}

// ---------------------------------------------------------
// 返回根节点
func (this *S_Doc) Root() *S_Node {
	return this.rootPtr
}

// 获取指定别名的名字空间
func (this *S_Doc) GetNamespace(nsName string) string {
	if this.rootPtr == nil {
		return ""
	}
	return this.rootPtr.getNameSpace(nsName)
}

// ---------------------------------------------------------
// 生成字符串形式的 xml 文档
func (this *S_Doc) ToXML() (string, error) {
	bxml, err := this.ToXMLData()
	if err != nil {
		return "", err
	}
	return convert.Bytes2String(bxml), nil
}

// 生成字节组形式的 xml 文档
func (this *S_Doc) ToXMLData() ([]byte, error) {
	var bxml = []byte{}
	w := bytes.NewBuffer(bxml)
	err := writeDoc(this, w)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), err
}

// 生成字符串形式的带嵌套格式化的 xml 文档
func (this *S_Doc) ToXMLIndent(indent string, endline T_Endline) (string, error) {
	bxml, err := this.ToXMLDataIndent(indent, endline)
	if err != nil {
		return "", nil
	}
	return convert.Bytes2String(bxml), nil
}

// 生成字节组形式的带嵌套格式化的 xml 文档
func (this *S_Doc) ToXMLDataIndent(indent string, endline T_Endline) ([]byte, error) {
	var bxml = []byte{}
	w := bytes.NewBuffer(bxml)
	err := writeDocIndent(this, w, indent, endline)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), err

}

// 将 xml 文档写入到流
func (this *S_Doc) Write(w io.Writer) error {
	return writeDoc(this, w)
}

// 将带嵌套格式化的 xml 文档写入到流
func (this *S_Doc) WriteIndent(w io.Writer, indent string, endline T_Endline) error {
	return writeDocIndent(this, w, indent, endline)
}

// 将 xml 数据保存成 xml 格式文件
// 注意：如果文件已经存在，则会覆盖；创建的文件权限为 0666
func (this *S_Doc) Save(path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = f.Close()
		} else {
			f.Close()
		}
	}()

	data, err := this.ToXMLData()
	if err == nil {
		_, err = f.Write(data)
	}
	return
}

// 将 xml 数据保存成带嵌套格式化的 xml 格式文件
// 注意：如果文件已经存在，则会覆盖；创建的文件权限为 0666
func (this *S_Doc) SaveIndent(path string, indent string, endline T_Endline) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = f.Close()
		} else {
			f.Close()
		}
	}()

	data, err := this.ToXMLDataIndent(indent, endline)
	if err == nil {
		_, err = f.Write(data)
	}
	return
}

// ---------------------------------------------------------
// 克隆文档（深拷贝）
func (this *S_Doc) Clone() *S_Doc {
	doc := &S_Doc{
		Header:     this.Header,
		ProcInsts:  this.ProcInsts[:],
		Directives: this.Directives[:],
	}
	doc.rootPtr = this.rootPtr.Clone(doc)
	return doc
}
