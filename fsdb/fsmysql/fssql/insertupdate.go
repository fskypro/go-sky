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
)

type s_InsertUpdate struct {
	s_SQL
	table         *S_Table    // 要更新的表
	insertMembers []*S_Member // 要插入的对象成员名
	updateMembers []*S_Member // 要更新的对象成员名
}

// -------------------------------------------------------------------
// InsertUpdate
// 请求插入单条记录，如果记录已经存在则更新
// 可调用链：
//	 [InsertUpdate()/InserUpdateAll()/InsertUpdateBesides()].[Values()/Object()].OrUpdate().[With()/WithObject()].Where().End()
//	 [InsertUpdate()/InserUpdateAll()/InsertUpdateBesides()].[Values()/Object()].OrUpdate().[With()/WithObject()].End()
//
// 注意：
//   End() 标示构建SQL语句结束
// -------------------------------------------------------------------
// 更新指定成员值
func InsertOrUpdate(table *S_Table, mnames ...string) *s_InsertUpdateValues {
	if len(mnames) == 0 {
		return InsertAllOrUpdate(table)
	}

	this := &s_InsertUpdate{
		table:         table,
		insertMembers: make([]*S_Member, 0),
	}
	this.sqlText = "INSERT INTO " + table.quote()
	dbkeys := []string{}
	for _, name := range mnames {
		m := this.table.Member(name)
		if m == nil {
			this.errorf("table %s has no object member named %q", this.table, name)
			return (*s_InsertUpdateValues)(this)
		}
		this.insertMembers = append(this.insertMembers, m)
		dbkeys = append(dbkeys, m.quote())
	}
	this.sqlText += fmt.Sprintf("(%s)", strings.Join(dbkeys, ","))
	return (*s_InsertUpdateValues)(this)
}

// 更新所有成员值
func InsertAllOrUpdate(table *S_Table) *s_InsertUpdateValues {
	this := &s_InsertUpdate{
		table:         table,
		insertMembers: make([]*S_Member, 0),
	}
	this.sqlText = "INSERT INTO " + table.quote()
	this.insertMembers = this.table.orderMembers
	dbkeys := []string{}
	for _, m := range this.insertMembers {
		dbkeys = append(dbkeys, m.quote())
	}
	this.sqlText += fmt.Sprintf("(%s)", strings.Join(dbkeys, ","))
	return (*s_InsertUpdateValues)(this)
}

// 更新除指定成员以外的其他所有成员值
func InsertBesidesOrUpdate(table *S_Table, mnames ...string) *s_InsertUpdateValues {
	if len(mnames) == 0 {
		return InsertAllOrUpdate(table)
	}

	this := &s_InsertUpdate{
		table:         table,
		insertMembers: make([]*S_Member, 0),
	}
	this.sqlText = "INSERT INTO " + table.quote()

	dbkeys := []string{}
L:
	for _, m := range this.table.orderMembers {
		for _, n := range mnames {
			if n == m.name {
				continue L
			}
		}
		this.insertMembers = append(this.insertMembers, m)
		dbkeys = append(dbkeys, m.quote())
	}
	this.sqlText += fmt.Sprintf("(%s)", strings.Join(dbkeys, ","))
	return (*s_InsertUpdateValues)(this)
}

// -------------------------------------------------------------------
// Values/Object
// -------------------------------------------------------------------
type s_InsertUpdateValues s_InsertUpdate

func (this *s_InsertUpdateValues) Values(values ...interface{}) *s_InsertUpdateOrUpdate {
	if this.notOK() {
		return (*s_InsertUpdateOrUpdate)(this)
	}
	count := len(this.insertMembers)
	if len(this.insertMembers) != count {
		this.errorf("the number of insert values is not consistent with the number of columns(count=%d) to be inserted", count)
		return (*s_InsertUpdateOrUpdate)(this)
	}
	this.sqlText += " VALUES(" + strings.Repeat(",?", count)[1:] + ")"
	this.addInValues(values...)
	return (*s_InsertUpdateOrUpdate)(this)
}

func (this *s_InsertUpdateValues) Object(obj interface{}) *s_InsertUpdateOrUpdate {
	if this.notOK() {
		return (*s_InsertUpdateOrUpdate)(this)
	}
	tobj := reflect.TypeOf(obj)
	vobj := reflect.ValueOf(obj)
	if tobj == nil {
		this.errorf("insert object is not allow to be nil")
		return (*s_InsertUpdateOrUpdate)(this)
	}
	if tobj.Kind() == reflect.Ptr {
		if vobj.IsNil() {
			this.errorf("insert object is not allow to be nil")
			return (*s_InsertUpdateOrUpdate)(this)
		}
		vobj = vobj.Elem()
		tobj = tobj.Elem()
	}
	if tobj != this.table.tobj {
		this.errorf("the input object %#v type is not the same as the table's map object type(%s) to be inserted", obj, this.table)
		return (*s_InsertUpdateOrUpdate)(this)
	}
	for _, m := range this.insertMembers {
		this.addInValues(m.value(vobj))
	}
	this.sqlText += " VALUES(" + strings.Repeat(",?", len(this.insertMembers))[1:] + ")"
	return (*s_InsertUpdateOrUpdate)(this)
}

// -------------------------------------------------------------------
// on duplicate key update
// -------------------------------------------------------------------
type s_InsertUpdateOrUpdate s_InsertUpdate

func (this *s_InsertUpdateOrUpdate) OrUpdate(mnames ...string) *s_InsertUpdateWith {
	for _, name := range mnames {
		m := this.table.Member(name)
		if m == nil {
			this.errorf("table %s has no object member named %q", this.table, name)
			return (*s_InsertUpdateWith)(this)
		}
		this.updateMembers = append(this.updateMembers, m)
	}
	this.sqlText += " ON DUPLICATE KEY UPDATE"
	return (*s_InsertUpdateWith)(this)
}

func (this *s_InsertUpdateOrUpdate) OrUpdateAll() *s_InsertUpdateWith {
	this.updateMembers = this.table.orderMembers
	this.sqlText += " ON DUPLICATE KEY UPDATE"
	return (*s_InsertUpdateWith)(this)
}

func (this *s_InsertUpdateOrUpdate) OrUpdateBesides(mnames ...string) *s_InsertUpdateWith {
L:
	for _, m := range this.table.orderMembers {
		for _, name := range mnames {
			if name == m.name {
				continue L
			}
		}
		this.updateMembers = append(this.updateMembers, m)
	}
	this.sqlText += " ON DUPLICATE KEY UPDATE"
	return (*s_InsertUpdateWith)(this)
}

// -------------------------------------------------------------------
// update values/object
// -------------------------------------------------------------------
type s_InsertUpdateWith s_InsertUpdate

func (this *s_InsertUpdateWith) With(values ...interface{}) *s_InsertUpdateEnd {
	count := len(this.updateMembers)
	if len(values) != count {
		this.errorf("the number of update values is not consistent with the number of columns(count=%d) to be inserted", count)
		return (*s_InsertUpdateEnd)(this)
	}
	items := []string{}
	for _, m := range this.updateMembers {
		items = append(items, fmt.Sprintf("%s=?", m.quote()))
	}
	this.sqlText += " " + strings.Join(items, ",")
	this.addInValues(values...)
	return (*s_InsertUpdateEnd)(this)
}

// 用表对应对象更新记录
func (this *s_InsertUpdateWith) WithObject(obj interface{}) *s_InsertUpdateEnd {
	if this.notOK() {
		return (*s_InsertUpdateEnd)(this)
	}
	vobj := reflect.ValueOf(obj)
	tobj := reflect.TypeOf(obj)
	if tobj.Kind() == reflect.Ptr {
		vobj = vobj.Elem()
		tobj = tobj.Elem()
	}
	if tobj != this.table.tobj {
		this.errorf("update object type is not the same as the table %s binds object type", this.table)
		return (*s_InsertUpdateEnd)(this)
	}
	items := []string{}
	for _, m := range this.updateMembers {
		this.addInValues(m.value(vobj))
		items = append(items, fmt.Sprintf("%s=?", m.quote()))
	}
	this.sqlText += " " + strings.Join(items, ",")
	return (*s_InsertUpdateEnd)(this)
}

// -------------------------------------------------------------------
// end
// -------------------------------------------------------------------
type s_InsertUpdateEnd s_InsertUpdate

func (this *s_InsertUpdateEnd) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}
