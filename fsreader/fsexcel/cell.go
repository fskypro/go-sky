/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: cell
@author: fanky
@version: 1.0
@date: 2023-09-13
**/

package fsexcel

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
