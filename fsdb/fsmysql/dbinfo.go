/**
@copyright: fantasysky 2016
@brief: 数据库初始化信息
@author: fanky
@version: 1.0
@date: 2019-03-15
**/

package fsmysql

import "errors"

type S_DBInfo struct {
	Charset string // 连接编码，默认：utf8mb4
	Collate string // 字符检索策略，默认：<Charset>_general_ci

	Host     string // 数据库主机地址，默认为 localhost
	Port     int    // 连接端口，默认为 3306
	Password string // 登录密码
	// 必填字段
	User   string // 登录用户名
	DBName string // 数据库名称
}

func (this *S_DBInfo) check() error {
	if this.User == "" {
		return errors.New("db user mustn't be empty.")
	}
	if this.DBName == "" {
		return errors.New("database name mustn't be emptya.")
	}
	return nil
}

func (this *S_DBInfo) host() string {
	if this.Host == "" {
		return "localhost"
	}
	return this.Host
}

func (this *S_DBInfo) port() int {
	if this.Port <= 0 {
		return 3306
	}
	return this.Port
}

func (this *S_DBInfo) charset() string {
	if this.Charset == "" {
		return "utf8mb4"
	}
	return this.Charset
}

func (this *S_DBInfo) collate() string {
	if this.Collate == "" {
		return this.Charset + "_general_ci"
	}
	return this.Collate
}
