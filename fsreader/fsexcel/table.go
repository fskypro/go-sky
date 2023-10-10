/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: sheet as table
@author: fanky
@version: 1.0
@date: 2023-09-13
**/

package fsexcel

// -------------------------------------------------------------------
// 表头
// -------------------------------------------------------------------
// 列表头
type S_Header struct {
	name  string // 表头名称
	index int    // 对应列索引
}

func (this *S_Header) Name() string {
	return this.name
}

type S_RowHeader struct {
	S_Header
}

func (this * S_RowHeader) Order() string {
	return IndexToColOrder(this.index)
}

type S_ColHeader struct {
	S_Header
}

func (this *S_ColHeader) Order() int {
	return this.index + 1
}

func newHeader[T S_RowHeader | S_ColHeader](name string, index int) *T  {
	return &T{S_Header{name, index}}
}

// -------------------------------------------------------------------
// 带表头的表格
// -------------------------------------------------------------------
type s_Table[T S_RowHeader | S_ColHeader] struct {
	headers     []*T
	nameHeaders map[string]*T
	data        [][]string
	validIndex  int
}

func newTable[T S_RowHeader | S_ColHeader](data [][]string, headerIndex, validIndex int) *s_Table[T] {
	if validIndex >= headerIndex {
		validIndex = headerIndex + 1
	}
	headers := []*T{}
	nameHeaders := map[string]*T{}
	if data == nil || headerIndex >= len(data) {
		data = [][]string{}
	} else {
		for idx, text := range data[headerIndex] {
			if text == "" { continue }
			header := newHeader[T](text, idx)
			nameHeaders[text] = header
			headers = append(headers, header)
		}
	}
	return &s_Table[T]{
		headers:     headers,
		nameHeaders: nameHeaders,
		data:        data,
		validIndex:  validIndex,
	}
}

func (this *s_Table[T]) Headers() []*T {
	return this.headers[:]
}

// 总记录数
func (this *s_Table[T]) Count() int {
	return max(len(this.data)-this.validIndex-1, 0)
}

// 通过表头名称获取表头
func (this *s_Table[T])GetHeaderByName(name string) *T {
	return this.nameHeaders[name]
} 

// -------------------------------------------------------------------
// 以某一行为表头的表格
// -------------------------------------------------------------------
type S_RowsTable struct{ *s_Table[S_RowHeader] }

func newRowsTable(data [][]string, headerIndex, validIndex int) *S_RowsTable {
	return &S_RowsTable{newTable[S_RowHeader](data, headerIndex, validIndex)}
}

// 通过表头序号获取表头
func (this *S_RowsTable) GetHeaderByOrder(order string) *S_RowHeader {
	for _ , header := range this.headers {
		if header.Order() == order {
			return header
		}
	}
	return nil
}

// 获取指定序号的行
// 注意：
//   第一行为 1
//   如果序号超出范围，则返回 nil
func (this *S_RowsTable) GetRow(order int) *S_Row {
	if order <= this.validIndex { return nil }
	if order > len(this.data)   { return nil }
	return newRow(this.s_Table, order-1)
}

// 获取一个新的记录迭代器
func (this *S_RowsTable) GetIter() *S_RowsIter {
	return newRowsIter(this.s_Table)
}

// 遍历
func (this *S_RowsTable) For(f func(*S_Row) bool) {
	count := len(this.data)
	for i := this.validIndex; i < count; i++ {
		if !f(newRow(this.s_Table, i)) {
			break
		}
	}
}

// -------------------------------------------------------------------
// 以某一列为表头的表格
// -------------------------------------------------------------------
type S_ColsTable struct{ *s_Table[S_ColHeader] }

func newColsTable(data [][]string, headerIndex, validIndex int) *S_ColsTable {
	return &S_ColsTable{newTable[S_ColHeader](data, headerIndex, validIndex)}
}

func (this *S_ColsTable) GetHeaderByOrder(order int ) *S_ColHeader {
	for _, header := range this.headers {
		if header.Order() == order {
			return header
		}
	}
	return nil
}

// 获取指定序号的列
// 注意：
//   第一列为 A
//   如果序号超出范围，则返回 nil
func (this *S_ColsTable) GetCol(order string) *S_Col {
	index := ColOrderToIndex(order)
	if index < this.validIndex { return nil }
	if index >= len(this.data) { return nil }
	return newCol(this.s_Table, index)
}

func (this *S_ColsTable) GetIter() *S_ColsIter {
	return newColsIter(this.s_Table)
}

// 遍历
func (this *S_ColsTable) For(f func(*S_Col) bool) {
	count := len(this.data)
	for i := this.validIndex; i < count; i++ {
		if !f(newCol(this.s_Table, i)) {
			break
		}
	}
}
