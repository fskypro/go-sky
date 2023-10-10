/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: object table wrapper
@author: fanky
@version: 1.0
@date: 2023-01-02
**/

package fssql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"fsky.pro/fsreflect"
)

var _colBuildTags = []string{"pgsqltd", "dbtd"} // 列名定义 tag 前缀
var _colTags = []string{"pgsql", "db"}          // 列名映射 tag 前缀
var _colOutTags = []string{"pgsqlout", "dbout"} // 列名传出映射 tag 前缀

// ---------------------------------------------------------------------------------------
// scheme
// ---------------------------------------------------------------------------------------
type S_DBField struct {
	DB   string // 字段名称
	DBTD string // 字段类型和描述
}

// {成员名称: *S_DBField}
type T_DBFields map[string]*S_DBField

// ---------------------------------------------------------------------------------------
// sql 语法中的成员
// ---------------------------------------------------------------------------------------
type S_Member struct {
	table     *S_Table     // 成员所属的表
	name      string       // 成员名称
	buildInfo string       // 构建列描述
	dbkey     string       // 在数据库中对应的列关键字
	dbout     string       // 列传出时的表达式
	vtype     reflect.Type // 类型
	offset    uintptr      // 字段偏移值
}

var memberType = reflect.TypeOf(new(S_Member)).Elem()

// ----------------------------------------------------------------------------
// private
// ----------------------------------------------------------------------------
func (this *S_Member) quote() string {
	return this.dbkey
}

func (this *S_Member) outexp() string {
	if this.dbout != "" {
		return this.dbout
	}
	return this.quote()
}

func (this *S_Member) quoteWithTable() string {
	return this.table.quote() + "." + this.quote()
}

// vobj 必须是结构体对象的 reflect.Value (非 reflect.Ptr)
func (this *S_Member) value(vobj reflect.Value) interface{} {
	p := unsafe.Pointer(this.offset + vobj.UnsafeAddr())
	return reflect.NewAt(this.vtype, p).Elem().Interface()
}

// vobj 必须是结构体对象的 reflect.Value (非 reflect.Ptr)
func (this *S_Member) valuePtr(vobj reflect.Value) interface{} {
	p := unsafe.Pointer(this.offset + vobj.UnsafeAddr())
	return reflect.NewAt(this.vtype, p).Interface()
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
func (this *S_Member) String() string {
	return this.table.String() + "." + this.name
}

// 对应是数据库列名
func (this *S_Member) DBKey() string {
	return this.dbkey
}

// ---------------------------------------------------------------------------------------
// sql 语法中的表
// ---------------------------------------------------------------------------------------
// 缓存 table，以免每次都要重新解释数据库记录对象成员
var _tbcache = map[reflect.Type]*S_Table{}

type S_Table struct {
	name         string               // 表在 DB 中的名称
	dbkey        string               // 对应的数据库键
	tobj         reflect.Type         // 表对应的 go 对象反射类型
	dbFields     T_DBFields           // 数据库字段映射
	members      map[string]*S_Member // 数据库记录对象成员名称对应的 go 结构成员
	orderMembers []*S_Member          // 有序成员列表
}

var tableType = reflect.TypeOf(new(S_Table)).Elem()

// 创建对象映射数据库表
func NewTable(name string, obj any) (*S_Table, error) {
	table := &S_Table{
		name:         name,
		dbkey:        `"` + name + `"`,
		members:      map[string]*S_Member{},
		orderMembers: []*S_Member{},
	}
	err := table.bindObjType(obj)
	if err != nil {
		return nil, fmt.Errorf("create table fail, %v", err)
	}
	return table, nil
}

// 创建对象映射数据库表，不使用结构体 tag，单独指定字段映射
func NewTableWithScheme(name string, obj any, dbFields T_DBFields) (*S_Table, error) {
	table := &S_Table{
		name:         name,
		dbkey:        `"` + name + `"`,
		dbFields:     dbFields,
		members:      map[string]*S_Member{},
		orderMembers: []*S_Member{},
	}
	err := table.bindObjType(obj)
	if err != nil {
		return nil, fmt.Errorf("create table fail, %v", err)
	}
	return table, nil
}

// -----------------------------------------------------------------------------
// private
// -----------------------------------------------------------------------------
func (this *S_Table) getMemberDBKey(f reflect.StructField) string {
	dbkey := ""
	if this.dbFields != nil {
		if dbf, ok := this.dbFields[f.Name]; ok {
			dbkey = dbf.DB
			goto L
		}
	}
	for _, tag := range _colTags {
		key := f.Tag.Get(tag)
		if key != "" {
			dbkey = key
			goto L
		}
	}
	if dbkey == "" {
		dbkey = f.Name
	}
L:
	if dbkey == "-" {
		return dbkey
	}
	dbkey = strings.ReplaceAll(dbkey, `"`, `"`)
	dbkey = strings.ReplaceAll(dbkey, ".", `"."`)
	return `"` + dbkey + `"`
}

func (this *S_Table) getMemberDBOut(f reflect.StructField) string {
	for _, tag := range _colOutTags {
		dbout := f.Tag.Get(tag)
		if dbout != "" {
			return dbout
		}
	}
	return ""
}

func (this *S_Table) getDBBuildInfo(f reflect.StructField) string {
	if this.dbFields != nil {
		if dbf, ok := this.dbFields[f.Name]; ok {
			return dbf.DBTD
		}
	}
	for _, tag := range _colBuildTags {
		buildInfo := f.Tag.Get(tag)
		if buildInfo != "" {
			return buildInfo
		}
	}
	return ""
}

func (this *S_Table) takeMembers(table *S_Table, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		tfield := t.Field(i)
		if tfield.Anonymous { // 匿名结构
			continue
		}

		dbkey := this.getMemberDBKey(tfield)
		if dbkey == "-" { // “-” 为，排除标记，表示该成员不映射到数据库
			continue
		}

		buildInfo := this.getDBBuildInfo(tfield)

		m := &S_Member{
			table:     table,
			name:      tfield.Name,
			buildInfo: buildInfo,
			dbkey:     dbkey,
			dbout:     this.getMemberDBOut(tfield),
			vtype:     tfield.Type,
			offset:    tfield.Offset,
		}

		this.members[m.name] = m
		this.orderMembers = append(this.orderMembers, m)
	}
}

func (this *S_Table) bindObjType(obj any) error {
	tobj, err := fsreflect.BaseRefType(obj)
	if err != nil {
		return errors.New("database table map object mustn't be a nil value")
	}
	if tobj.Kind() != reflect.Struct {
		return errors.New("database table map object must be a struct")
	}

	this.tobj = tobj
	this.members = make(map[string]*S_Member)
	this.orderMembers = make([]*S_Member, 0)
	this.takeMembers(this, tobj)
	if len(this.members) == 0 {
		return fmt.Errorf(`object "%v" has no map db members, may be is not a db map object`, tobj)
	}
	_tbcache[this.tobj] = this
	return nil
}

func (this *S_Table) quote() string {
	return this.dbkey
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
func (this *S_Table) Name() string {
	return this.name
}

func (this *S_Table) String() string {
	return fmt.Sprintf("%s[db:%s]", this.ObjectPath(), this.name)
}

// table 对应的 go 结构体所在包路径
func (this *S_Table) ObjectPath() string {
	return this.tobj.PkgPath() + ":" + this.tobj.String()
}

// -------------------------------------------------------------------
// 通过数据库记录映射对象的成员名称，获取对应的 Member 对象
func (this *S_Table) Member(mname string) *S_Member {
	if m, ok := this.members[mname]; ok {
		return m
	}
	return nil
}

func (this *S_Table) M(mname string) *S_Member {
	return this.Member(mname)
}

func (this *S_Table) HasMember(m *S_Member) bool {
	return m.table.tobj == this.tobj
}

func (this *S_Table) HasMemberName(mname string) bool {
	return this.members[mname] != nil
}

// -------------------------------------------------------------------
func (this *S_Table) CreateObject() interface{} {
	return reflect.New(this.tobj).Interface()
}

// -----------------------------------------------------------------------------
// module functions
// -----------------------------------------------------------------------------
func getObjTable(obj any) (*S_Table, error) {
	if obj == nil {
		return nil, errors.New("the db table's maps go object mustn't be a nil value")
	}
	tobj := reflect.TypeOf(obj)
	if tobj.Kind() == reflect.Ptr {
		tobj = tobj.Elem()
	}
	tb := _tbcache[tobj]
	if tb == nil {
		return NewTable("", obj)
	}
	return tb, nil
}
