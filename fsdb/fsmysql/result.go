/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: query return result
@author: fanky
@version: 1.0
@date: 2022-01-21
**/

package fsmysql

import (
	"database/sql"
	"fmt"

	"fsky.pro/fsmysql/fssql"

	"fsky.pro/fsreflect"
)

type I_SQLInfo interface {
	Err() error
	SQLText() string
	FmtSQLText(string) string
}

type s_SQLInfo struct {
	sql string
}

func newSQLInfo(sql string) *s_SQLInfo {
	return &s_SQLInfo{sql}
}

func newSQLInfof(sql string, args ...any) *s_SQLInfo {
	return &s_SQLInfo{fmt.Sprintf(sql, args...)}
}

func (this *s_SQLInfo) SQLText() string {
	return this.sql
}

func (this *s_SQLInfo) FmtSQLText(string) string {
	return this.sql
}

func (this *s_SQLInfo) Err() error {
	return nil
}

// -------------------------------------------------------------------
// Result
// -------------------------------------------------------------------
type S_OPResult struct {
	err     error
	sqlInfo I_SQLInfo
}

func newOPResult(sqlInfo I_SQLInfo, err error) *S_OPResult {
	return &S_OPResult{err, sqlInfo}
}

func (this *S_OPResult) Err() error {
	return this.err
}

func (this *S_OPResult) SQLText() string {
	return this.sqlInfo.SQLText()
}

func (this *S_OPResult) FmtSQLText(indent string) string {
	return this.sqlInfo.FmtSQLText(indent)
}

// -------------------------------------------------------------------
// SelectResult
// -------------------------------------------------------------------
type F_IterValues func(error, ...any) bool
type F_IterObjects func(error, any) bool

type S_OPSelectResult struct {
	*S_OPResult
	rows *sql.Rows
}

func newOPSelectResult(sqlInfo *fssql.S_SelectInfo, err error) *S_OPSelectResult {
	return &S_OPSelectResult{
		S_OPResult: newOPResult(sqlInfo, err),
		rows:       nil,
	}
}

func newOPSelectResult2(sqlInfo *fssql.S_SelectInfo, rows *sql.Rows) *S_OPSelectResult {
	return &S_OPSelectResult{
		S_OPResult: newOPResult(sqlInfo, nil),
		rows:       rows,
	}
}

func (this *S_OPSelectResult) ForValues(fun F_IterValues, values ...interface{}) error {
	if this.err != nil {
		return this.err
	}
	defer this.rows.Close()
	return this.ForValuesTx(fun, values...)
}

// 该方法给事务调用
func (this *S_OPSelectResult) ForValuesTx(fun F_IterValues, values ...any) error {
	if this.err != nil {
		return this.err
	}
	for this.rows.Next() {
		err := this.rows.Scan(values...)
		if err != nil {
			if !fun(fmt.Errorf("scan row error: %v", err)) {
				return nil
			}
		} else if !fun(nil, values...) {
			return nil
		}
	}
	return nil
}

func (this *S_OPSelectResult) ForObjects(fun F_IterObjects) error {
	if this.err != nil {
		return this.err
	}
	defer this.rows.Close()
	return this.ForObjectsTx(fun)
}

// 该方法给事务调用
func (this *S_OPSelectResult) ForObjectsTx(fun F_IterObjects) error {
	if this.err != nil {
		return this.err
	}
	sqlInfo := this.sqlInfo.(*fssql.S_SelectInfo)
	obj, mptrs, err := sqlInfo.CreateOutObject()
	if err != nil {
		return err
	}

	for this.rows.Next() {
		err := this.rows.Scan(mptrs...)
		if err != nil {
			if !fun(fmt.Errorf("scan row error: %v", err), nil) {
				return nil
			}
			continue
		}
		rowObj, err := fsreflect.CopyStructObject(obj)
		if !fun(err, rowObj) {
			return nil
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// ExecResult
// -----------------------------------------------------------------------------
// 插入操作返回
type S_OPExecResult struct {
	*S_OPResult
	rest sql.Result
}

func newOPExecResult(sqlInfo *fssql.S_ExecInfo, rest sql.Result, err error) *S_OPExecResult {
	return &S_OPExecResult{
		S_OPResult: newOPResult(sqlInfo, err),
		rest:       rest,
	}
}

func (this *S_OPExecResult) LastInsertId() (int64, error) {
	if this.rest == nil {
		return 0, this.Err()
	}
	return this.rest.LastInsertId()
}

func (this *S_OPExecResult) RowsAffected() (int64, error) {
	if this.rest == nil {
		return 0, this.Err()
	}
	return this.rest.RowsAffected()
}

// -----------------------------------------------------------------------------
// 行值扫描操作
// -----------------------------------------------------------------------------
type S_OPFetchResult struct {
	*S_OPResult
	rows *sql.Rows
}

func newOPFetchResult(sqlInfo *fssql.S_FetchInfo, err error) *S_OPFetchResult {
	return &S_OPFetchResult{
		S_OPResult: newOPResult(sqlInfo, err),
	}
}

func newOPFetchResult2(sqlInfo *fssql.S_FetchInfo, rows *sql.Rows) *S_OPFetchResult {
	return &S_OPFetchResult{
		S_OPResult: newOPResult(sqlInfo, nil),
		rows:       rows,
	}
}

// 如果 fun 返回 false，则停止扫描
// 该方法给事务调用
func (this *S_OPFetchResult) ForTx(fun F_IterValues, values ...any) error {
	if this.err != nil {
		return this.err
	}
	for this.rows.Next() {
		err := this.rows.Scan(values...)
		if err != nil {
			if !fun(fmt.Errorf("scan row error: %v", err)) {
				return nil
			}
		} else if !fun(nil, values...) {
			return nil
		}
	}
	return nil
}

func (this *S_OPFetchResult) For(fun F_IterValues, values ...any) error {
	if this.err != nil {
		return this.err
	}
	defer this.rows.Close()
	return this.For(fun, values...)
}

// -----------------------------------------------------------------------------
// 一次性返回所有值操作
// -----------------------------------------------------------------------------
type S_OPValueResult struct {
	*S_OPResult
	Value any
}

func newOPValueResult(sqlInfo *fssql.S_FetchInfo, err error) *S_OPValueResult {
	return &S_OPValueResult{
		S_OPResult: newOPResult(sqlInfo, err),
	}
}

func newOPValueResult2(sqlInfo *fssql.S_FetchInfo, value any) *S_OPValueResult {
	return &S_OPValueResult{
		S_OPResult: newOPResult(sqlInfo, nil),
		Value:      value,
	}
}
