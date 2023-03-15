/**
@copyright: fantasysky 2016
@brief: dom explain
@author: fanky
@version: 1.0
@date: 2020-02-02
**/

package fsxml

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"fsky.pro/fsserializer/fsxml/internal/xml"
)

// -------------------------------------------------------------------
// inners
// -------------------------------------------------------------------
func _toName(xname *xml.Name) string {
	if xname.SpaceName == "" {
		return xname.Local
	}
	return fmt.Sprintf("%s:%s", xname.SpaceName, xname.Local)
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func LoadString(xml string) (doc *S_Doc, err error) {
	return LoadReader(strings.NewReader(xml))
}

func LoadBytes(bxml []byte) (doc *S_Doc, err error) {
	return LoadReader(bytes.NewReader(bxml))
}

func LoadFile(path string) (*S_Doc, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadReader(file)
}

func LoadReader(reader io.Reader) (*S_Doc, error) {
	decoder := xml.NewDecoder(reader)
	doc := newEmptyDoc()
	var node *S_Node
	t, err := decoder.Token()
	err = parseError(err)
	if err != nil {
		return nil, err
	}

	for t != nil {
		switch token := t.(type) {
		// 起始标签
		case xml.StartElement:
			newNode := CreateNode(doc, _toName(&token.Name), "")
			for _, attr := range token.Attr {
				newNode.SetAttr(NewAttr(_toName(&attr.Name), attr.Value))
			}

			if node != nil {
				node.AddChild(newNode)
			}
			if doc.rootPtr == nil {
				doc.rootPtr = newNode
			}
			node = newNode

		// 标签数据
		case xml.CharData:
			// 这个地方有个问题：无法区分究竟是不是 CDATA 数据：<![CDATA[...]]>
			if node != nil {
				node.text = string(bytes.TrimSpace(token))
			}

		// 结束标签
		case xml.EndElement:
			node = node.parentPtr

		// 处理指令：<?target instruction?>
		case xml.ProcInst:
			inst := t.(xml.ProcInst)
			text := fmt.Sprintf("<?%s %s?>", inst.Target, inst.Inst)
			if inst.Target == "xml" {
				doc.Header = text
			} else {
				doc.ProcInsts = append(doc.ProcInsts, text)
			}

		// 指示：<!directive>
		case xml.Directive:
			doc.Directives = append(doc.Directives, fmt.Sprintf("<!%s>", t.(xml.Directive)))
		}

		t, err = decoder.Token()
	}

	err = parseError(err)
	if err != io.EOF {
		return nil, err
	}
	if doc.rootPtr == nil {
		return nil, fmt.Errorf("no xml root tag in document.")
	}
	return doc, nil
}
