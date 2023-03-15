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
	"fsky.pro/fsky"
)

// -----------------------------------------------------------------------------
// DB/Tx Wraper
// -----------------------------------------------------------------------------
type i_DBWrapper interface {
	Exec(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
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
func (this *s_Operator) QueryDBStatus(likes ...string) (map[string]string, error) {
	status := map[string]string{}
	sqlText := "SHOW STATUS"
	args := []any{}
	for idx, like := range likes {
		sqlText += fsky.IfElse(idx == 0, " WHERE", " OR") + " `Variable_name` LIKE ?"
		args = append(args, like)
	}
	rows, err := this.wrapper.Query(sqlText, args...)
	if err != nil {
		return status, err
	}
	defer rows.Close()

	var name, value string
	for rows.Next() {
		err = rows.Scan(&name, &value)
		if err == nil {
			status[name] = value
		}
	}
	return status, nil
}

func (this *s_Operator) QueryDBVariables(likes ...string) (map[string]string, error) {
	status := map[string]string{}
	sqlText := "SHOW VARIABLES"
	args := []any{}
	for idx, like := range likes {
		sqlText += fsky.IfElse(idx == 0, " WHERE", " OR") + " `Variable_name` LIKE ?"
		args = append(args, like)
	}

	rows, err := this.wrapper.Query(sqlText, args...)
	if err != nil {
		return status, err
	}
	defer rows.Close()

	var name, value string
	for rows.Next() {
		err = rows.Scan(&name, &value)
		if err == nil {
			status[name] = value
		}
	}
	return status, nil
}

func (this *s_Operator) SetDBVariable(name string, value any) error {
	sqlText := fmt.Sprintf("SET GLOBAL %s=?", name)
	_, err := this.wrapper.Exec(sqlText, value)
	return err
}

// -------------------------------------------------------------------
// 通过对象映射表创建数据库表
// -------------------------------------------------------------------
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
func (this *s_Operator) SelectRowValue(sqlInfo *fssql.S_SelectInfo, outValues ...any) *S_OPResult {
	if sqlInfo.Err() != nil {
		return newOPResult(sqlInfo, fmt.Errorf("select object member values fail, %v", sqlInfo.Err()))
	}
	row := this.wrapper.QueryRow(sqlInfo.SQLText(), sqlInfo.InValues...)
	err := row.Scan(outValues...)
	return newOPResult(sqlInfo, err)
}

// 查找单个对象
func (this *s_Operator) SelectRowObject(sqlInfo *fssql.S_SelectInfo, outObj any) *S_OPResult {
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
		return newOPSelectResult(sqlInfo, sqlInfo.Err())
	}
	rows, err := this.wrapper.Query(sqlInfo.SQLText(), sqlInfo.InValues...)
	result := newOPSelectResult2(sqlInfo, rows)
	if err != nil {
		result.err = errors.New("query objects fail, " + err.Error())
		return result
	}
	return result
}

// -------------------------------------------------------------------
// fetch
// -------------------------------------------------------------------
func (this *s_Operator) FetchRow(sqlInfo *fssql.S_FetchInfo, outValues ...any) *S_OPResult {
	if sqlInfo.Err() != nil {
		return newOPResult(sqlInfo, fmt.Errorf("search row value fail, %v", sqlInfo.Err()))
	}
	row := this.wrapper.QueryRow(sqlInfo.SQLText(), sqlInfo.InValues...)
	err := row.Scan(outValues...)
	return newOPResult(sqlInfo, err)
}

func (this *s_Operator) Fetch(sqlInfo *fssql.S_FetchInfo) *S_OPFetchResult {
	if sqlInfo.Err() != nil {
		return newOPFetchResult(sqlInfo, sqlInfo.Err())
	}
	rows, err := this.wrapper.Query(sqlInfo.SQLText(), sqlInfo.InValues...)
	if err != nil {
		return newOPFetchResult(sqlInfo, err)
	}
	return newOPFetchResult2(sqlInfo, rows)
}

// -------------------------------------------------------------------
// table scheme
// -------------------------------------------------------------------
// 表是否存在
func (this *s_Operator) HasTable(name string) *S_OPValueResult {
	sqlInfo := fssql.FetchTablesSQL(name)
	row := this.wrapper.QueryRow(sqlInfo.SQLText())
	var tmp string
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return newOPValueResult2(sqlInfo, false)
	}
	if err != nil {
		return newOPValueResult(sqlInfo, err)
	}
	return newOPValueResult2(sqlInfo, true)
}

// 修改表名
func (this *s_Operator) RenameTable(oldName, newName string) *S_OPResult {
	sqlInfo := newSQLInfof("RENAME TABLE `%s` TO `%s`", oldName, newName)
	_, err := this.wrapper.Exec(sqlInfo.SQLText())
	return newOPResult(sqlInfo, err)
}

// ---------------------------------------------------------
// 给指定表格添加字段
func (this *s_Operator) AddColumn(table *fssql.S_Table, colName, colType string, tail string) *S_OPResult {
	sqlInfo := table.AddColumnSQL(colName, colType, tail)
	_, err := this.wrapper.Exec(sqlInfo.SQLText())
	return newOPResult(sqlInfo, err)
}

// 删除指定表格的某字段
func (this *s_Operator) DelColumn(table *fssql.S_Table, colName string) *S_OPResult {
	sqlInfo := table.DelColumnsSQL(colName)
	_, err := this.wrapper.Exec(sqlInfo.SQLText())
	return newOPResult(sqlInfo, err)
}

// 重命名字段
func (this *s_Operator) RenameColumn(table *fssql.S_Table, oldName, newName string, ctype string, tail string) *S_OPResult {
	sqlInfo := table.RenameColumnSQL(oldName, newName, ctype, tail)
	_, err := this.wrapper.Exec(sqlInfo.SQLText())
	return newOPResult(sqlInfo, err)
}

// ---------------------------------------------------------
// 添加唯一约束
func (this *s_Operator) AddUniqueKey(table *fssql.S_Table, uniqueName string, mnames ...string) *S_OPResult {
	sqlInfo := table.AddUniqueKey(uniqueName, mnames...)
	_, err := this.wrapper.Exec(sqlInfo.SQLText())
	return newOPResult(sqlInfo, err)
}

// 删除唯一约束
func (this *s_Operator) DropUniqueKey(table *fssql.S_Table, uniqueName string) *S_OPResult {
	sqlInfo := table.DropUniqueKey(uniqueName)
	_, err := this.wrapper.Exec(sqlInfo.SQLText())
	return newOPResult(sqlInfo, err)
}

// ---------------------------------------------------------
// 判断指定表是否存在指定名称的约束
func (this *s_Operator) FetchConstrains(dbName string, table *fssql.S_Table, consName string) *S_OPValueResult {
	sqlInfo := table.Constrains(dbName, consName)
	rest := this.Fetch(sqlInfo)
	if rest.Err() != nil {
		return newOPValueResult(sqlInfo, rest.Err())
	}
	var con string
	cons := []string{}
	err := rest.For(func(err error, vs ...any) bool {
		if err != nil {
			cons = append(cons, *(vs[0].(*string)))
		}
		return true
	}, &con)
	return newOPValueResult(sqlInfo, err)
}

// 添加表外键约束
func (this *s_Operator) AddForeignKey(fkName string, table *fssql.S_Table, mName string, mainMember *fssql.S_Member, constrain string) *S_OPResult {
	sqlInfo := table.AddForeignKey(fkName, mName, mainMember, constrain)
	_, err := this.wrapper.Exec(sqlInfo.SQLText())
	return newOPResult(sqlInfo, err)
}
