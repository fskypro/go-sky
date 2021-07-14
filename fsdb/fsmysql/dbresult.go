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
	"errors"
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

func (this *S_QueryResult) Next() bool {
	if this.rows == nil {
		return false
	}
	return this.rows.Next()
}

func (this *S_QueryResult) Scan() (interface{}, error) {
	if this.rows == nil {
		return nil, this.err
	}
	err := this.rows.Scan(this.valuePtrs...)
	if err != nil {
		return nil, err
	}
	return this.obj, nil
}

func (this *S_QueryResult) Close() error {
	if this.rows == nil {
		return errors.New("no rows to close.")
	}
	return this.rows.Close()
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
