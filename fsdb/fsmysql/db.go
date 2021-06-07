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
	"reflect"
	"strings"
	"unsafe"

	_ "github.com/go-sql-driver/mysql"
)

// ----------------------------------------------------------------------------
// private
// ----------------------------------------------------------------------------
// 获取结构体成员 tag 与成员值或成员指针指针的映射：{tag:member ptr}
// 如果 ptr 参数为 true，则取成员的指针，否则取成员的值
func dbkeyMapValues(obj interface{}, members string, ptr bool) (tagMems map[string]interface{}, err error) {
	v := reflect.ValueOf(obj)
	if v.Type().Kind() != reflect.Ptr {
		err = errors.New("object type must be a pointer of struct.")
		return
	}
	v = v.Elem()
	if v.Type().Kind() != reflect.Struct {
		err = errors.New("object must be a struct instancei pointer.")
		return
	}

	tagMems = make(map[string]interface{})
	t := v.Type()
	var tfield reflect.StructField
	var vfield reflect.Value
	var dbkey string

	// = 号开头，成员名称与数据库关键之一致
	ignoreTag := false
	if len(members) > 1 && members[0:2] == "=|" {
		members = members[2:]
		ignoreTag = true
	}

	members = strings.TrimSpace(members)
	if members == "*" || members == "" {
		for i := 0; i < t.NumField(); i++ {
			tfield = t.Field(i)
			vfield = v.Field(i)
			if ignoreTag {
				dbkey = tfield.Name
			} else {
				dbkey = tfield.Tag.Get("mysql")
				if dbkey == "" {
					continue
				}
			}
			up := unsafe.Pointer(vfield.UnsafeAddr())
			pv := reflect.NewAt(vfield.Type(), up)
			if ptr {
				tagMems[dbkey] = pv.Interface()
			} else {
				tagMems[dbkey] = pv.Elem().Interface()
			}
		}
		return
	}

	for _, m := range strings.Split(members, ",") {
		m := strings.TrimSpace(m)
		tfield, ok := t.FieldByName(m)
		if !ok {
			err = fmt.Errorf("object has no field named %q.", m)
			return
		}
		vfield = v.FieldByName(m)
		if ignoreTag {
			dbkey = tfield.Name
		} else {
			dbkey = tfield.Tag.Get("mysql")
			if dbkey == "" {
				continue
			}
		}
		up := unsafe.Pointer(vfield.UnsafeAddr())
		pv := reflect.NewAt(vfield.Type(), up)
		if ptr {
			tagMems[dbkey] = pv.Interface()
		} else {
			tagMems[dbkey] = pv.Elem().Interface()
		}
	}
	return
}

// ----------------------------------------------------------------------------
// S_DB
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
// query 一个 struct 对象，如果对象成员名称与数据库对应列关键之不一样，则需要带 tag，并且 tag 名称必须为：mysql
// 参数：
//	obj：必须是一个 struct 指针，否则返回非 nil error
//	members：指定要 query 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//		如果 members 以 =| 开头，则可以不给结构成员添加 tag，表示结构中所有成员名称与数据库对列名称一致
//	from：一般为表名，可以是多表联合查询
//	tail：一般为 where 或 order by、limit 等子句，如果这些子句带有参数，则将参数放到 args 中
//	args：为 from、tail 中可能包含的参数列表
// 返回：
//	resutl.SQLText()：查找数据库时所使用的 sql 语句
//	result.Err()：如果记录不存在，则返回 sql.ErrNoRows
//
// 示例：
//	type AA struct {
//		v1 string `mysql:"col1"`		// col1 为数据库中对应的字段名
//		v2 int `mysql:"col2"`
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
//	obj：必须是一个 struct 指针，否则返回非 nil error
//	members：指定要 query 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//		如果 members 以 =| 开头，则可以不给结构成员添加 tag，表示结构中所有成员名称与数据库对列名称一致
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
//		a := new(AA)
//		err = query.Scan(a)
//	}
//	query.Close()
func (this *S_DB) QueryObjects(obj interface{}, members string, from, tail string, args ...interface{}) *S_QueryResult {
	tagMems, err := dbkeyMapValues(obj, members, true)
	if err != nil {
		return newQueryResult("", err, nil)
	}
	sqltx, valuePtrs := FmtSelectPrepare(tagMems, from, tail)
	rows, err := this.Query(sqltx, args...)
	scanner := newQueryResult(sqltx+";", err, rows)
	if err != nil {
		scanner.err = errors.New("query objects fail, " + err.Error())
		return scanner
	}
	scanner.valuePtrs = valuePtrs
	scanner.obj = obj
	return scanner
}

// ---------------------------------------------------------
// 插入一个 mysql 化对象，如果对象成员名称与数据库对应列关键之不一样，则需要带 tag，并且 tag 名称必须为：mysql
// 参数：
//	obj：必须是一个 struct 指针，否则返回非 nil error
//	members：指定要 insert 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//		如果 members 以 =| 开头，则可以不给结构成员添加 tag，表示结构中所有成员名称与数据库对列名称一致
// 返回：
//	result.SQLText()：返回预处理 sql 语句
//	result.Err()：如果构建或者执行 sql 语句错误，则返回非 nil 错误
//	result.RowsAffected()：返回影响的记录数，插入成功影响的记录数为 1，否则为 0
func (this *S_DB) InsertObject(tbName string, obj interface{}, members string, objs ...interface{}) *S_ExecResult {
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
//	obj：必须是一个 struct 指针，否则返回非 nil error
//	members：指定要 insert 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//		如果 members 以 =| 开头，则可以不给结构成员添加 tag，表示结构中所有成员名称与数据库对列名称一致
// 返回：
//	result.SQLText()：返回预处理 sql 语句
//	result.Err()：如果构建或者执行 sql 语句错误，则返回非 nil 错误
//	result.RowsAffected()：返回影响的记录数，插入成功影响的记录数为 1，否则为 0
func (this *S_DB) InsertIgnoreObject(tbName string, obj interface{}, members string) *S_ExecResult {
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
//	iobj：必须是一个 struct 指针，否则返回非 nil error
//	imems：指定要 insert 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//		如果 members 以 =| 开头，则可以不给结构成员添加 tag，表示结构中所有成员名称与数据库对列名称一致
//	uobj：必须是一个 struct 指针，否则返回非 nil error
//	umems：指定要 update 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
// 返回：
//	result.SQLText()：返回预处理 sql 语句
//	result.Err()：如果构建或者执行 sql 语句错误，则返回非 nil 错误
//	result.Affected：返回影响的记录数，插入成功影响的记录数为 1，否则为 0
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
//	obj：必须是一个 struct 指针，否则返回非 nil error
//	members：指定要 insert 的成员，多个成员用逗号隔开，如果要 query 全部成员可以："" 或 "*"
//		如果 members 以 =| 开头，则可以不给结构成员添加 tag，表示结构中所有成员名称与数据库对列名称一致
//	where：where 条件子句（不需要添加 where 关键字）
//	whereArgs: where 子句的参数
// 返回：
//	result.SQLText()：返回预处理 sql 语句
//	result.Err()：如果构建或者执行 sql 语句错误，则返回非 nil 错误
//	result.Affected：返回影响的记录数，插入成功影响的记录数为 1，否则为 0
func (this *S_DB) UpdateObject(tbName string, obj interface{}, members string, where string, whereArgs ...interface{}) *S_ExecResult {
	tagMems, err := dbkeyMapValues(obj, members, false)
	if err != nil {
		return newExecResult("", err, nil)
	}
	sqltx, values := FmtUpdatePrepare(tbName, tagMems, where, whereArgs)
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
