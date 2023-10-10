/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: rows or cols
@author: fanky
@version: 1.0
@date: 2023-09-13
**/

package fsexcel

// -------------------------------------------------------------------
// 一条记录
// -------------------------------------------------------------------
type s_Item[T S_RowHeader|S_ColHeader] struct {
	table *s_Table[T]
	index int
}

func newItem[T S_RowHeader|S_ColHeader](table *s_Table[T], index int) *s_Item[T] {
	return &s_Item[T]{
		table: table,
		index: index,
	}
}

// ---------------------------------------------------------
// 记录迭代器
type s_ItemIter[T S_RowHeader|S_ColHeader] struct {
	table *s_Table[T]
	index int
}

func newItemsIter[T S_RowHeader|S_ColHeader](table *s_Table[T]) *s_ItemIter[T] {
	return &s_ItemIter[T]{
		table: table,
		index: table.validIndex - 1,
	}
}

func (this *s_ItemIter[T]) Next() bool {
	this.index++
	return this.index < len(this.table.data)
}

// -------------------------------------------------------------------
// 一行记录
// -------------------------------------------------------------------
type S_Row s_Item[S_RowHeader]

func newRow(table *s_Table[S_RowHeader], index int) *S_Row {
	return (*S_Row)(newItem(table, index))
}

func (this *S_Row) Order() int {
	return this.index + 1
}

// 通过表头名称获取单元格
func (this *S_Row) GetCell(name string) *S_Cell {
	header := this.table.nameHeaders[name]
	if header == nil { return nil }
	rowData := this.table.data[this.index]
	text := ""
	if header.index < len(rowData) {
		text = this.table.data[this.index][header.index]
	}
	return newCell(header.name, this.index, header.index, text)
}

// 获取所有单元格
func (this *S_Row) GetCells() []*S_Cell {
	cells := []*S_Cell{}
	for _, h := range this.table.headers {
		rowData := this.table.data[this.index]
		text := ""
		if h.index < len(rowData) {
			text = this.table.data[this.index][h.index]
		}
		cell := newCell(h.name, this.index, h.index, text)
		cells = append(cells, cell)
	}
	return cells
}

// 遍历所有单元格
func (this *S_Row) For(f func(*S_Cell)) {
	for _, h := range this.table.headers {
		rowData := this.table.data[this.index]
		text := ""
		if h.index < len(rowData) {
			text = this.table.data[this.index][h.index]
		}
		cell := newCell(h.name, this.index, h.index, text)
		f(cell)
	}
}

// ---------------------------------------------------------
type S_RowsIter struct{ *s_ItemIter[S_RowHeader] }

func newRowsIter(table *s_Table[S_RowHeader]) *S_RowsIter {
	return &S_RowsIter{newItemsIter(table)}
}

func (this *S_RowsIter) Row() *S_Row {
	if this.index < len(this.table.data) {
		return newRow(this.table, this.index)
	}
	return nil

}

// -------------------------------------------------------------------
// 一列记录
// -------------------------------------------------------------------
type S_Col s_Item[S_ColHeader]

func newCol(table *s_Table[S_ColHeader], index int) *S_Col {
	return (*S_Col)(newItem(table, index))
}

func (this *S_Col) Order() string {
	return IndexToColOrder(this.index)
}

// 通过表头名称获取单元格
func (this *S_Col) GetCell(name string) *S_Cell {
	header := this.table.nameHeaders[name]
	if header == nil { return nil }
	text := ""
	if header.index < len(this.table.data[this.index]) {
		text = this.table.data[this.index][header.index]
	}
	return newCell(header.name, header.index, this.index, text)
}

// 获取所有单元格
func (this *S_Col) GetCells() []*S_Cell {
	cells := []*S_Cell{}
	for _, h := range this.table.headers {
		text := ""
		if h.index < len(this.table.data[this.index]) {
			text = this.table.data[this.index][h.index]
		}
		cell := newCell(h.name, h.index, this.index, text)
		cells = append(cells, cell)
	}
	return cells
}

// 遍历所有单元格
func (this *S_Col) For(f func(*S_Cell)) {
	for _, h := range this.table.headers {
		text := ""
		if h.index < len(this.table.data[this.index]) {
			text = this.table.data[this.index][h.index]
		}
		cell := newCell(h.name, h.index, this.index, text)
		f(cell)
	}
}

// ---------------------------------------------------------
type S_ColsIter struct{ *s_ItemIter[S_ColHeader] }

func newColsIter(table *s_Table[S_ColHeader]) *S_ColsIter {
	return &S_ColsIter{newItemsIter(table)}
}

// 通过表头名称
func (this *S_ColsIter) Col() *S_Col {
	if this.index < len(this.table.data) {
		return newCol(this.table, this.index)
	}
	return nil
}
