/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: tx
@author: fanky
@version: 1.0
@date: 2021-12-31
**/

package fsmysql

import (
	"database/sql"
	"fmt"

	"fsky.pro/fsmysql/fssql"
)

type S_Tx struct {
	*s_Operator
	*sql.Tx
}

func newTx(tx *sql.Tx) *S_Tx {
	return &S_Tx{
		s_Operator: newOperator(tx),
		Tx:         tx,
	}
}

// 执行 fssql.ExecSQLInfo 所构建的 SQL 语句
func (this *S_Tx) ExecSQLInfo(sqlInfo *fssql.S_ExecInfo) *S_OPExecResult {
	if sqlInfo.Err() != nil {
		return newOPExecResult(sqlInfo, nil, fmt.Errorf("exec sql fail, %v", sqlInfo.Err()))
	}
	stmt, err := this.wrapper.Prepare(sqlInfo.SQLText())
	if err != nil {
		return newOPExecResult(sqlInfo, nil, err)
	}
	// 事务不需要关闭 stmt
	// defer stmt.Close()
	rest, err := stmt.Exec(sqlInfo.InValues...)
	return newOPExecResult(sqlInfo, rest, err)
}

// 查找包含指定字符串的字段
// reptn 是 mysql 的正则表达式
func (this *S_Tx) FetchColumns(table *fssql.S_Table, reptn string) *S_OPValueResult {
	sqlInfo := table.FetchColumnsSQL(reptn)
	rows, err := this.wrapper.Query(sqlInfo.SQLText())
	if err != nil {
		return newOPValueResult(sqlInfo, err)
	}
	// 事务不需要关闭 rows
	// defer rows.Close()
	var col string
	cols := []string{}
	for rows.Next() {
		if err := rows.Scan(&col); err == nil {
			cols = append(cols, col)
		}
	}
	return newOPValueResult2(sqlInfo, cols)
}
