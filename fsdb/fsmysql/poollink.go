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

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func _open(dbInfo *S_DBInfo) (*sql.DB, error) {
	lntext := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		dbInfo.User, dbInfo.Password, dbInfo.Host, dbInfo.Port, dbInfo.DBName, dbInfo.Encoding)
	// 初始化连接池
	link, err := sql.Open("mysql", lntext)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("create mysql link pool fail: %s", err.Error()))
	}
	link.SetMaxOpenConns(dbInfo.MaxOpenConns)
	link.SetMaxIdleConns(dbInfo.MaxIdleConns)
	link.SetConnMaxLifetime(dbInfo.ConnMaxLifetime)
	return link, nil
}

// -------------------------------------------------------------------
// S_PoolLink
// -------------------------------------------------------------------
type S_PoolLink struct {
	DBInfo *S_DBInfo
	DB     *sql.DB // 连接串
}

// ---------------------------------------------------------
// public
// ---------------------------------------------------------
// 打开数据库连接
func Open(dbInfo *S_DBInfo) (*S_PoolLink, error) {
	db, err := _open(dbInfo)
	if err != nil {
		return nil, err
	}
	link := &S_PoolLink{
		DBInfo: dbInfo,
		DB:     db,
	}
	return link, nil
}

// 打开数据库连接池，如果数据库不存在，则创建
// charset 为数据库默认字符集，如果传入空字符串，则默认使用：utf8mb4
// collate 为数据库默认字符集，如果传入空字符串，则默认使用：utf8mb4_general_ci
func OpenAndCreate(dbInfo *S_DBInfo, charset, collate string) (*S_PoolLink, error) {
	db, err := _open(dbInfo)
	if err != nil {
		return nil, err
	}
	if charset == "" {
		charset = "utf8mb4"
	}
	if collate == "" {
		collate = "utf8mb4_general_ci"
	}

	sqltx := fmt.Sprintf("CREATE DATABASE IF NO EXISTS `%s` DEFAULT CHARACTER SET %s DEFAULT COLLATE %s;", dbInfo.DBName, charset, collate)
	_, err = db.Exec(sqltx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("create database %q fail: %v", dbInfo.DBName, err)
	}
	link := &S_PoolLink{
		DBInfo: dbInfo,
		DB:     db,
	}
	return link, nil
}

// close sql pool
func (this *S_PoolLink) Close() error {
	return this.DB.Close()
}
