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

	"fsky.pro/fsdb"
	_ "github.com/go-sql-driver/mysql"
)

type S_DBMgr struct {
	fsdb.S_DBMgr
}

func (this *S_DBMgr) Open(dbInfo *fsdb.S_DBInfo) error {
	lntext := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		dbInfo.User, dbInfo.Password, dbInfo.Host, dbInfo.Port, dbInfo.DBName, dbInfo.Encoding)

	// 初始化连接池
	dbPtr, err := sql.Open("mysql", lntext)
	if err != nil {
		return errors.New(fmt.Sprintf("can't create mysql link pool: %s", err.Error()))
	}

	// 连接测试
	if err = dbPtr.Ping(); err != nil {
		dbPtr.Close()
		return errors.New(fmt.Sprintf("can't connect to mysql database: %s", err.Error()))
	}
	dbPtr.SetMaxOpenConns(dbInfo.MaxOpenConns)
	dbPtr.SetMaxIdleConns(dbInfo.MaxIdleConns)
	dbPtr.SetConnMaxLifetime(dbInfo.ConnMaxLifetime)
	this.DBPtr = dbPtr
	return nil
}

func (this *S_DBMgr) OpenWithLink(link string) error {
	dbPtr, err := sql.Open("mysql", link)
	if err != nil {
		return errors.New(fmt.Sprintf("can't create mysql link pool: %s", err.Error()))
	}

	// 连接测试
	if err = dbPtr.Ping(); err != nil {
		dbPtr.Close()
		return errors.New(fmt.Sprintf("can't connect to mysql database: %s", err.Error()))
	}
	return nil
}
