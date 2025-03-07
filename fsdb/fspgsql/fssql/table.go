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
var _colOutTags = []string{"pgsqlout", "dbout"} // 列名传出映射 tag 前缀，譬如：SELECT count(0) FROM <table>，这里的 count(0) 可以用 "dbout" 修饰

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
	table      *S_Table              // 成员所属的表
	name       string                // 成员名称
	buildInfo  string                // 构建列描述
	dbkey      string                // 在数据库中对应的列关键字
	tbkey      string                // 所属表的key(从 *S_Table.Join 函数传入)
	dbout      string                // 列传出时的表达式
	field      reflect.StructField   // 成员所属结构体的域
	pathFields []reflect.StructField // 继承链的域列表
}

var memberType = reflect.TypeOf(new(S_Member)).Elem()

// ----------------------------------------------------------------------------
// private
// ----------------------------------------------------------------------------
func (this *S_Member) quote() string {
	return fmt.Sprintf("%q", this.dbkey)
}

// 带表名引用字段
func (this *S_Member) quoteWithTable() string {
	if this.tbkey == "" {
		return this.table.quote() + "." + this.quote()
	}
	var tbName = this.tbkey
	tb := this.table.joinTables[this.tbkey]
	if tb != nil {
		tbName = tb.name
	}
	return fmt.Sprintf("%q.%q", tbName, this.dbkey)
}

// 输出表达式中，通配符：
// ${n} ：表示字段名称
func (this *S_Member) outexp(name string) string {
	if this.dbout == "" {
		return name
	}
	return strings.ReplaceAll(this.dbout, "${n}", name)
}

func (this *S_Member) isList() bool {
	t := this.field.Type
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Array || t.Kind() == reflect.Slice
}

// vobj 必须是结构体对象的 reflect.Value (非 reflect.Ptr)
func (this *S_Member) value(obj any) (any, error) {
	var state int
	vobj := reflect.ValueOf(obj)
	for _, field := range this.pathFields {
		state, vobj = fsreflect.BaseRefValue(vobj)
		if state < 1 {
			return nil, fmt.Errorf("nil object or nil inner object")
		}
		vobj = vobj.FieldByName(field.Name)
	}
	state, vobj = fsreflect.BaseRefValue(vobj)
	if state < 1 {
		return nil, fmt.Errorf("nil object or nil inner object")
	}
	fobj := vobj.FieldByName(this.field.Name)
	if !fobj.IsValid() {
		return nil, fmt.Errorf("no filed name %q in object of type %v", this.field.Name, fobj.Type())
	}
	if !fobj.CanAddr() {
		return nil, fmt.Errorf("field %q in object of type %v is not accessable", this.field.Name, vobj.Type())
	}
	if fobj.CanInterface() {
		return fobj.Interface(), nil
	}
	return reflect.NewAt(this.field.Type, unsafe.Pointer(fobj.UnsafeAddr())).Elem().Interface(), nil
}

// vobj 必须是结构体对象的 reflect.Value (非 reflect.Ptr)
func (this *S_Member) valuePtr(obj any) (any, error) {
	state, vobj := fsreflect.BaseRefValue(reflect.ValueOf(obj))
	if state < 1 {
		return nil, fmt.Errorf("nil object")
	}
	for _, field := range this.pathFields {
		fobj := vobj.FieldByName(field.Name)
		state, fobj := fsreflect.BaseRefValue(fobj)
		if state < 0 {
			return nil, fmt.Errorf("object of %v has no field named %q", vobj.Type(), field.Name)
		}
		if state > 0 {
			// 只要继承的不是结构体指针，则会走这里
			vobj = fobj
			continue
		}
		// 进入这里意味着继承的一定是指针
		if fobj.CanSet() {
			// 因此，这里不需要判断 field.Type 是不是指针(一定是指针)
			fobj.Set(reflect.New(field.Type.Elem()))
		} else if !fobj.CanAddr() {
			return nil, fmt.Errorf("object of type %v can't be accessed", field.Type)
		} else {
			// 这里 filed.Type 也一定是指针
			p := reflect.NewAt(field.Type, unsafe.Pointer(fobj.UnsafeAddr()))
			p.Elem().Set(reflect.New(field.Type.Elem()))
		}
		_, fobj = fsreflect.BaseRefValue(fobj)
		vobj = fobj
	}
	state, vobj = fsreflect.BaseRefValue(vobj)
	if state < 1 {
		// 原则上永远也不会进入这里
		return nil, fmt.Errorf("nil object or nil inner object")
	}
	fv := vobj.FieldByName(this.field.Name)
	if !fv.IsValid() {
		return nil, fmt.Errorf("object of %v has no field named %q", vobj.Type(), this.field.Name)
	} else if !fv.CanAddr() {
		return nil, fmt.Errorf("field %q in object of %v can't be accessed", this.field.Name, vobj.Type())
	}
	return reflect.NewAt(this.field.Type, unsafe.Pointer(fv.UnsafeAddr())).Interface(), nil
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
func (this *S_Member) Name() string {
	return this.name
}

func (this *S_Member) String() string {
	return this.table.String() + "." + this.name
}

// 对应是数据库列名
func (this *S_Member) DBKey() string {
	return this.dbkey
}

// 所属表名
func (this *S_Member) TableName() string {
	if this.tbkey == "" {
		return this.table.name
	}
	tb := this.table.joinTables[this.tbkey]
	if tb == nil {
		return "<unknow>"
	}
	return tb.name
}

// 获取字段引用
func (this *S_Member) AutoQuote() string {
	if len(this.table.joinTables) > 0 {
		return this.quoteWithTable()
	}
	if this.tbkey == "" || this.table.name == "" {
		return this.quote()
	}
	return this.quoteWithTable()
}

// ---------------------------------------------------------------------------------------
// sql 语法中的表
// ---------------------------------------------------------------------------------------
// 缓存 table，以免每次都要重新解释数据库记录对象成员
var _tbcache = map[reflect.Type]*S_Table{}

type S_Table struct {
	name       string              // 表在 DB 中的名称
	dbkey      string              // 对应的数据库键
	tobj       reflect.Type        // 表对应的 go 对象反射类型
	dbFields   T_DBFields          // 数据库字段映射
	members    []*S_Member         // 成员列表
	joinTables map[string]*S_Table // 关联表
}

var tableType = reflect.TypeOf(new(S_Table)).Elem()

// 创建对象映射数据库表
func NewTable(name string, obj any) (*S_Table, error) {
	table := &S_Table{
		name:       name,
		dbkey:      `"` + name + `"`,
		members:    []*S_Member{},
		joinTables: map[string]*S_Table{},
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
		name:       name,
		dbkey:      `"` + name + `"`,
		dbFields:   dbFields,
		members:    []*S_Member{},
		joinTables: map[string]*S_Table{},
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
func (this *S_Table) getMemberDBKey(f reflect.StructField) (tbkey string, dbkey string) {
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
		return "", dbkey
	}
	fields := strings.Split(dbkey, ".")
	if len(fields) == 1 {
		return "", fields[0]
	} else if len(fields) > 1 {
		return fields[0], fields[1]
	}
	return "", "-"
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

func (this *S_Table) takeMembers(table *S_Table, obj any) {
	members := map[string]bool{}
	// 遍历顺序为，优先遍历子结构体成员
	fsreflect.TrivalStructMembers(obj, false, func(info *fsreflect.S_TrivalStructInfo) bool {
		if info.IsBase {
			return true
		}
		tbkey, dbkey := this.getMemberDBKey(info.Field)
		if dbkey == "-" {
			return true
		} // “-” 为，排除标记，表示该成员不映射到数据库
		fieldName := info.Field.Name
		// 如果继承的父结构里有同名成员，则取子的，忽略父的
		if members[fieldName] {
			return true
		}

		m := &S_Member{
			table:      table,
			name:       fieldName,
			buildInfo:  this.getDBBuildInfo(info.Field),
			dbkey:      dbkey,
			tbkey:      tbkey,
			dbout:      this.getMemberDBOut(info.Field),
			field:      info.Field,
			pathFields: info.PathFields,
		}
		members[m.name] = true
		this.members = append(this.members, m)
		return true
	})
}

func (this *S_Table) bindObjType(obj any) error {
	tobj := fsreflect.BaseRefType(reflect.TypeOf(obj))
	if tobj == nil {
		return errors.New("database table map object mustn't be a nil value")
	}
	if tobj.Kind() != reflect.Struct {
		return errors.New("database table map object must be a struct")
	}

	this.tobj = tobj
	this.members = make([]*S_Member, 0)
	this.takeMembers(this, obj)
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
func (this *S_Table) Join(key string, table *S_Table) *S_Table {
	this.joinTables[key] = table
	return this
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
// 通过数据库记录映射对象的成员名称，获取对应的 Member 对象
func (this *S_Table) Member(mname string) *S_Member {
	for _, member := range this.members {
		if member.name == mname {
			return member
		}
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
	for _, member := range this.members {
		if member.name == mname {
			return true
		}
	}
	return false
}

func (this *S_Table) CreateSQL(tails ...string) string {
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %q(%%s)", this.name)
	cols := []string{}
	for _, member := range this.members {
		col := fmt.Sprintf("%s %s", member.dbkey, member.buildInfo)
		cols = append(cols, col)
	}
	items := append(cols, tails...)
	return fmt.Sprintf(sql, strings.Join(items, ","))
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
