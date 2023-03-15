/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: database sql
@author: fanky
@version: 1.0
@date: 2022-07-10
**/

package fssql

import "fmt"

// 构建创建数据库 sql 语句
func CreateDBSQL(name string, charset string) *S_ExecInfo {
	sqltx := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` charset %s", name, charset)
	return newExecInfo2(nil, sqltx)
}

// 获取指定表格名称 sql 语句，如果 like 参数传入空串，则获取所有表格名称
func FetchTablesSQL(like string) *S_FetchInfo {
	var sqltx string
	if like == "" {
		sqltx = "SHOW TABLES"
	} else {
		sqltx = fmt.Sprintf("SHOW TABLES LIKE '%s'", like)
	}
	return newFetchInfo(nil, sqltx)
}
