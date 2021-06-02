/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 格式化 sql 语句
@author: fanky
@version: 1.0
@date: 2021-05-31
**/

package fsmysql

import (
	"fmt"
	"reflect"
	"strings"

	"fsky.pro/fsreflect"
)

// 构建 Inser、Update 语句时，对于值为字符串类型，会在字符串两端加上单引号。
// 如果不需要单引号，则需要用 Unquote 包裹字符串类型值
type Unquote struct {
	v interface{}
}

// FmtSelectPrepare 格式化一个 prepare sql 查询语句
// 参数：
//	colValuePtrs：数据库列名与传出值映射（注意：传出值，必须为指针）
//		colValuePtrs 中的值为 scan 后可存储值的指针
//	from：一般为表名
//	tail：一般为 where、order by、limit 等
// 如，colValuePtrs = {"tb.col1": &a, "tb.col2": &b}、from = "table"，则调用函数后返回：
//		sqltx == "select `tb`.`col1`,`tb`.`col2` from table"
//		valuePtrs == []interface{}{[&a, &b}
// 注意：
//	返回的 sqltx 并不会在语句最后加分号
func FmtSelectPrepare(colValuePtrs map[string]interface{}, from string, tail string) (sqltx string, valuePtrs []interface{}) {
	colNames := make([]string, 0)
	valuePtrs = make([]interface{}, 0)
	for colName, valuePtr := range colValuePtrs {
		colName = "`" + strings.Replace(colName, ".", "`.`", 1) + "`"
		colNames = append(colNames, colName)
		valuePtrs = append(valuePtrs, valuePtr)
	}
	if tail != "" {
		sqltx = fmt.Sprintf("select %s from %s %s", strings.Join(colNames, ","), from, tail)
	} else {
		sqltx = fmt.Sprintf("select %s from %s", strings.Join(colNames, ","), from)
	}
	return
}

// FmtInsertPrepare 格式化一个 Prepare sql 插入语句
// 参数：
//	tbName：表名
//	colValues：数据库列名与插入值映射
// 返回值：
//	values: 为插入列的值数组（按插入时的顺序传出）
// 注意：
//	1、返回的 sqltx 并不会在语句最后加分号
//	2、colValues 中的值，如果是 string、[]byte、[]rune。则会自动添加双引号，如果不需要添加双引号，需要这样写：
//		colValue[xxx] = fsmysql.Unquote{"abc"}
func FmtInsertPrepare(tbName string, colValues map[string]interface{}) (sqltx string, values []interface{}) {
	sqltx = fmt.Sprintf("insert into `%s`(%%s) values(%%s)", tbName)
	colNames := make([]string, 0)
	values = make([]interface{}, 0)
	qms := make([]string, 0)

	for colName, value := range colValues {
		colNames = append(colNames, "`"+colName+"`")
		values = append(values, value)
		if reflect.TypeOf(value) == reflect.TypeOf(Unquote{}) {
			qms = append(qms, fmt.Sprintf("%v", value.(Unquote).v))
		} else if fsreflect.CanConvertToTypeOf(value, "") {
			qms = append(qms, fmt.Sprintf("'%v'", value))
		} else {
			qms = append(qms, fmt.Sprintf("%v", value))
		}
	}
	sqltx = fmt.Sprintf(sqltx, strings.Join(colNames, ","), strings.Join(qms, ","))
	return
}

// FmtUpdatePrepare 格式化一个 prepare sql 更新语句
// 参数：
//	tbName：表名
//	colValues：数据库列名与更新值映射
//	where：条件语句
// 返回值：
//	values: 为更新列的值数组（按更新时的顺序传出）
// 注意：
//	1、返回的 sqltx 并不会在语句最后加分号
//	2、colValues 中的值，如果是 string、[]byte、[]rune。则会自动添加双引号，如果不需要添加双引号，需要这样写：
//		colValue[xxx] = fsmysql.Unquote{"abc"}
func FmtUpdatePrepare(tbName string, colValues map[string]interface{}, where string) (sqltx string, values []interface{}) {
	sqltx = fmt.Sprintf("update `%s` set %%s where %%s", tbName)
	sets := make([]string, 0)
	values = make([]interface{}, 0)

	for colName, value := range colValues {
		if reflect.TypeOf(value) == reflect.TypeOf(Unquote{}) {
			sets = append(sets, fmt.Sprintf("`%s`=%v", colName, value.(Unquote).v))
		} else if fsreflect.CanConvertToTypeOf(value, "") {
			sets = append(sets, fmt.Sprintf("`%s`='%v'", colName, value))
		} else {
			sets = append(sets, fmt.Sprintf("`%s`=%v", colName, value))
		}
		values = append(values, value)
	}
	sqltx = fmt.Sprintf(sqltx, strings.Join(sets, ","), where)
	return
}
