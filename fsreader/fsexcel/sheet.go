/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: excel sheet
@author: fanky
@version: 1.0
@date: 2023-09-13
**/

package fsexcel

import (
	"regexp"
)

// -------------------------------------------------------------------
// 分页
// -------------------------------------------------------------------
type S_Sheet struct {
	owner *S_Excel
	name  string
	index int
}

func (this *S_Sheet) Name() string {
	return this.name
}

func (this *S_Sheet) Index() int {
	return this.index
}

// 返回表头在指定行的表格
// headerRow: 表头所在行(注意：第一行为 1)
// validRow ：有效起始行(注意：第一行为 1)
// 如果表头不存在，则返回 nul
func (this *S_Sheet) AsRowsTable(headerRow, validRow int) *S_RowsTable {
	data, err := this.owner.GetRows(this.name)
	if err != nil { return nil }
	return newRowsTable(data, headerRow-1, validRow-1)
}

// 返回表头在指定列的表格
// headerCol: 表头所在行(注意：第一行为 "A")
// validCol ：有效起始行(注意：第一行为 "A")
// 如果表头不存在，则返回 nul
func (this *S_Sheet) AsColsTable(headerCol, validCol string) *S_ColsTable {
	data, err := this.owner.GetCols(this.name)
	if err != nil { return nil }
	re := regexp.MustCompile("^[A-Z]+$")
	if !re.MatchString(headerCol) || !re.MatchString(validCol) {
		return nil
	}
	headerIndex := ColOrderToIndex(headerCol)
	validIndex := ColOrderToIndex(validCol)
	return newColsTable(data, headerIndex, validIndex)
}
