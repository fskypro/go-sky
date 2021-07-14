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
)

// FmtCreateTableSQL 格式化一个创建表格的 SQL 语句
// 参数：
//	tbName：表格名称
//	cols：[[列名1, 类型和默认值、修饰等], [列名2, 类型和默认值、修饰等], ...}
//	makeups：指定特殊键，譬如主键、外键、唯一键等
//	tbAttr：表格属性，如数据引擎、字符集等
// 返回：
//	返回创建表格的 sql 语句，不带分号结尾
func FmtCreateTableSQL(tbName string, cols [][2]string, makeups []string, tbAttr string) string {
	sqltx := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(%%s) %s", tbName, tbAttr)
	items := make([]string, 0)
	for _, col := range cols {
		items = append(items, fmt.Sprintf("`%s` %s", col[0], col[1]))
	}
	if makeups != nil {
		for _, extra := range makeups {
			items = append(items, extra)
		}
	}
	return fmt.Sprintf(sqltx, strings.Join(items, ","))
}

// -------------------------------------------------------------------
// FmtSelectPrepare 格式化一个 select prepare sql 查询语句
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
		sqltx = fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(colNames, ","), from, tail)
	} else {
		sqltx = fmt.Sprintf("SELECT %s FROM %s", strings.Join(colNames, ","), from)
	}
	return
}

// -------------------------------------------------------------------
// FmtInsertPrepare 格式化一个可以插入多条记录的 insert Prepare sql 插入语句
// 参数：
//	tbName：表名
//	colValues：要插入的列和值映射：{列名:列值}
// 返回值：
//	values: 为插入列的值数组（按插入时的顺序传出）
// 注意：
//	返回的 sqltx 并不会在语句最后加分号
// 注意：
//	1、插入的列以参数 colValues 的 key 为准
//	2、如果某个 colValuess 中不存在某些 key，则使用 colValues 中对应 key 值所属类型的默认值
func FmtInsertPrepare(tbName string, colValues map[string]interface{}, colValuess ...map[string]interface{}) (sqltx string, values []interface{}) {
	sqltx = fmt.Sprintf("INSERT INTO `%s`(`%%s`) VALUES%%s", tbName)
	colNames := make([]string, 0)
	values = make([]interface{}, 0)
	items := make([]string, 0, len(colValuess)+1)

	qms := make([]string, 0)
	for colName, value := range colValues {
		colNames = append(colNames, colName)
		values = append(values, value)
		qms = append(qms, "?")
	}
	items = append(items, "("+strings.Join(qms, ",")+")")

	for _, cvs := range colValuess {
		for _, colName := range colNames {
			value, ok := cvs[colName]
			if !ok {
				rv := reflect.ValueOf(colValues[colName])
				value = reflect.New(rv.Type()).Elem().Interface()
			}
			values = append(values, value)
		}
		items = append(items, "("+strings.Join(qms, ",")+")")
	}
	sqltx = fmt.Sprintf(sqltx, strings.Join(colNames, "`,`"), strings.Join(items, ","))
	return
}

// FmtInsertIgnorePrepare 格式化一个 insert ignore prepare sql 插入语句
// 参数：
//	tbName：表名
//	colValues：要插入的列和值映射：{列名:列值}
// 返回值：
//	values: 为插入列的值数组（按插入时的顺序传出）
// 注意：
//	返回的 sqltx 并不会在语句最后加分号
func FmtInsertIgnorePrepare(tbName string, colValues map[string]interface{}) (sqltx string, values []interface{}) {
	sqltx = fmt.Sprintf("INSERT IGNORE `%s`(%%s) VALUES(%%s)", tbName)
	colNames := make([]string, 0)
	values = make([]interface{}, 0)
	qms := make([]string, 0)

	for colName, value := range colValues {
		colNames = append(colNames, "`"+colName+"`")
		values = append(values, value)
		qms = append(qms, "?")
	}
	sqltx = fmt.Sprintf(sqltx, strings.Join(colNames, ","), strings.Join(qms, ","))
	return
}

// FmtInsertUpdatePrepare 格式化一个 insert on duplicate key update prepare sql 插入语句
// 参数：
//	tbName：表名
//	colIValues：要插入的列和值映射：{列名:列值}
//	colUValues：要更新的列和值映射：{列名:列值}
// 返回值：
//	values: 为插入和更新列的值数组（按插入时的顺序传出）
// 注意：
//	返回的 sqltx 并不会在语句最后加分号
func FmtInsertUpdatePrepare(tbName string, colIValues map[string]interface{}, colUValues map[string]interface{}) (sqltx string, values []interface{}) {
	sqltx = fmt.Sprintf("INSERT INTO `%s`(%%s) VALUES(%%s) ON DUPLICATE KEY UPDATE %%s", tbName)
	colNames := make([]string, 0)
	values = make([]interface{}, 0)
	iqms := make([]string, 0)
	uqms := make([]string, 0)

	for colName, value := range colIValues {
		colNames = append(colNames, "`"+colName+"`")
		values = append(values, value)
		iqms = append(iqms, "?")
	}
	for colName, value := range colUValues {
		values = append(values, value)
		uqms = append(uqms, fmt.Sprintf("`%s`=?", colName))
	}
	sqltx = fmt.Sprintf(sqltx,
		strings.Join(colNames, ","),
		strings.Join(iqms, ","),
		strings.Join(uqms, ","))
	return
}

// -------------------------------------------------------------------
// FmtUpdatePrepare 格式化一个 update prepare sql 更新语句
// 参数：
//	tbName：表名
//	colValues：数据库列名与更新值映射
//	where：条件语句
//	whereArgs：where 子句中的参数
// 返回值：
//	values: 为更新列的值数组（按更新时的顺序传出）
// 注意：
//	返回的 sqltx 并不会在语句最后加分号
func FmtUpdatePrepare(tbName string, colValues map[string]interface{}, where string, whereArgs ...interface{}) (sqltx string, values []interface{}) {
	sqltx = fmt.Sprintf("UPDATE `%s` SET %%s WHERE %%s", tbName)
	sets := make([]string, 0)
	values = make([]interface{}, 0)

	for colName, value := range colValues {
		sets = append(sets, fmt.Sprintf("`%s`=?", colName))
		values = append(values, value)
	}
	values = append(values, whereArgs...)
	sqltx = fmt.Sprintf(sqltx, strings.Join(sets, ","), where)
	return
}

// -------------------------------------------------------------------
// FmtDeletePrepare 格式化一个 delete prepare sql 删除语句
// 参数：
//	tbName：要删除记录所在表明
//	where：要删除记录的条件 where 子句
//	whereArgs：where 条件子句中的参数
// 返回:
//	sqltx：预处理 sql 语句
//	values：与参数 whereArgs 一致
func FmtDeletePrepare(tbName string, where string, whereArgs ...interface{}) (sqltx string, values []interface{}) {
	return fmt.Sprintf("DELETE FROM `%s` WHERE %s", tbName, where), whereArgs
}
