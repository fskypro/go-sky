/**
@copyright: fantasysky 2016
@brief: 数据库管理器
@author: fanky
@version: 1.0
@date: 2021-04-02
**/

package fspostgres

import (
	"database/sql"
	"errors"
	"fmt"

	"fsky.pro/fsdb"
	_ "github.com/lib/pq"
)

type S_DBMgr struct {
	fsdb.S_DBMgr
}

func (this *S_DBMgr) Open(dbInfo *fsdb.S_DBInfo) error {
	// see also:
	//    https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING
	lntext := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbInfo.Host, dbInfo.Port, dbInfo.User, dbInfo.Password, dbInfo.DBName)
	if dbInfo.Timeout > 0 {
		lntext += fmt.Sprintf(" connect_timeout=%d", dbInfo.Timeout)
	}

	// 初始化连接池
	dbPtr, err := sql.Open("postgres", lntext)
	if err != nil {
		return errors.New(fmt.Sprintf("can't create postgresql link pool: %s", err.Error()))
	}

	// 连接测试
	if err = dbPtr.Ping(); err != nil {
		dbPtr.Close()
		return errors.New(fmt.Sprintf("can't connect to postgresql database: %s", err.Error()))
	}
	dbPtr.SetMaxOpenConns(dbInfo.MaxOpenConns)
	dbPtr.SetMaxIdleConns(dbInfo.MaxIdleConns)
	dbPtr.SetConnMaxLifetime(dbInfo.ConnMaxLifetime)
	this.DBPtr = dbPtr
	return nil
}

/*
使用 CA 证书时，需要对 PostgreSQL 中的配置进行设置：
	postgresql.conf
		ssl = on
		ssl_ca_file = '/etc/postgres/security/root.crt'
		ssl_cert_file = '/etc/postgres/security/server.crt'
		ssl_key_file = '/etc/postgres/security/server.key'
		password_encryption = scram-sha-256

	pg_hba.conf
		local     all      all                md5
		host      all      all  127.0.0.1/32  md5
		hostssl   all      all  0.0.0.0/0     cert clientcert=1
*/
func (this *S_DBMgr) OpenWithLink(link string) error {
	/*
	   connection := fmt.Sprint(
	           " host=localhost",
	           " port=5432",
	           " user=pguser",
	           " dbname=securitylearning",
	           " sslmode=verify-full",
	           " sslrootcert=root.crt",
	           " sslkey=client.key",
	           " sslcert=client.crt",
	       )
	*/
	// 初始化连接池
	dbPtr, err := sql.Open("postgres", link)
	if err != nil {
		return errors.New(fmt.Sprintf("can't create postgresql link pool: %s", err.Error()))
	}

	// 连接测试
	if err = dbPtr.Ping(); err != nil {
		dbPtr.Close()
		return errors.New(fmt.Sprintf("can't connect to postgresql database: %s", err.Error()))
	}
	return nil
}
