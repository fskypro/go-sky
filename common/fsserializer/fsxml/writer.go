/**
@copyright: fantasysky 2016
@brief: xml writer
@author: fanky
@version: 1.0
@date: 2020-02-09
**/

package fsxml

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

import (
	"fsky.pro/fsserializer/fsxml/internal/xml"
	"fsky.pro/fsstr/convert"
)

func _writeAttribute(attr *S_Attr, writer *bufio.Writer) error {
	_, err := writer.WriteString(fmt.Sprintf(` %s="`, attr.name))
	if err != nil {
		return nil
	}
	err = xml.EscapeText(writer, attr.TextBytes())
	if err != nil {
		return err
	}
	_, err = writer.WriteString(`"`)
	return err
}

func _writeNode(node *S_Node, writer *bufio.Writer) error {
	_, err := writer.WriteString(fmt.Sprintf("<%s", node.name))
	if err != nil {
		return err
	}
	for _, attr := range node.attrPtrs {
		if err := _writeAttribute(attr, writer); err != nil {
			return err
		}
	}
	err = writer.WriteByte('>')
	if err != nil {
		return err
	}

	if node.text != "" {
		if node.isCData {
			_, err = writer.WriteString(fmt.Sprintf("<![CDATA[%s]]>", node.text))
			if err != nil {
				return err
			}
		} else {
			err = xml.EscapeText(writer, convert.String2Bytes(node.text))
			if err != nil {
				return err
			}
		}
	}

	for _, ch := range node.childPtrs {
		err = _writeNode(ch, writer)
		if err != nil {
			return err
		}
	}

	_, err = writer.WriteString(fmt.Sprintf("</%s>", node.name))
	return err
}

func writeDoc(doc *S_Doc, w io.Writer) (err error) {
	writer := bufio.NewWriter(w)

	defer func() {
		if err == nil {
			err = writer.Flush()
		}
	}()

	if doc.Header != "" {
		_, err = writer.WriteString(doc.Header)
		if err != nil {
			return
		}
	}
	for _, inst := range doc.ProcInsts {
		_, err = writer.WriteString(inst)
		if err != nil {
			return
		}
	}
	for _, direct := range doc.Directives {
		_, err = writer.WriteString(direct)
		if err != nil {
			return
		}
	}
	if doc.rootPtr != nil {
		err = _writeNode(doc.rootPtr, writer)
	}
	return
}

// -------------------------------------------------------------------
func _writeNodeIndent(node *S_Node, writer *bufio.Writer, layer int, indent string, endline T_Endline) error {
	_, err := writer.WriteString(fmt.Sprintf("%s<%s", strings.Repeat(indent, layer), node.name))
	if err != nil {
		return err
	}
	for _, attr := range node.attrPtrs {
		if err := _writeAttribute(attr, writer); err != nil {
			return err
		}
	}
	err = writer.WriteByte('>')
	if err != nil {
		return err
	}

	if len(node.childPtrs) == 0 {
		if node.text == "" {
			_, err = writer.WriteString(fmt.Sprintf("</%s>%s", node.name, endline))
		} else if node.isCData {
			space := strings.Repeat(indent, layer)
			_, err = writer.WriteString(fmt.Sprintf("%s%s<![CDATA[%s]]>%s%s</%s>%s", endline, space+indent, node.text, endline, space, node.name, endline))
			if err != nil {
				return err
			}
		} else {
			err = writer.WriteByte(' ')
			if err != nil {
				return err
			}
			err = xml.EscapeText(writer, convert.String2Bytes(node.text))
			if err != nil {
				return err
			}
			_, err = writer.WriteString(fmt.Sprintf(" </%s>%s", node.name, endline))
			if err != nil {
				return err
			}
		}
	} else {
		if node.text != "" {
			if node.isCData {
				_, err = writer.WriteString(fmt.Sprintf("%s%s<![CDATA[%s]]>", endline, strings.Repeat(indent, layer+1), node.text))
				if err != nil {
					return err
				}
			} else {
				err = writer.WriteByte(' ')
				if err != nil {
					return err
				}
				err = xml.EscapeText(writer, convert.String2Bytes(node.text))
				if err != nil {
					return err
				}
			}
		}
		_, err = writer.WriteString(string(endline))
		if err != nil {
			return err
		}

		for _, ch := range node.childPtrs {
			err = _writeNodeIndent(ch, writer, layer+1, indent, endline)
			if err != nil {
				return err
			}
		}
		_, err = writer.WriteString(fmt.Sprintf("%s</%s>%s", strings.Repeat(indent, layer), node.name, endline))
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDocIndent(doc *S_Doc, w io.Writer, indent string, endline T_Endline) (err error) {
	writer := bufio.NewWriter(w)

	defer func() {
		if err == nil {
			err = writer.Flush()
		}
	}()

	if doc.Header != "" {
		_, err = writer.WriteString(doc.Header + string(endline))
		if err != nil {
			return
		}
	}
	for _, inst := range doc.ProcInsts {
		_, err = writer.WriteString(inst + string(endline))
		if err != nil {
			return
		}
	}
	for _, direct := range doc.Directives {
		_, err = writer.WriteString(direct + string(endline))
		if err != nil {
			return
		}
	}
	if doc.rootPtr != nil {
		err = _writeNodeIndent(doc.rootPtr, writer, 0, indent, endline)
	}
	return
}
