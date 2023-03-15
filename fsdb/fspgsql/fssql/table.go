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
)

var _colBuildTags = []string{"pgsqltd", "dbtd"} // 列名定义 tag 前缀
var _colTags = []string{"pgsql", "db"}          // 列名映射 tag 前缀

// 提取对象中映射到 DB 的成员
func _iterMembers(tobj reflect.Type, fun func(*S_Member)) {
}

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

func (this *S_Member) quoteWithTable() string {
	if this.table.IsLink() {
		return this.quote()
	}
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
type S_Table struct {
	name         string               // 表在 DB 中的名称
	dbkey        string               // 对应的数据库键
	tobj         reflect.Type         // 表对应的 go 对象反射类型
	dbFields     T_DBFields           // 数据库字段映射
	members      map[string]*S_Member // 表字段对应的 go 结构成员
	orderMembers []*S_Member          // 有序成员列表

	// 映射表成员
	schemes   []string
	buildInfo string

	// 连接表成员
	link         string     // 连接表达式
	linkInValues []any      // 连接表达式中包含的传入参数
	linkTables   []*S_Table // 连接表
}

var tableType = reflect.TypeOf(new(S_Table)).Elem()

// 创建对象映射数据库表
func NewTable(name string, obj any) (*S_Table, error) {
	table := &S_Table{
		name:         name,
		dbkey:        `"` + name + `"`,
		members:      map[string]*S_Member{},
		orderMembers: []*S_Member{},

		schemes:   []string{},
		buildInfo: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",

		linkInValues: []any{},
		linkTables:   []*S_Table{},
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

		schemes:   []string{},
		buildInfo: "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",

		linkInValues: []any{},
		linkTables:   []*S_Table{},
	}
	err := table.bindObjType(obj)
	if err != nil {
		return nil, fmt.Errorf("create table fail, %v", err)
	}
	return table, nil
}

// 创建连接表
// link 是多表连接表达式，如：
//
//	newLinkTable("#[1] join #[2] ON $[3]=$[4]", tb1, tb2, tb1.M("ID"), tb2.M("ID"))
func NewLinkTable(obj interface{}, link string, args ...any) (*S_Table, error) {
	link, inValues, err := explainExp(nil, link, args...)
	if err != nil {
		return nil, fmt.Errorf("create LinkTable fail, %v", err)
	}
	// 找出所有连接表
	linkTables := []*S_Table{}
	for _, v := range args {
		switch v.(type) {
		case *S_Table:
			linkTables = append(linkTables, v.(*S_Table))
		}
	}

	tbName := fmt.Sprintf("LinkTable[%s]", link)
	table := &S_Table{
		name:         tbName,
		dbkey:        link,
		members:      map[string]*S_Member{},
		orderMembers: []*S_Member{},

		schemes: []string{},

		linkInValues: inValues,
		linkTables:   linkTables,
	}
	err = table.bindObjType(obj)
	if err != nil {
		return nil, fmt.Errorf("create table fail, %v", err)
	}
	return table, nil
}

// -----------------------------------------------------------------------------
// private
// -----------------------------------------------------------------------------
// 根据 go 对象类型，找出对应的连接表
func (this *S_Table) getLinkTable(tobj reflect.Type) *S_Table {
	for _, tb := range this.linkTables {
		if tb.tobj == tobj {
			return tb
		}
	}
	return this
}

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
			atype := tfield.Type
			if atype.Kind() == reflect.Ptr {
				atype = atype.Elem()
			}
			if atype.Kind() == reflect.Struct {
				this.takeMembers(table.getLinkTable(atype), atype)
			}
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
			vtype:     tfield.Type,
			offset:    tfield.Offset,
		}

		this.members[m.name] = m
		this.orderMembers = append(this.orderMembers, m)
	}
}

func (this *S_Table) bindObjType(obj interface{}) error {
	if obj == nil {
		return errors.New("the db table's maps go object mustn't be a nil value")
	}
	tobj := reflect.TypeOf(obj)
	if tobj.Kind() == reflect.Ptr {
		tobj = tobj.Elem()
	}
	this.tobj = tobj
	this.members = make(map[string]*S_Member)
	this.orderMembers = make([]*S_Member, 0)
	this.takeMembers(this, tobj)
	return nil
}

func (this *S_Table) quote() string {
	return this.dbkey
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
func (this *S_Table) IsLink() bool {
	return len(this.linkTables) > 0
}

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
// create table settings
// -----------------------------------------------------------------------------
func (this *S_Table) AddSchemes(sc ...string) {
	this.schemes = append(this.schemes, sc...)
}

func (this *S_Table) SetBuildInfo(info string) {
	this.buildInfo = info
}
