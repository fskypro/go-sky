/**
@copyright: fantasysky 2016
@brief: sql 语句格式化
@author: fanky
@version: 1.0
@date: 2019-03-10
**/

package fssql

import "fmt"
import "strings"

// FmtSelectPrepare 格式化一个 prepare sql 查询语句
// 参数：
//	tbName：表名
//	colValuePtrs：数据库列名与传出值映射（注意：传出值，必须为指针）
//	vm：prepare 查询值标识符，mysql 为 ?
func FmtSelectPrepare(tbName string, colValuePtrs map[string]interface{}, where, vm string) (sqltx string, valuePtrs []interface{}) {
	sqltx = fmt.Sprintf("select %%s from %s where %%s", tbName)
	colNames := make([]string, 0)
	valuePtrs = make([]interface{}, 0)
	for colName, valuePtr := range colValuePtrs {
		colNames = append(colNames, colName)
		valuePtrs = append(valuePtrs, valuePtr)
	}
	sqltx = fmt.Sprintf(sqltx, strings.Join(colNames, ","), where)
	return
}

// FmtInsertPrepare 格式化一个 Prepare sql 插入语句
// 参数：
//	tbName：表名
//	colValues：数据库列名与插入值映射
//	vm：prepare 插入值表示符，mysql 为 ?
// 返回值：
//	values: 为插入列的值数组（按插入时的顺序传出）
func FmtInsertPrepare(tbName string, colValues map[string]interface{}, vm string) (sqltx string, values []interface{}) {
	sqltx = fmt.Sprintf("insert into %s(%%s) values(%%s)", tbName)
	colNames := make([]string, 0)
	values = make([]interface{}, 0)
	qms := make([]string, 0)
	for colName, value := range colValues {
		colNames = append(colNames, colName)
		values = append(values, value)
		qms = append(qms, vm)
	}
	sqltx = fmt.Sprintf(sqltx, strings.Join(colNames, ","), strings.Join(qms, ","))
	return
}

// FmtUpdatePrepare 格式化一个 prepare sql 更新语句
// 参数：
//	tbName：表名
//	colValues：数据库列名与更新值映射
//	where：条件语句
//	vm：propare 插入值标识符，mysq 为 ?
// 返回值：
//	values: 为更新列的值数组（按更新时的顺序传出）
func FmtUpdatePrepare(tbName string, colValues map[string]interface{}, where, vm string) (sqltx string, values []interface{}) {
	sqltx = fmt.Sprintf("update %s set %%s where %%s", tbName)
	sets := make([]string, 0)
	values = make([]interface{}, 0)
	for colName, value := range colValues {
		sets = append(sets, fmt.Sprintf("%s=?", colName))
		values = append(values, value)
	}
	sqltx = fmt.Sprintf(sqltx, strings.Join(sets, ","), where)
	return
}
