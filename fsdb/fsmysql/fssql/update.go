/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: updater
@author: fanky
@version: 1.0
@date: 2022-01-17
**/

package fssql

import (
	"fmt"
	"reflect"
	"strings"

	"fsky.pro/fsstr"
)

type s_Update struct {
	s_SQL
	table   *S_Table    // 要更新的表
	members []*S_Member // 要更新的对象成员名

	whered bool
}

// -------------------------------------------------------------------
// Update
// 可调用链：
//   更新单个表单条记录：
//	   [Update()/UpdateAll()/UpdateBesides()].[Set()/SetObject()/SetExp()].Where().End()
//	   [Update()/UpdateAll()/UpdateBesides()].[Set()/SetObject()/SetExp()].End()
//
//   用一个表更新另一个表：
//	   UpdateTables().SetExp().Where().End()
//	   UpdateTables().SetExp().End()
//
// 注意：
//   End() 标示构建SQL语句结束
// -------------------------------------------------------------------
// 更新指定成员值
func Update(table *S_Table, mnames ...string) *s_UpdateSet {
	if len(mnames) == 0 {
		return UpdateAll(table)
	}

	this := &s_Update{
		table:   table,
		members: make([]*S_Member, 0),
	}
	this.sqlText = "UPDATE " + table.quote()
	for _, name := range mnames {
		m := this.table.Member(name)
		if m == nil {
			this.errorf("table %s has no object member named %q", this.table, name)
			return (*s_UpdateSet)(this)
		}
		this.members = append(this.members, m)
	}
	return (*s_UpdateSet)(this)
}

// 更新所有成员值
func UpdateAll(table *S_Table) *s_UpdateSet {
	this := &s_Update{
		table:   table,
		members: make([]*S_Member, 0),
	}
	this.sqlText = "UPDATE " + table.quote()
	this.members = this.table.orderMembers
	return (*s_UpdateSet)(this)
}

// 更新除指定成员以外的其他所有成员值
func UpdateBesides(table *S_Table, mnames ...string) *s_UpdateSet {
	if len(mnames) == 0 {
		return UpdateAll(table)
	}

	this := &s_Update{
		table:   table,
		members: make([]*S_Member, 0),
	}
	this.sqlText = "UPDATE " + table.quote()

L:
	for _, m := range this.table.orderMembers {
		for _, n := range mnames {
			if n == m.name {
				continue L
			}
		}
		this.members = append(this.members, m)
	}
	return (*s_UpdateSet)(this)
}

// ---------------------------------------------------------
// 更新关联表，如：
// update t1, t2 set t2.a = t1.a where t2.id = ti.id;
// 注意：调用 UpdateTables 后，必须要调用 SetExp
func UpdateTables(tb1 *S_Table, tb2 *S_Table, tbs ...*S_Table) *s_UpdateSetExp {
	tables := append([]*S_Table{tb1, tb2}, tbs...)
	this := &s_Update{
		table:   tb1,
		members: make([]*S_Member, 0),
	}
	this.sqlText = "UPDATE " + fsstr.JoinFunc(tables, ",", func(t *S_Table) string { return t.quote() })
	return (*s_UpdateSetExp)(this)
}

// -------------------------------------------------------------------
// Set
// -------------------------------------------------------------------
type s_UpdateSet s_Update

// 以值更新表记录
func (this *s_UpdateSet) Set(values ...interface{}) *S_UpdateWhere {
	if this.notOK() {
		return (*S_UpdateWhere)(this)
	}
	if len(values) != len(this.members) {
		this.errorf("the number of update values must be %d", len(this.members))
		return (*S_UpdateWhere)(this)
	}
	this.addInValues(values...)
	items := []string{}
	for _, m := range this.members {
		items = append(items, fmt.Sprintf("%s=?", m.quote()))
	}
	this.sqlText += " SET " + strings.Join(items, ",")
	return (*S_UpdateWhere)(this)
}

// 用表对应对象更新记录
func (this *s_UpdateSet) SetObject(obj interface{}) *S_UpdateWhere {
	if this.notOK() {
		return (*S_UpdateWhere)(this)
	}
	vobj := reflect.ValueOf(obj)
	tobj := reflect.TypeOf(obj)
	if tobj.Kind() == reflect.Ptr {
		vobj = vobj.Elem()
		tobj = tobj.Elem()
	}
	if tobj != this.table.tobj {
		this.errorf("input object type is not the same as the table %s binds object type", this.table)
		return (*S_UpdateWhere)(this)
	}

	items := []string{}
	for _, m := range this.members {
		this.addInValues(m.value(vobj))
		items = append(items, fmt.Sprintf("%s=?", m.quote()))
	}
	this.sqlText += " SET " + strings.Join(items, ",")
	return (*S_UpdateWhere)(this)
}

func (this *s_UpdateSet) SetExp(exp string, args ...interface{}) *S_UpdateWhere {
	return (*s_UpdateSetExp)(this).SetExp(exp, args...)
}

// ---------------------------------------------------------
type s_UpdateSetExp s_Update

// 更新关联表，如：
// update t1, t2 set t2.a = t1.a where t2.id = ti.id;
// SetExp("$[1]=$[2]", t1.M("A"), t2.M("A"))
func (this *s_UpdateSetExp) SetExp(exp string, args ...interface{}) *S_UpdateWhere {
	if this.notOK() {
		return (*S_UpdateWhere)(this)
	}
	exp = this.explainExp(this.table, exp, args...)
	if this.notOK() {
		this.errorf("error update exp, %v", this.err.Error())
		return (*S_UpdateWhere)(this)
	}
	this.sqlText += " SET " + exp
	return (*S_UpdateWhere)(this)
}

// -------------------------------------------------------------------
// Where
// -------------------------------------------------------------------
type S_UpdateWhere s_Update

func (this *S_UpdateWhere) where(exp string, args ...interface{}) (string, bool) {
	if this.notOK() {
		return "", false
	}
	exp = this.explainExp(this.table, exp, args...)
	if this.notOK() {
		this.errorf("error update where condition, %v", this.err.Error())
		return "", false
	}
	return exp, true
}

func (this *S_UpdateWhere) concat(link string, exp string) *S_UpdateWhere {
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

// -------------------------------------------------------------------
// 前括号
func (this *S_UpdateWhere) Quote() *S_UpdateWhere {
	return this.concat("", "(")
}

// 与前括号
func (this *S_UpdateWhere) AndQuote() *S_UpdateWhere {
	return this.concat("AND", "(")
}

// 或前括号
func (this *S_UpdateWhere) OrQuote() *S_UpdateWhere {
	return this.concat("OR", "(")
}

// 后括号
func (this *S_UpdateWhere) RQuote() *S_UpdateWhere {
	this.sqlText += ")"
	return this
}

func (this *S_UpdateWhere) Where(exp string, args ...interface{}) *S_UpdateWhere {
	return this.AndWhere(exp, args...)
}

func (this *S_UpdateWhere) AndWhere(exp string, args ...interface{}) *S_UpdateWhere {
	exp, ok := this.where(exp, args...)
	if !ok {
		return this
	}
	return this.concat("AND", exp)
}

func (this *S_UpdateWhere) OrWhere(exp string, args ...interface{}) *S_UpdateWhere {
	exp, ok := this.where(exp, args...)
	if !ok {
		return this
	}
	return this.concat("OR", exp)
}

func (p *S_UpdateWhere) End() *S_ExecInfo {
	return (*s_UpdateEnd)(p).End()
}

// -------------------------------------------------------------------
// End
// -------------------------------------------------------------------
type s_UpdateEnd s_Update

func (this *s_UpdateEnd) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}
