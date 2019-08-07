/**
@copyright: fantasysky 2016
@brief: 错误类型
@author: fanky
@version: 1.0
@date: 2019-05-31
**/

package freejson

import "fmt"

// 解释错误
type ParseError struct {
	msg string
}

func newParseError(site int, rowIdx, colIdx int) *ParseError {
	return &ParseError{
		msg: fmt.Sprintf("error json format. site=%d, row=%d, col=%d", site, rowIdx+1, colIdx+1),
	}
}

func (this *ParseError) Error() string {
	return this.msg
}

// 非正确 json 文件
type JsonFileError struct {
	msg string
}

func newJsonFileError(path string, errParse error) *JsonFileError {
	return &JsonFileError{
		msg: fmt.Sprintf("invalid json file(%q). ", path) + errParse.Error(),
	}
}

func (this *JsonFileError) Error() string {
	return this.msg
}
