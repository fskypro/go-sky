/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: inserter
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

type s_Insert struct {
	s_SQL
	table        *S_Table    // 要更新的表
	members      []*S_Member // 要更新的对象成员名
	inValueOrder int         // 传入值序号
}

// -------------------------------------------------------------------
// Insert
// 可调用链：
//	 [Insert()/InsertAll()/InsertBesides()].Values().End()
//	 [Insert()/InsertAll()/InsertBesides()].Objects().End()
//
//	 [InsertIgnore()/InsertIgnoreAll()/InsertIgnoreBesides()].Values().End()
//	 [InsertIgnore()/InsertIgnoreAll()/InsertIgnoreBesides()].Objects().End()
//
// 注意：
//   End() 标示构建SQL语句结束
// -------------------------------------------------------------------
// 更新指定成员值
func Insert(table *S_Table, mnames ...string) *s_InsertValues {
	if len(mnames) == 0 {
		return InsertAll(table)
	}

	this := &s_Insert{
		table:        table,
		inValueOrder: 1,
	}
	this.sqlText = "INSERT INTO " + table.quote()
	dbkeys := []string{}
	for _, name := range mnames {
		m := this.table.Member(name)
		if m == nil {
			this.errorf("table %s has no object member named %q", this.table, name)
			return (*s_InsertValues)(this)
		}
		this.members = append(this.members, m)
		dbkeys = append(dbkeys, m.quote())
	}
	this.sqlText += "(" + strings.Join(dbkeys, ",") + ")"
	return (*s_InsertValues)(this)
}

// 更新所有成员值
func InsertAll(table *S_Table) *s_InsertValues {
	this := &s_Insert{
		table:        table,
		inValueOrder: 1,
	}
	this.sqlText = "INSERT INTO " + table.quote()
	this.members = this.table.orderMembers
	dbkeys := []string{}
	for _, m := range this.table.orderMembers {
		dbkeys = append(dbkeys, m.quote())
	}
	this.sqlText += "(" + strings.Join(dbkeys, ",") + ")"
	return (*s_InsertValues)(this)
}

// 更新除指定成员以外的其他所有成员值
func InsertBesides(table *S_Table, mnames ...string) *s_InsertValues {
	if len(mnames) == 0 {
		return InsertAll(table)
	}

	this := &s_Insert{
		table:        table,
		inValueOrder: 1,
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

		this.members = append(this.members, m)
		dbkeys = append(dbkeys, m.quote())
	}
	this.sqlText += "(" + strings.Join(dbkeys, ",") + ")"
	return (*s_InsertValues)(this)
}

// -------------------------------------------------------------------
// Values/Objects
// -------------------------------------------------------------------
type s_InsertValues s_Insert

func (this *s_InsertValues) Values(values ...[]any) *s_InsertConflict {
	if this.notOK() {
		return (*s_InsertConflict)(this)
	}

	items := []string{}
	count := len(this.members)
	item := ""
	for range this.members {
		item += fmt.Sprintf(",$%d", this.inValueOrder)
		this.inValueOrder++
	}
	item = "(" + item[1:] + ")"
	for _, row := range values {
		if count != len(row) {
			this.errorf("the number of incoming values is not consistent with the number of columns(count=%d) to be inserted", count)
			return (*s_InsertConflict)(this)
		}
		this.addInValues(row...)
		items = append(items, item)
	}
	this.sqlText += " VALUES" + strings.Join(items, ",")
	return (*s_InsertConflict)(this)
}

func (this *s_InsertValues) Objects(objs ...any) *s_InsertConflict {
	if this.notOK() {
		return (*s_InsertConflict)(this)
	}

	buildItem := func(o any) (bool, error) {
		vobj := reflect.ValueOf(o)
		tobj := reflect.TypeOf(o)
		if !vobj.IsValid() {
			return false, nil
		}
		if tobj.Kind() == reflect.Ptr {
			if vobj.IsNil() {
				return false, nil
			}
			vobj = vobj.Elem()
			tobj = tobj.Elem()
		}
		if tobj != this.table.tobj {
			return false, fmt.Errorf("the input object %#v type is not the same as the table's map object type(%s) to be inserted", o, this.table)
		}
		for _, m := range this.members {
			this.addInValues(m.value(vobj))
		}
		return true, nil
	}

	items := []string{}
	item := ""
	for range this.members {
		item += fmt.Sprintf(",$%d", this.inValueOrder)
		this.inValueOrder++
	}
	item = "(" + item[1:] + ")"
	for _, obj := range objs {
		if ok, err := buildItem(obj); err != nil {
			this.errorf(err.Error())
			return (*s_InsertConflict)(this)
		} else if ok {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		this.errorf("not a valid input object to update")
		return (*s_InsertConflict)(this)
	}
	this.sqlText += " VALUES" + strings.Join(items, ",")
	return (*s_InsertConflict)(this)
}

// -------------------------------------------------------------------
// conflict
// -------------------------------------------------------------------
type s_InsertConflict s_Insert

func (this *s_InsertConflict) OnConflict(mnames ...string) *s_InsertConflictDo {
	if this.notOK() { return (*s_InsertConflictDo)(this) }

	this.sqlText += " ON CONFLICT"
	if len(mnames) == 0 { return (*s_InsertConflictDo)(this) }
	dbkeys := []string{}
	for _, name := range mnames {
		m := this.table.Member(name)
		if m == nil {
			this.errorf("table %s has no object member named %q", this.table, name)
			return (*s_InsertConflictDo)(this)
		}
		dbkeys = append(dbkeys, m.quote())
	}
	this.sqlText += fmt.Sprintf("(%s)", strings.Join(dbkeys, ","))
	return (*s_InsertConflictDo)(this)
}

func (this *s_InsertConflict) Returning(mname string) *s_InsertEnd {
	if this.notOK() { return (*s_InsertEnd)(this) }
	member := this.table.Member(mname)
	if member == nil {
		this.errorf("table %s has no object member named %q", this.table, mname)
		return (*s_InsertEnd)(this)
	}
	this.sqlText += " RETURNING " + member.quote()
	return (*s_InsertEnd)(this)
}

func (this *s_InsertConflict) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}

// ---------------------------------------------------------
type s_InsertConflictDo s_Insert

// 存在则不做任何事
func (this *s_InsertConflictDo) DoNothing() *s_InsertReturning {
	if this.notOK() { return (*s_InsertReturning)(this) }

	this.sqlText += " DO NOTHING"
	return (*s_InsertReturning)(this)
}

// 以值更新表记录
func (this *s_InsertConflictDo) DoUpdateSet(values ...interface{}) *s_InsertReturning {
	if this.notOK() { return (*s_InsertReturning)(this) }

	if len(values) != len(this.members) {
		this.errorf("the number of update values must be %d", len(this.members))
		return (*s_InsertReturning)(this)
	}
	this.addInValues(values...)
	items := []string{}
	for _, m := range this.members {
		items = append(items, fmt.Sprintf("%s=$%d", m.quote(), this.inValueOrder))
		this.inValueOrder++
	}
	this.sqlText += " DO UPDATE SET " + strings.Join(items, ",")
	return (*s_InsertReturning)(this)
}

// 用表对应对象更新记录
func (this *s_InsertConflictDo) DoUpdateSetObject(obj interface{}) *s_InsertReturning {
	if this.notOK() { return (*s_InsertReturning)(this) }

	vobj := reflect.ValueOf(obj)
	tobj := reflect.TypeOf(obj)
	if tobj.Kind() == reflect.Ptr {
		vobj = vobj.Elem()
		tobj = tobj.Elem()
	}
	if tobj != this.table.tobj {
		this.errorf("input object type is not the same as the table %s binds object type", this.table)
		return (*s_InsertReturning)(this)
	}

	items := []string{}
	for _, m := range this.members {
		this.addInValues(m.value(vobj))
		items = append(items, fmt.Sprintf("%s=$%d", m.quote(), this.inValueOrder))
		this.inValueOrder++
	}
	this.sqlText += " DO UPDATE SET " + strings.Join(items, ",")
	return (*s_InsertReturning)(this)
}

func (this *s_InsertConflictDo) DoUpdateSetExp(exp string, args ...interface{}) *s_InsertReturning {
	if this.notOK() { return (*s_InsertReturning)(this) }

	exp = this.explainExp(this.table, exp, args...)
	if this.notOK() {
		this.errorf("error update exp, %v", this.err.Error())
		return (*s_InsertReturning)(this)
	}
	this.sqlText += " DO UPDATE SET " + exp
	return (*s_InsertReturning)(this)
}

func (this *s_InsertConflictDo) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}

// -------------------------------------------------------------------
// Returning
// -------------------------------------------------------------------
type s_InsertReturning s_Insert

func (this *s_InsertReturning) Returning(mname string) *s_InsertEnd {
	if this.notOK() { return (*s_InsertEnd)(this) }
	member := this.table.Member(mname)
	if member == nil {
		this.errorf("table %s has no object member named %q", this.table, mname)
		return (*s_InsertEnd)(this)
	}
	this.sqlText += " RETURNING " + member.quote()
	return (*s_InsertEnd)(this)
}

func (this *s_InsertReturning) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}

// -------------------------------------------------------------------
// End
// -------------------------------------------------------------------
type s_InsertEnd s_Insert

func (this *s_InsertEnd) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}
