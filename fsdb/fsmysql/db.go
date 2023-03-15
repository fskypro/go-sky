/**
@copyright: fantasysky 2016
@brief: 数据库管理器
@author: fanky
@version: 1.0
@date: 2019-01-12
**/

package fsmysql

import (
	"database/sql"
	"fmt"

	"fsky.pro/fsmysql/fssql"
	_ "github.com/go-sql-driver/mysql"
)

// ----------------------------------------------------------------------------
// S_DB
// S_DB 可以以对象的方式操作数据库，定义对象结构如下：
// type S_MysqlTable struct {
//		col1 int       `mysql:"col_1"`           // 对应数据库中的列为：col_1
//		col2 float32                             // 不定义 mysql tag，则在数据库表在的列名与成员名称一致，即：col2
//		col3 string    `mysql:"-"`               // mysql 的 tag 为减号 “-” 的话，则表示该成员不映射为数据库表中的任何列
// }
// 或：
// type S_MysqlTable struct {
//		col1 int       `db:"col_1"`              // 对应数据库中的列为：col_1
//		col2 float32                             // 不定义 mysql tag，则在数据库表在的列名与成员名称一致，即：col2
//		col3 string    `db:"-"`                  // mysql 的 tag 为减号 “-” 的话，则表示该成员不映射为数据库表中的任何列
// }
// 即：tag 为 “mysql” 或 “db” 都可以，如果 “mysql” 和 “db” 两个 tag 同时存在，优先考虑 “mysql”
// ----------------------------------------------------------------------------
type S_DB struct {
	*s_Operator
	*sql.DB // 连接串
	DBInfo  *S_DBInfo
}

// 启动事务
func (this *S_DB) Begin() (*S_Tx, error) {
	tx, err := this.DB.Begin()
	if err != nil {
		return nil, err
	}
	return newTx(tx), nil
}

// 执行 fssql.ExecSQLInfo 所构建的 SQL 语句
func (this *S_DB) ExecSQLInfo(sqlInfo *fssql.S_ExecInfo) *S_OPExecResult {
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

// 查找包含指定字符串的字段
func (this *S_DB) FetchColumns(table *fssql.S_Table, like string) *S_OPValueResult {
	sqlInfo := table.FetchColumnsSQL(like)
	rows, err := this.wrapper.Query(sqlInfo.SQLText())
	if err != nil {
		return newOPValueResult(sqlInfo, err)
	}
	defer rows.Close()
	var col string
	cols := []string{}
	for rows.Next() {
		if err := rows.Scan(&col); err != nil {
			cols = append(cols, col)
		}
	}
	return newOPValueResult2(sqlInfo, cols)
}
