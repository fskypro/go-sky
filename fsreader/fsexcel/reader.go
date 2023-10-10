/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: excel reader
@author: fanky
@version: 1.0
@date: 2023-09-12
**/

package fsexcel

import (
	"io"

	"github.com/xuri/excelize/v2"
)

// -------------------------------------------------------------------
// 读取
// -------------------------------------------------------------------
type S_Excel struct {
	*excelize.File
	sheets []*S_Sheet
}

func newExcel(excelFile *excelize.File) *S_Excel {
	excel := &S_Excel{
		File:   excelFile,
		sheets: make([]*S_Sheet, 0),
	}
	for index, name := range excelFile.GetSheetList() {
		excel.sheets = append(excel.sheets, &S_Sheet{
			owner: excel,
			name:  name,
			index: index,
		})
	}
	return excel
}

// 通过索引获取分页(第一页为 0)
func (this *S_Excel) GetSheetByIndex(index int) *S_Sheet {
	if index < len(this.sheets) {
		return this.sheets[index]
	}
	return nil
}

// 通过分页名称获取分页
func (this *S_Excel) GetSheetByName(name string) *S_Sheet {
	for _, sheet := range this.sheets {
		if name == sheet.name {
			return sheet
		}
	}
	return nil
}

// 获取所有分页
func (this *S_Excel) Sheets() []*S_Sheet {
	return this.sheets
}

// ---------------------------------------------------------
// 打开 excel 文件流
func OpenReader(r io.Reader) (*S_Excel, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	return newExcel(f), nil
}

// 打开 excel 文件
func OpenFile(path string) (*S_Excel, error) {
	excel, err := excelize.OpenFile(path)
	if err != nil { return nil, err }
	return newExcel(excel), nil
}
