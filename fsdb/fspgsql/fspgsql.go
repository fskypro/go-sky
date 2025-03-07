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
	"net"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type s_Dialer struct {
	*net.Dialer
	timeout time.Duration
}

func (this *s_Dialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}

func open(dbInfo *S_DBInfo) (*sql.DB, error) {
	dialer := &s_Dialer{
		Dialer: &net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP: net.ParseIP(dbInfo.LocalAddr), // 本地网卡 IP
			},
			Timeout: time.Duration(dbInfo.Timeout) * time.Second,
		},
	}
	// 使用自定义拨号器创建 Connector
	conn, err := pq.NewConnector(dbInfo.ConnString())
	if err != nil {
		return nil, fmt.Errorf("Failed to create connector: %v", err)
	}

	conn.Dialer(dialer)

	// 使用 Connector 创建数据库连接
	sqldb := sql.OpenDB(conn)

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
	sqldb, err := open(dbInfo)
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
