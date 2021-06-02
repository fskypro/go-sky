/**
@copyright: fantasysky 2016
@brief: 数据库初始化信息
@author: fanky
@version: 1.0
@date: 2019-03-15
**/

package fsmysql

import "time"

type S_DBInfo struct {
	Encoding        string        // 连接编码
	MaxOpenConns    int           // 打开数据库的最大连接数
	MaxIdleConns    int           // 连接池中保持连接的最大连接数
	ConnMaxLifetime time.Duration // 控线连接的最长有效期

	Host     string
	Port     int
	User     string
	Password string

	DBName  string // 数据库名称，可以省略
	Timeout uint   // 超时（秒）
}
