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
	"errors"
	"fmt"

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
//		col1 int       `db:"col_1"`           // 对应数据库中的列为：col_1
//		col2 float32                             // 不定义 mysql tag，则在数据库表在的列名与成员名称一致，即：col2
//		col3 string    `db:"-"`               // mysql 的 tag 为减号 “-” 的话，则表示该成员不映射为数据库表中的任何列
// }
// 即：tag 为 “mysql” 或 “db” 都可以，如果 “mysql” 和 “db” 两个 tag 同时存在，优先考虑 “mysql”
// ----------------------------------------------------------------------------
type S_DB struct {
	*sql.DB // 连接串
	DBInfo  *S_DBInfo
}

func open(dbInfo *S_DBInfo) (*sql.DB, error) {
	if err := dbInfo.check(); err != nil {
		return nil, err
	}

	lntext := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		dbInfo.User, dbInfo.Password, dbInfo.host(), dbInfo.port(), dbInfo.DBName, dbInfo.charset())
	// 初始化连接池
	link, err := sql.Open("mysql", lntext)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("create mysql link pool fail: %s", err.Error()))
	}
	return link, nil
}

// 打开数据库连接
func Open(dbInfo *S_DBInfo) (*S_DB, error) {
	db, err := open(dbInfo)
	if err != nil {
		return nil, err
	}
	link := &S_DB{
		DBInfo: dbInfo,
		DB:     db,
	}
	return link, nil
}

// 打开数据库连接池，如果数据库不存在，则创建
func OpenAndCreate(dbInfo *S_DBInfo) (*S_DB, error) {
	db, err := open(dbInfo)
	if err != nil {
		return nil, err
	}

	sqltx := fmt.Sprintf("CREATE DATABASE IF NO EXISTS `%s` DEFAULT CHARACTER SET %s DEFAULT COLLATE %s;", dbInfo.DBName, dbInfo.charset(), dbInfo.collate())
	_, err = db.Exec(sqltx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("create database %q fail: %v", dbInfo.DBName, err)
	}
	link := &S_DB{
		DBInfo: dbInfo,
		DB:     db,
	}
	return link, nil
}

// -------------------------------------------------------------------
// methods
// -------------------------------------------------------------------
// CreateObjectTable 根据对象属性创建数据库表
// 参数：
//  tbName：表名
//	obj：必须是一个结构体对象，或者是结构体对象指针(可以为 nil 指针)，否则产生错误
//	makeups：指定特殊键，譬如主键、外键、唯一键等
//	tbAttr：表格属性，如数据引擎、字符集等
// 返回 result=fsmysql.S_Result
//	result.SQLText()：创建表格的 sql 语句
//	result.Err()：创建成功时，返回 nil，创建失败时，返回具体错误
// 注意：
//	结构成员必须带有类型修饰符，如下：
//	type S struct{
//		Abc int    `mysql:"abc",mysqltd:"int not null default 0"`   // 在数据库中创建的列名为 abc，类型为 int
//		def string `mysqltd:"varchar(64)" default ''`               // 在数据库中创建的列名为 def，类型为 varchar(64)
//		ghi int                                                     // 不指定类型装饰 mysqltd，将不会创建该字段
//	}
// 或：
//	type S struct{
//		Abc int    `db:"abc",mysqltd:"int not null default 0"`      // 在数据库中创建的列名为 abc，类型为 int
//		def string `dbtd:"varchar(64)" default ''`                  // 在数据库中创建的列名为 def，类型为 varchar(64)
//		ghi int                                                     // 不指定类型装饰 mysqltd，将不会创建该字段
//	}
// 即：tag 为 “mysqltd” 或 “dbtd” 都可以，如果 “mysqltd” 和 “dbtd” 两个 tag 同时存在，优先考虑 “mysqltd”
func (this *S_DB) CreateObjectTable(tbName string, obj interface{}, makeups []string, tbAttr string) *S_Result {
	keyTypes, err := dbkeyTypes(obj)
	if err != nil {
		return newResult("", err)
	}
	sqltx := FmtCreateTableSQL(tbName, keyTypes, makeups, tbAttr)
	_, err = this.Exec(sqltx)
	return newResult(sqltx, err)
}

// query 一个 struct 对象，如果对象成员名称与数据库对应列关键之不一样，则需要带 tag，并且 tag 名称必须为：mysql
// 参数：
//	obj：必须是一个非 nil struct 指针，否则返回非 nil error
//	members：指定要 query 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//	from：一般为表名，可以是多表联合查询
//	tail：一般为 where 或 order by、limit 等子句，如果这些子句带有参数，则将参数放到 args 中
//	args：为 from、tail 中可能包含的参数列表
// 返回：
//	resutl.SQLText()：查找数据库时所使用的 sql 语句
//	result.Err()：如果记录不存在，则返回 sql.ErrNoRows，如果是其他查询错误则返回其他错误
//
// 示例：
//	type AA struct {
//		v1 string  `mysql:"col1"`		// col1 为数据库中对应的字段名
//		v2 int     `mysql:"col2"`
//		v3 float32 `mysql:"col3"`
//	}
//	aa := new(AA)
//	_, err := QueryObject(aa, "v2,v3", "table", "where v1=?", aa.v1)	// 只查询 v1 和 v2
//	_, err := QueryObject(aa, "*", "table", "where v1=?", aa.v1)		// 查询全部列
func (this *S_DB) QueryObject(obj interface{}, members string, from, tail string, args ...interface{}) *S_Result {
	tagMems, err := dbkeyMapValues(obj, members, true)
	if err != nil {
		return newResult("", err)
	}
	sqltx, valuePtrs := FmtSelectPrepare(tagMems, from, tail)
	row := this.QueryRow(sqltx+";", args...)
	err = row.Scan(valuePtrs...)
	if err == nil {
		return newResult(sqltx, err)
	}
	return newResult(sqltx, err)
}

// query 所有符合条件的 struct 对象，如果对象成员名称与数据库对应列关键之不一样，则需要带 tag，并且 tag 名称必须为：mysql
// 参数：
//	obj：必须是一个非 nil struct 指针，否则返回非 nil error
//	members：指定要 query 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//	from：一般为表名，可以是多表联合查询
//	tail：一般为 where 或 order by、limit 等子句，如果这些子句带有参数，则将参数放到 args 中
//	args：为 from、tail 中可能包含的参数列表
// 返回：
//	query.Err()：查询错误
//	query.SQLText()：查找数据库时所使用的 sql 语句
//
// 示例：
//	type AA struct {
//		v1 string `mysql:"col1"`		// col1 为数据库中对应的字段名
//		v2 int `mysql:"col2"`
//		v3 float32 `mysql:"col3"`
//	}
//	aa := new(AA)
//	query := QueryObjects(aa, "v2,v3", "table", "where v1=?", aa.v1)	// 只查询 v1 和 v2
//	if query.Err() != nil {
//		obj, err := query.Scan()
//		if err != nil {
//			// error
//		} else {
//			// 这里可以获取 obj 的成员值，但不能对 obj 进行存储使用，该函数每次返回的 obj 都是同一个对象指针，只是每次传回时，其成员值不一样
//		}
//	}
//	query.Close()
//
// 注意：
//	1、通过调用返回值 result 的 Next() 方法可以判断搜索记录是否还存在
//	2、通过调用返回值 resutl 的 Scan() 方法可以获得带有记录值的 object 对象，但是，该 object 正是该函数第一个参数传入的 obj 值
//	   因此每次调用 Scan() 方法返回的都是同一个对象指针，所以，不能直接对返回的 object 进行存储使用！
func (this *S_DB) QueryObjects(obj interface{}, members string, from, tail string, args ...interface{}) *S_QueryResult {
	tagMems, err := dbkeyMapValues(obj, members, true)
	if err != nil {
		return newQueryResult("", err, nil)
	}
	sqltx, valuePtrs := FmtSelectPrepare(tagMems, from, tail)
	rows, err := this.Query(sqltx+";", args...)
	result := newQueryResult(sqltx, err, rows)
	if err != nil {
		result.err = errors.New("query objects fail, " + err.Error())
		return result
	}
	result.valuePtrs = valuePtrs
	result.obj = obj
	return result
}

// ---------------------------------------------------------
// 插入一个 mysql 化对象，如果对象成员名称与数据库对应列关键之不一样，则需要带 tag，并且 tag 名称必须为：mysql
// 参数：
//	obj：必须是一个非 nil struct 指针，否则返回非 nil error
//	members：指定要 insert 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//	objs：每一个 obj 必须是一个非 nil struct 指针
// 返回：
//	result.SQLText()：返回预处理 sql 语句
//	result.Err()：如果构建或者执行 sql 语句错误，则返回非 nil 错误
//	result.RowsAffected()：返回影响的记录数，插入成功影响的记录数为 1，否则为 0
func (this *S_DB) InsertObjects(tbName string, obj interface{}, members string, objs ...interface{}) *S_ExecResult {
	objs = append([]interface{}{obj}, objs...)
	items := make([]map[string]interface{}, 0, len(objs)+1)
	for _, obj := range objs {
		tagMems, err := dbkeyMapValues(obj, members, false)
		if err != nil {
			return newExecResult("", err, nil)
		}
		items = append(items, tagMems)
	}

	sqltx, values := FmtInsertPrepare(tbName, items[0], items[1:]...)
	stmt, err := this.Prepare(sqltx + ";")
	if err != nil {
		return newExecResult(sqltx, err, nil)
	}
	defer stmt.Close()
	rest, err := stmt.Exec(values...)
	return newExecResult(sqltx, err, rest)
}

// 插入一个 mysql 化对象，如果唯一键已经存在，则放弃插入操作，如果对象成员名称与数据库对应列关键之不一样，则需要带 tag，并且 tag 名称必须为：mysql
// 参数：
//	obj：必须是一个非 nil struct 指针，否则返回非 nil error
//	members：指定要 insert 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
// 返回：
//	result.SQLText()：返回预处理 sql 语句
//	result.Err()：如果构建或者执行 sql 语句错误，则返回非 nil 错误
//	result.RowsAffected()：返回影响的记录数，插入成功影响的记录数为 1，否则为 0
func (this *S_DB) InsertIgnoreObject(tbName string, obj interface{}, members string) *S_ExecResult {
	if obj == nil {
		return newExecResult("", errors.New("the obj argument mustn't be a nil pointer."), nil)
	}
	tagMems, err := dbkeyMapValues(obj, members, false)
	if err != nil {
		return newExecResult("", err, nil)
	}
	sqltx, values := FmtInsertIgnorePrepare(tbName, tagMems)
	stmt, err := this.Prepare(sqltx + ";")
	if err != nil {
		return newExecResult(sqltx, err, nil)
	}
	defer stmt.Close()
	rest, err := stmt.Exec(values...)
	return newExecResult(sqltx, err, rest)
}

// 插入一个 mysql 化对象，如果唯一键已经存在，则改为更新操作,如果对象成员名称与数据库对应列关键之不一样，则需要带 tag，并且 tag 名称必须为：mysql
// 参数：
//	iobj：必须是一个非 nil struct 指针，否则返回非 nil error
//	imems：指定要 insert 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//	uobj：必须是一个非 nil struct 指针，否则返回非 nil error
//	umems：指定要 update 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
// 返回：
//	result.SQLText()：返回预处理 sql 语句
//	result.Err()：如果构建或者执行 sql 语句错误，则返回非 nil 错误
//	result.RowsAffected()：返回影响的记录数，插入成功影响的记录数为 1，否则为 0
func (this *S_DB) InsertUpdateObject(tbName string, iobj interface{}, imems string, uobj interface{}, umems string) *S_ExecResult {
	iTagMems, err := dbkeyMapValues(iobj, imems, false)
	if err != nil {
		return newExecResult("", err, nil)
	}
	uTagMems, err := dbkeyMapValues(uobj, umems, false)
	if err != nil {
		return newExecResult("", err, nil)
	}

	sqltx, values := FmtInsertUpdatePrepare(tbName, iTagMems, uTagMems)
	stmt, err := this.Prepare(sqltx + ";")
	if err != nil {
		return newExecResult(sqltx, err, nil)
	}
	defer stmt.Close()
	rest, err := stmt.Exec(values...)
	return newExecResult(sqltx, err, rest)
}

// ---------------------------------------------------------
// 更新一个 mysql 化对象到数据库，如果对象成员名称与数据库对应列关键之不一样，则需要带 tag，并且 tag 名称必须为：mysql
// 参数：
//	tbName: 数据库表名
//	obj：必须是一个非 nil struct 指针，否则返回非 nil error
//	members：指定要 insert 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//	where：where 条件子句（不需要添加 where 关键字）
//	whereArgs: where 子句的参数
// 返回：
//	result.SQLText()：返回预处理 sql 语句
//	result.Err()：如果构建或者执行 sql 语句错误，则返回非 nil 错误
//	result.RowsAffected()：返回影响的记录数，插入成功影响的记录数为 1，否则为 0
func (this *S_DB) UpdateObject(tbName string, obj interface{}, members string, where string, whereArgs ...interface{}) *S_ExecResult {
	tagMems, err := dbkeyMapValues(obj, members, false)
	if err != nil {
		return newExecResult("", err, nil)
	}
	sqltx, values := FmtUpdatePrepare(tbName, tagMems, where, whereArgs...)
	stmt, err := this.Prepare(sqltx)
	if err != nil {
		return newExecResult(sqltx+";", err, nil)
	}
	defer stmt.Close()
	rest, err := stmt.Exec(values...)
	return newExecResult(sqltx, err, rest)
}

// ---------------------------------------------------------
// 删除记录
// 参数：
//	tbName：要删除记录所在表明
//	where：要删除记录的条件 where 子句
//	whereArgs：where 条件子句中的参数
// 返回:
//	sqltx：预处理 sql 语句
//	values：与参数 whereArgs 一致
func (this *S_DB) Delete(tbName string, where string, whereArgs ...interface{}) *S_ExecResult {
	sqltx, values := FmtDeletePrepare(tbName, where, whereArgs...)
	stmt, err := this.Prepare(sqltx + ";")
	if err != nil {
		return newExecResult(sqltx, err, nil)
	}
	defer stmt.Close()
	rest, err := stmt.Exec(values...)
	return newExecResult(sqltx, err, rest)
}
