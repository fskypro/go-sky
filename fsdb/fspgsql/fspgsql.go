/**
@copyright: fantasysky 2016
@brief: go postgresql
@author: fanky
@version: 1.0
@date: 2023-02-01
**/

package fspgsql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

func open(linktext string) (*sql.DB, error) {
	sqldb, err := sql.Open("postgres", linktext)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't create postgresql link pool: %s", err.Error()))
	}

	// 连接测试
	if err = sqldb.Ping(); err != nil {
		sqldb.Close()
		return nil, errors.New(fmt.Sprintf("can't connect to postgresql database: %s", err.Error()))
	}
	return sqldb, nil
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func Open(dbInfo *S_DBInfo) (*S_DB, error) {
	// see also:
	//    https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	// 初始化连接池
	sqldb, err := open(dbInfo.LinkText())
	if err != nil {
		return nil, err
	}
	if dbInfo.MaxOpenConns > 0 {
		sqldb.SetMaxOpenConns(dbInfo.MaxOpenConns)
	}
	if dbInfo.MaxIdleConns > 0 {
		sqldb.SetMaxIdleConns(dbInfo.MaxIdleConns)
	}
	if dbInfo.ConnMaxLifetime > 0 {
		sqldb.SetConnMaxLifetime(dbInfo.ConnMaxLifetime)
	}
	return &S_DB{
		s_Operator: newOperator(sqldb),
		DB:         sqldb,
		DBInfo:     dbInfo,
	}, nil
}
