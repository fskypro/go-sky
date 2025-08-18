/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: cell
@author: fanky
@version: 1.0
@date: 2023-09-13
**/

package fsexcel

import (
	"strings"

	"fsky.pro/fsconv"
)

// -------------------------------------------------------------------
// 单元格
// -------------------------------------------------------------------
type S_Cell struct {
	string
	Header   string // 对应的表头名称
	rowIndex int    // 行索引
	colIndex int    // 列索引
}

func newCell(headerText string, rowIndex, colIndex int, text string) *S_Cell {
	return &S_Cell{text, headerText, rowIndex, colIndex}
}

// 在 excel 表格中的行号
// 注：第一行为 1
func (this *S_Cell) RowOrder() int {
	return this.rowIndex + 1
}

// 在 excel 表格中的列序号
// 注：第一列为 A，第 26 列为 AA
func (this *S_Cell) ColOrder() string {
	return IndexToColOrder(this.colIndex)
}

// 表格中的文本内容
func (this *S_Cell) Text() string {
	return this.string
}

// 去除两边空白字符后的文本内容
func (this *S_Cell) TrimText() string {
	return strings.TrimSpace(this.string)
}

// ---------------------------------------------------------
// 将指定行，指定列名下的格子内容，转换为泛型类型的值
// 如果格子内容为空，则返回参数中的默认值
// 参数：
// -> row  ：Excel 行
// -> def  ：当表格内容为空时，使用的默认值
// 可能返回错误：
// -> RangeError   ：转换为数值超出范围，如果不转换为数值的话，不会出现该错误
// -> SyntaxError  ：无法转换为数值，如果不转换为数值的话，不会出现该错误
func CellAsDefault[T fsconv.T_Base](cell *S_Cell, def T) (T, error) {
	text := cell.TrimText()
	if text == "" {
		return def, nil
	}
	return fsconv.StrTo[T](text)
}

// 将指定行，指定列名下的格子内容，转换为泛型类型的值
// 如果格子内容为空，则返回 EmptyCellError 错误
// 参数：
// -> row  ：Excel 行
// -> def  ：当表格内容为空时，使用的默认值
// 可能返回错误：
// -> EmptyCellError：格子内容为空
// -> RangeError    ：转换为数值超出范围，如果不转换为数值的话，不会出现该错误
// -> SyntaxError   ：无法转换为数值，如果不转换为数值的话，不会出现该错误
func CellAs[T fsconv.T_Base](cell *S_Cell) (T, error) {
	var v T
	text := cell.TrimText()
	if text == "" {
		return v, EmptyCellError{cell.rowIndex + 1, cell.Header}
	}
	return fsconv.StrTo[T](text)
}
