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
	this.sqlText = "UPDATE " + fsstr.JoinFunc(tables, ",", func(e interface{}) string { return e.(*S_Table).quote() })
	return (*s_UpdateSetExp)(this)
}

// -------------------------------------------------------------------
// Set
// -------------------------------------------------------------------
type s_UpdateSet s_Update

// 以值更新表记录
func (this *s_UpdateSet) Set(values ...interface{}) *s_UpdateWhere {
	if this.notOK() {
		return (*s_UpdateWhere)(this)
	}
	if len(values) != len(this.members) {
		this.errorf("the number of update values must be %d", len(this.members))
		return (*s_UpdateWhere)(this)
	}
	this.addInValues(values...)
	items := []string{}
	for _, m := range this.members {
		items = append(items, fmt.Sprintf("%s=?", m.quote()))
	}
	this.sqlText += " SET " + strings.Join(items, ",")
	return (*s_UpdateWhere)(this)
}

// 用表对应对象更新记录
func (this *s_UpdateSet) SetObject(obj interface{}) *s_UpdateWhere {
	if this.notOK() {
		return (*s_UpdateWhere)(this)
	}
	vobj := reflect.ValueOf(obj)
	tobj := reflect.TypeOf(obj)
	if tobj.Kind() == reflect.Ptr {
		vobj = vobj.Elem()
		tobj = tobj.Elem()
	}
	if tobj != this.table.tobj {
		this.errorf("input object type is not the same as the table %s binds object type", this.table)
		return (*s_UpdateWhere)(this)
	}

	items := []string{}
	for _, m := range this.members {
		this.addInValues(m.value(vobj))
		items = append(items, fmt.Sprintf("%s=?", m.quote()))
	}
	this.sqlText += " SET " + strings.Join(items, ",")
	return (*s_UpdateWhere)(this)
}

func (this *s_UpdateSet) SetExp(exp string, args ...interface{}) *s_UpdateWhere {
	return (*s_UpdateSetExp)(this).SetExp(exp, args...)
}

// ---------------------------------------------------------
type s_UpdateSetExp s_Update

// 更新关联表，如：
// update t1, t2 set t2.a = t1.a where t2.id = ti.id;
// SetExp("$[1]=$[2]", t1.M("A"), t2.M("A"))
func (this *s_UpdateSetExp) SetExp(exp string, args ...interface{}) *s_UpdateWhere {
	if this.notOK() {
		return (*s_UpdateWhere)(this)
	}
	exp = this.explainExp(this.table, exp, args...)
	if this.notOK() {
		this.errorf("error update exp, %v", this.err.Error())
		return (*s_UpdateWhere)(this)
	}
	this.sqlText += " SET " + exp
	return (*s_UpdateWhere)(this)
}

// -------------------------------------------------------------------
// Where
// -------------------------------------------------------------------
type s_UpdateWhere s_Update

// 更新条件
func (this *s_UpdateWhere) Where(exp string, args ...interface{}) *s_UpdateEnd {
	if this.notOK() {
		return (*s_UpdateEnd)(this)
	}
	exp = this.explainExp(this.table, exp, args...)
	if this.notOK() {
		this.errorf("error update where condition, %v", this.err.Error())
		return (*s_UpdateEnd)(this)
	}
	this.sqlText += " WHERE " + exp
	return (*s_UpdateEnd)(this)
}

func (p *s_UpdateWhere) End() *S_ExecInfo {
	return (*s_UpdateEnd)(p).End()
}

// -------------------------------------------------------------------
// End
// -------------------------------------------------------------------
type s_UpdateEnd s_Update

func (this *s_UpdateEnd) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}
