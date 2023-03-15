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
	table   *S_Table    // 要更新的表
	members []*S_Member // 要更新的对象成员名
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
		table: table,
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
		table: table,
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
		table: table,
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
// InsertIgnore
// -------------------------------------------------------------------
func InsertIgnore(table *S_Table, mnames ...string) *s_InsertValues {
	this := Insert(table, mnames...)
	this.sqlText = strings.Replace(this.sqlText, "INSERT INTO", "INSERT IGNORE ", 1)
	return this
}

func InsertIgnoreAll(table *S_Table) *s_InsertValues {
	this := InsertAll(table)
	this.sqlText = strings.Replace(this.sqlText, "INSERT INTO", "INSERT IGNORE ", 1)
	return this
}

func InsertIgnoreBesides(table *S_Table, mnames ...string) *s_InsertValues {
	this := InsertBesides(table, mnames...)
	this.sqlText = strings.Replace(this.sqlText, "INSERT INTO", "INSERT IGNORE ", 1)
	return this
}

// -------------------------------------------------------------------
// Values/Objects
// -------------------------------------------------------------------
type s_InsertValues s_Insert

func (this *s_InsertValues) Values(values ...[]any) *s_InsertEnd {
	if this.notOK() {
		return (*s_InsertEnd)(this)
	}

	items := []string{}
	count := len(this.members)
	item := "(" + strings.Repeat(",?", count)[1:] + ")"
	for _, row := range values {
		if count != len(row) {
			this.errorf("the number of incoming values is not consistent with the number of columns(count=%d) to be inserted", count)
			return (*s_InsertEnd)(this)
		}
		this.addInValues(row...)
		items = append(items, item)
	}
	this.sqlText += " VALUES" + strings.Join(items, ",")
	return (*s_InsertEnd)(this)
}

func (this *s_InsertValues) Objects(objs ...any) *s_InsertEnd {
	if this.notOK() {
		return (*s_InsertEnd)(this)
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
	item := "(" + strings.Repeat(",?", len(this.members))[1:] + ")"
	for _, obj := range objs {
		if ok, err := buildItem(obj); err != nil {
			this.errorf(err.Error())
			return (*s_InsertEnd)(this)
		} else if ok {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		this.errorf("not a valid input object to update")
		return (*s_InsertEnd)(this)
	}
	this.sqlText += " VALUES" + strings.Join(items, ",")
	return (*s_InsertEnd)(this)
}

// -------------------------------------------------------------------
// End
// -------------------------------------------------------------------
type s_InsertEnd s_Insert

func (this *s_InsertEnd) End() *S_ExecInfo {
	return newExecInfo(this.createSQLInfo())
}
