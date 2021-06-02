/**
@copyright: fantasysky 2016
@brief: definations
@author: fanky
@version: 1.0
@date: 2020-02-18
**/

package xmlex

const (
	// xml 命名空间
	xmlNameSpace = "http://www.w3.org/XML/1998/namespace"
)

// 格式化 xml 时使用的换行符
type T_Endline string

const (
	CR   T_Endline = "\r"
	LF             = "\n"
	CRLF           = "\r\n"
)
