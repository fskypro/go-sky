/**
@copyright: fantasysky 2016
@brief: errors
@author: fanky
@version: 1.0
@date: 2020-02-01
**/

package xmlex

import (
	"fmt"
	"strconv"
)

import (
	"fsky.pro/fsserializer/xmlex/internal/xml"
)

// -------------------------------------------------------------------
// SyntaxyError
// -------------------------------------------------------------------
type SyntaxError struct {
	err *xml.SyntaxError
}

func (this *SyntaxError) Error() string {
	return this.err.Error()
}

func parseError(err error) error {
	switch err.(type) {
	case *xml.SyntaxError:
		return &SyntaxError{err.(*xml.SyntaxError)}
	}
	return err
}

// -------------------------------------------------------------------
// RangeError/TypeError
// -------------------------------------------------------------------
// 数值内容超出数值表示范围
type RangeError struct {
	Text     string
	TypeName string
}

func (this *RangeError) Error() string {
	return fmt.Sprintf("value %q is out of range, can't convert to type %s", this.Text, this.TypeName)
}

// --------------------------------------------------------
// 非数值不能转换
type TypeError struct {
	Text     string
	TypeName string
}

func (this *TypeError) Error() string {
	return fmt.Sprintf("can't convert value %q to %s type value.", this.Text, this.TypeName)
}

func convertError(err error, text string, typeName string) error {
	if err == strconv.ErrRange {
		return &RangeError{text, typeName}
	} else if err == strconv.ErrSyntax {
		return &TypeError{text, typeName}
	}
	return err
}
