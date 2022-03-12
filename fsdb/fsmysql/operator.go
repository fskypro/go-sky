/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: operator interface
@author: fanky
@version: 1.0
@date: 2021-12-31
**/

package fsmysql

import (
	"database/sql"
	"errors"
	"fmt"

	"fsky.pro/fsmysql/fssql"
)

// -----------------------------------------------------------------------------
// DB/Tx Wraper
// -----------------------------------------------------------------------------
type i_DBWrapper interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
}

// -----------------------------------------------------------------------------
// Operator
// -----------------------------------------------------------------------------
type s_Operator struct {
	wrapper i_DBWrapper
}

func newOperator(wrapper i_DBWrapper) *s_Operator {
	return &s_Operator{wrapper}
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
// 通过对象映射表创建数据库表
func (this *s_Operator) CreateTable(table *fssql.S_Table) *S_OPResult {
	sqlInfo := table.CreateTableSQLInfo()
	if sqlInfo.Err() != nil {
		return newOPResult(sqlInfo, fmt.Errorf("create table fail, %v", sqlInfo.Err()))
	}
	_, err := this.wrapper.Exec(sqlInfo.SQLText())
	return newOPResult(sqlInfo, err)
}

// -------------------------------------------------------------------
// select
// -------------------------------------------------------------------
// 查找一行，并将结果送到传出值
func (this *s_Operator) SelectOneValue(sqlInfo *fssql.S_SelectInfo, outValues ...interface{}) *S_OPResult {
	if sqlInfo.Err() != nil {
		return newOPResult(sqlInfo, fmt.Errorf("select object fail, %v", sqlInfo.Err()))
	}
	row := this.wrapper.QueryRow(sqlInfo.SQLText(), sqlInfo.InValues...)
	err := row.Scan(outValues...)
	return newOPResult(sqlInfo, err)
}

// 查找单个对象
func (this *s_Operator) SelectOneObject(sqlInfo *fssql.S_SelectInfo, outObj interface{}) *S_OPResult {
	if sqlInfo.Err() != nil {
		return newOPResult(sqlInfo, fmt.Errorf("select object fail, %v", sqlInfo.Err()))
	}
	outs, err := sqlInfo.ScanMembers(outObj)
	if err != nil {
		return newOPResult(sqlInfo, fmt.Errorf("select object fail, %v", err))
	}
	row := this.wrapper.QueryRow(sqlInfo.SQLText(), sqlInfo.InValues...)
	err = row.Scan(outs...)
	return newOPResult(sqlInfo, err)
}

// 查找符合条件的所有对象
func (this *s_Operator) Select(sqlInfo *fssql.S_SelectInfo) *S_OPSelectResult {
	if sqlInfo.Err() != nil {
		return newOPSelectResult(sqlInfo, nil, sqlInfo.Err())
	}
	rows, err := this.wrapper.Query(sqlInfo.SQLText(), sqlInfo.InValues...)
	result := newOPSelectResult(sqlInfo, rows, err)
	if err != nil {
		result.err = errors.New("query objects fail, " + err.Error())
		return result
	}
	return result
}

// ---------------------------------------------------------
// 执行 fssql.ExecSQLInfo 所构建的 SQL 语句
func (this *s_Operator) ExecSQLInfo(sqlInfo *fssql.S_ExecInfo) *S_OPExecResult {
	if sqlInfo.Err() != nil {
		return newOPExecResult(sqlInfo, nil, fmt.Errorf("exec sql fail, %v", sqlInfo.Err()))
	}
	stmt, err := this.wrapper.Prepare(sqlInfo.SQLText())
	if err != nil {
		return newOPExecResult(sqlInfo, nil, err)
	}
	defer stmt.Close()
	rest, err := stmt.Exec(sqlInfo.InValues...)
	return newOPExecResult(sqlInfo, rest, err)
}
