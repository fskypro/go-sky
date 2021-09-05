/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: operate db result
@author: fanky
@version: 1.0
@date: 2021-06-04
**/

package fsmysql

import (
	"database/sql"
	"fmt"

	"fsky.pro/fsreflect"
)

// ----------------------------------------------------------------------------
// S_Result
// ----------------------------------------------------------------------------
type S_Result struct {
	sqltx string
	err   error
}

func newResult(sqltx string, err error) *S_Result {
	return &S_Result{
		sqltx: sqltx,
		err:   err,
	}
}

func (this *S_Result) SQLText() string {
	return this.sqltx
}
func (this *S_Result) Err() error {
	return this.err
}

// -----------------------------------------------------------------------------
// S_QueryResult
// -----------------------------------------------------------------------------
// 查询返回结果
type S_QueryResult struct {
	*S_Result
	rows      *sql.Rows     // sql.db 返回的 Rows 对象
	valuePtrs []interface{} // obj 中所有需要查找成员的指针
	obj       interface{}   // 需要查找的结构体
}

func newQueryResult(sqltx string, err error, rows *sql.Rows) *S_QueryResult {
	return &S_QueryResult{
		S_Result: newResult(sqltx, err),
		rows:     rows,
	}
}

func (this *S_QueryResult) ForEach(fun func(interface{}, error) bool) {
	if this.rows == nil {
		return
	}
	defer this.rows.Close()
	for this.rows.Next() {
		err := this.rows.Scan(this.valuePtrs...)
		if err != nil {
			if !fun(nil, fmt.Errorf("scan row error: %v", err)) {
				return
			}
			continue
		}
		obj, err := fsreflect.CopyStructObject(this.obj)
		if err != nil {
			if !fun(nil, fmt.Errorf("deserialize to object error: %v", err)) {
				return
			}
			continue
		}
		if !fun(obj, nil) {
			return
		}
	}
}

// -----------------------------------------------------------------------------
// S_ExecResult
// -----------------------------------------------------------------------------
// 插入操作返回
type S_ExecResult struct {
	*S_Result
	rest sql.Result
}

func newExecResult(sqltx string, err error, rest sql.Result) *S_ExecResult {
	return &S_ExecResult{
		S_Result: newResult(sqltx, err),
		rest:     rest,
	}
}

func (this *S_ExecResult) LastInsertId() (int64, error) {
	if this.rest == nil {
		return 0, this.Err()
	}
	return this.rest.LastInsertId()
}

func (this *S_ExecResult) RowsAffected() (int64, error) {
	if this.rest == nil {
		return 0, this.Err()
	}
	return this.rest.RowsAffected()
}
