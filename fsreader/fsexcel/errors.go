/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: errors
@author: fanky
@version: 1.0
@date: 2025-07-16
**/

package fsexcel

import (
	"fmt"

	"fsky.pro/fsconv"
)

// -------------------------------------------------------------------
// 空格子错误
// -------------------------------------------------------------------
type EmptyCellError struct {
	RowOrder int
	Header   string
}

func (self EmptyCellError) Error() string {
	return fmt.Sprintf("cell(row=%d, header=%q) is empty", self.RowOrder, self.Header)
}

// -------------------------------------------------------------------
// 数据转换错误
// -------------------------------------------------------------------
// 表格内容为数值时，超出指定类型范围
type RangeError = fsconv.RangeError

// 将表格内容转换为指定类型时，转换失败
type SyntaxError = fsconv.SyntaxError
