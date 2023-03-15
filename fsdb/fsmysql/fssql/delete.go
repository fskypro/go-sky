/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: deleter
@author: fanky
@version: 1.0
@date: 2022-01-17
**/

package fssql

import (
	"fmt"
	"strings"
)

type s_Delete struct {
	s_SQL
	table  *S_Table // 要更新的表
	whered bool
}

// 更新指定成员值
func Delete(table *S_Table) *S_DeleteWhere {
	this := &s_Delete{
		table: table,
	}
	this.sqlText = "DELETE FROM " + table.quote()
	return (*S_DeleteWhere)(this)
}

// -------------------------------------------------------------------
// DeleteWehere
// -------------------------------------------------------------------
type S_DeleteWhere s_Delete

func (this *S_DeleteWhere) where(exp string, args ...interface{}) (string, bool) {
	if this.notOK() {
		return "", false
	}
	exp = this.explainExp(this.table, exp, args...)
	if this.notOK() {
		this.errorf("error delete where condition, %v", this.err.Error())
		return "", false
	}
	return exp, true
}

func (this *S_DeleteWhere) concat(link string, exp string) *S_DeleteWhere {
	if !this.whered {
		this.sqlText += " WHERE "
		this.whered = true
		this.sqlText += exp
		return this
	}
	if strings.HasSuffix(this.sqlText, "(") {
		this.sqlText += exp
	} else {
		this.sqlText += fmt.Sprintf(" %s %s", link, exp)
	}
	return this
}

// ---------------------------------------------------------
// 前括号
func (this *S_DeleteWhere) Quote() *S_DeleteWhere {
	return this.concat("", "(")
}

// 与前括号
func (this *S_DeleteWhere) AndQuote() *S_DeleteWhere {
	return this.concat("AND", "(")
}

// 或前括号
func (this *S_DeleteWhere) OrQuote() *S_DeleteWhere {
	return this.concat("OR", "(")
}

// 后括号
func (this *S_DeleteWhere) RQuote() *S_DeleteWhere {
	this.sqlText += ")"
	return this
}

func (this *S_DeleteWhere) Where(exp string, args ...interface{}) *S_DeleteWhere {
	return this.AndWhere(exp, args...)
}

func (this *S_DeleteWhere) AndWhere(exp string, args ...interface{}) *S_DeleteWhere {
	exp, ok := this.where(exp, args...)
	if !ok {
		return this
	}
	return this.concat("AND", exp)
}

func (this *S_DeleteWhere) OrWhere(exp string, args ...interface{}) *S_DeleteWhere {
	exp, ok := this.where(exp, args...)
	if !ok {
		return this
	}
	return this.concat("OR", exp)
}

func (this *S_DeleteWhere) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}

// ---------------------------------------------------------
type s_DeleteEnd s_Delete

func (this *s_DeleteEnd) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}
