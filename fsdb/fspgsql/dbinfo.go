/**
@copyright: fantasysky 2016
@brief: 数据库初始化信息
@author: fanky
@version: 1.0
@date: 2023-02-01
**/

package fspgsql

import (
	"fmt"
	"strings"
	"time"
)

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
type S_DBInfo struct {
	MaxOpenConns    int           // 打开数据库的最大连接数
	MaxIdleConns    int           // 连接池中保持连接的最大连接数
	ConnMaxLifetime time.Duration // 控线连接的最长有效期
	SSLMode         bool          // 是否使用 SSL 私密连接
	Timeout         int           // 超时（秒）
	LocalAddr       string        // 本地地址

	Host     string // 数据库主机地址
	Port     int    // 数据库连接端口
	User     string // 连接用户
	Password string // 连接密码
	DBName   string // 数据库名称，可以省略
}

func (this *S_DBInfo) ConnString() string {
	segs := []string{}
	if strings.TrimSpace(this.Host) != "" {
		segs = append(segs, "host="+strings.TrimSpace(this.Host))
	}
	if this.Port > 0 {
		segs = append(segs, fmt.Sprintf("port=%d", this.Port))
	}
	if strings.TrimSpace(this.DBName) != "" {
		segs = append(segs, "dbname="+strings.TrimSpace(this.DBName))
	}
	if this.SSLMode {
		segs = append(segs, "sslmode=enable")
	} else {
		segs = append(segs, "sslmode=disable")
	}
	if this.Timeout > 0 {
		segs = append(segs, fmt.Sprintf("connect_timeout=%d", this.Timeout))
	}
	segs = append(segs, "user="+this.User)
	if this.Password != "" {
		segs = append(segs, "password="+this.Password)
	}
	return strings.Join(segs, " ")
}
