/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: deleter
@author: fanky
@version: 1.0
@date: 2022-01-17
**/

package fssql

type s_Delete struct {
	s_SQL
	table *S_Table // 要更新的表
}

// 更新指定成员值
func Delete(table *S_Table) *s_DeleteWhere {
	this := &s_Delete{
		table: table,
	}
	this.sqlText = "DELETE FROM " + table.quote()
	return (*s_DeleteWhere)(this)
}

// ---------------------------------------------------------
type s_DeleteWhere s_Delete

// 更新条件
func (this *s_DeleteWhere) Where(exp string, args ...interface{}) *s_DeleteEnd {
	exp = this.explainExp(this.table, exp, args...)
	if this.notOK() {
		this.errorf("error delete where condition, %v", this.err.Error())
		return (*s_DeleteEnd)(this)
	}
	this.sqlText += " WHERE " + exp
	return (*s_DeleteEnd)(this)
}

// ---------------------------------------------------------
type s_DeleteEnd s_Delete

func (this *s_DeleteEnd) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}
