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

// ----------------------------------------------------------------------------
// functions
// ----------------------------------------------------------------------------
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
		s_Operator: newOperator(db),
		DBInfo:     dbInfo,
		DB:         db,
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
		s_Operator: newOperator(db),
		DBInfo:     dbInfo,
		DB:         db,
	}
	return link, nil
}
