/**
@copyright: fantasysky 2016
@brief: 数据库管理器
@author: fanky
@version: 1.0
@date: 2019-01-12
**/

package fsdb

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

// ---------------------------------------------------------
// I_DBMgr
// ---------------------------------------------------------
type I_DBMgr interface {
	Open(*S_DBInfo) error
	Close()
	Ping() error
	DBStats() *sql.DBStats
}

// ---------------------------------------------------------
// S_DBMgr
// ---------------------------------------------------------
type S_DBMgr struct {
	DBPtr *sql.DB // 连接串
}

func (this *S_DBMgr) Open() error {
	return errors.New("Open method is not implemented")
}

func (this *S_DBMgr) Close() {
	if this.DBPtr != nil {
		this.DBPtr.Close()
	}
}

func (this *S_DBMgr) Opened() bool {
	return this.DBPtr != nil
}

func (this *S_DBMgr) Ping() error {
	if this.DBPtr == nil {
		return errors.New("dbmgr has not initialized")
	}
	return this.DBPtr.Ping()
}

func (this *S_DBMgr) DBStats() *sql.DBStats {
	if this.DBPtr == nil {
		return nil
	}
	stats := this.DBPtr.Stats()
	return &stats
}
