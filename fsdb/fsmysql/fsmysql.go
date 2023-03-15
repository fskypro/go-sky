/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: go mysql
@author: fanky
@version: 1.0
@date: 2022-07-10
**/

package fsmysql

import (
	"database/sql"
	"errors"
	"fmt"
)

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

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
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

func OpenWithArgs(host string, port int, user, pwd string, dbName string) (*S_DB, error) {
	dbInfo := &S_DBInfo{
		User:     user,
		Password: pwd,
		Host:     host,
		Port:     port,
		DBName:   dbName,
	}
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
