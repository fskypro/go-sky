/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: sql builder base
@author: fanky
@version: 1.0
@date: 2022-01-17
**/

package fssql

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fsky.pro/fscollection"
	"fsky.pro/fsky"
	"fsky.pro/fsregexp"
	"fsky.pro/fsstr"
	"fsky.pro/fstype"
	"github.com/lib/pq"
)

// 允许使用的引导字符
const leads = "!@#$%^&*_-+?|:"

// group[1] ：引导字符前面还有个引导字符
// group[2] ：参数索引
// group[3] ：转义字符（s|q|v|o|TN|TM|TV|TO）
var reptns = []string{
	// 不需要指定 content 的正则模式
	`(%[1]s)?%[1]s(?:\[\s*(\d+)\s*\])?(s|q|v|TN)`,
	// 必须指定 conent 的正则模式（group[4]=""; group[5]="", group[6]="", group[7]=content）
	`(%[1]s)?%[1]s(?:\[\s*(\d+)\s*\])?(o)()()()\{(.+?)\}`,
	// content 内容可选的（group[4]=alias, group[5] = "-"；group[]="."; group[7]=content）
	`(%[1]s)?%[1]s(?:\[\s*(\d+)\s*\])?(TM|TV|TE|TO)(?:\(\s*(\w+)\s*\))?(-)?(\.)?(?:\{([_0-9A-z\., ]*)\})?`,
}

var bytea = reflect.TypeOf([]byte{})

// -------------------------------------------------------------------
// SQL
// -------------------------------------------------------------------
type S_SQL struct {
	SQLTxt  string // sql
	Inputs  []any  // 传入参数
	Outputs []any  // 传出参数
	Error   error  // 拼接SQL错误

	recombine *regexp.Regexp
	resubs    []*regexp.Regexp
	argOrder  int // 当前解释到到参数索引
	inOrder   int // $ 符序号
}

func newSQL(lead byte) *S_SQL {
	sql := &S_SQL{
		Inputs:   make([]any, 0),
		Outputs:  make([]any, 0),
		argOrder: 1,
		inOrder:  1,
	}
	relead := string(lead)
	if fsregexp.IsEscapeChar(lead) {
		relead = `\` + relead
	}

	ptns := []string{}
	for _, ptn := range reptns {
		ptn := fmt.Sprintf(ptn, relead)
		ptns = append(ptns, fmt.Sprintf("(?:%s)", ptn))
		sql.resubs = append(sql.resubs, regexp.MustCompile(ptn))
	}
	sql.recombine = regexp.MustCompile(strings.Join(ptns, "|"))
	return sql
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func (this *S_SQL) errorf(m string, args ...any) {
	this.Error = errors.Join(this.Error, fmt.Errorf(m, args...))
}

// 将大括号中的内容转换为数据对象成员名称列表
func (this *S_SQL) takeContentMembers(content string) []string {
	members := []string{}
	for _, m := range strings.Split(content, ",") {
		m = strings.TrimSpace(m)
		if m == "" { continue }
		members = append(members, m)
	}
	return members
}

// 添加传入值
func (this *S_SQL) addInput(value any) (string, error) {
	rt := reflect.TypeOf(value)
	if rt == nil { return "NULL", nil }
	// 如果传入 []byte 类型的值，则认为数据库中对应的字段类型为 bytea
	if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
		// 注意：
		//  如果 sql 语句中有如：
		//    SELECT * FROM <table> WHERE col IN (1,2,3);
		//  其中 (1,2,3) 是不定的传入 slice，则以上语句需要这样写：
		//    SELECT * FROM %TN WHERE col=ANY(%v);
		//  然后 []int{1,2,3} 作为 SQL(...) 函数的第二个参数传入
		//  同理 NOT IN 的话，就是 <>ANY(%v)
		value = pq.Array(value)
	} else if rt == reflect.TypeOf(time.Now()) {
		value = pq.FormatTimestamp(value.(time.Time))
	}
	this.Inputs = append(this.Inputs, value)
	ret := fmt.Sprintf("$%d", this.inOrder)
	this.inOrder++
	return ret, nil
}

// 引用成员名称
func (this *S_SQL) getMemberName(m *S_Member, tbAlias string, withTable bool) string {
	if tbAlias != "" {
		// 指定表别名
		return fmt.Sprintf("%s.%s", tbAlias, m.quote())
	}
	if withTable {
		return m.quoteWithTable()
	}
	return m.quote()
}

// 作为传出引用成员名称
func (this *S_SQL) getOutMemberName(m *S_Member, tbAlias string, withTable bool) string {
	if tbAlias != "" {
		// 指定表别名
		return fmt.Sprintf("%s.%s", tbAlias, m.quote())
	}
	if withTable {
		return m.quoteWithTable()
	}
	return m.outexp()
}

// ---------------------------------------------------------
// 将指定参数直接取其字符串放入 SQL 语句
// 格式化字符串：
//   %s 或 %[index]s
func (this *S_SQL) parse_s(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	return fmt.Sprintf("%v", arg)
}

// 将指定参数直接取其字符串并加上双引号放入 SQL 语句
// 格式化字符串：
//   %q 或 %[index]q
func (this *S_SQL) parse_q(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	return fmt.Sprintf(`"%v"`, arg)
}

// 格式化预传入值
// 格式化字符串：
//   %v 或 %[index]v
//     表示将指定参数作为 sql 的传入值
func (this *S_SQL) parse_v(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	subsql, err := this.addInput(arg)
	if err != nil {
		this.errorf("take value of argument %d for expression %q fail, %v", argOrder, exp, err)
		return exp
	}
	return subsql
}

// 设置传出参数
// 格式化字符串：
//   %o(mname) 或 %[index](mname)
//     表示将 table.col 传出给指定的参数
func (this *S_SQL) parse_o(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	if content == "" {
		this.errorf("no content in expression %q", exp)
		return exp
	}
	rarg := reflect.ValueOf(arg)
	if !rarg.IsValid() {
		this.errorf("argument %d must be a pointer for expression %q", argOrder, exp)
		return exp
	}
	if rarg.Type().Kind() != reflect.Ptr {
		this.errorf("argument %d must be a pointer for expression %q", argOrder, exp)
		return exp
	}
	this.Outputs = append(this.Outputs, arg)
	return content
}

// 根据表对象，传入表名
// 格式化字符串：
//   %TN 或 %[index]TN
//     将转义符转换为对应参数表格的表名，对应参数必须是 S_Table 指针对象
func (this *S_SQL) parse_TN(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	if !fstype.IsType[*S_Table](arg) {
		this.errorf("argument %d must be a pointer of %v for expression %q", argOrder, tableType, exp)
		return exp
	}
	if fsky.IsNil(arg) {
		this.errorf("argument %d must be a not nil pointer of %v for expression %q", argOrder, tableType, exp)
		return exp
	}
	return arg.(*S_Table).quote()
}

// 提取表记录对象成员所对应的数据库字段，并将这些字段排列成一串以逗号分隔开的字符串
// 格式化字符串：
//   1、%TM{} 或 %[index]{}
//     表示获取指定参数表记录对象，其所有成员对应的数据库字段名称，并用逗号分隔开，将其放到 SQL 语句中。
//   2、%TM{mname1, mname2, ...} 或 %TM[index]{mname1, mname2, ...}
//     表示获取指定参数表记录对象，其括号中指定的成员对应的数据库字段名称，并用逗号分隔开，将其放到 SQL 语句中。
//   3、%TM-{mname1, mname2, ...} 或 %TM[index]-{mname1, mname2, ...}
//     表示获取指定参数表记录对象，其除了括号中的成员以外的所有成员对应的数据库字段名称，并用逗号分隔开，将其放到 SQL 语句中。
//   4、以上三种方式下都可以通过在 %TM 后面加上一个括号括起来的别名，如：%TM(t){M1, M2}，这样在生成 sql 语句时，会字字段前面
//     用别名限定，如：t.m1, t.m2
//
//  注意：
//     对应参数必须是一个 S_Table 指针对象，或者是数据库表记录映射对象
func (this *S_SQL) parse_TM(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	members := []string{}
	for _, m := range strings.Split(content, ",") {
		m = strings.TrimSpace(m)
		if m == "" { continue }
		members = append(members, m)
	}
	var tb *S_Table
	if fstype.IsType[*S_Table](arg) {
		tb = arg.(*S_Table)
	} else {
		t, err := getObjTable(arg)
		if err != nil {
			this.errorf("argument %d must be a not nil pointer of %v for expression %q", argOrder, tableType, exp)
			return exp
		}
		tb = t
	}
	// 所有成员
	if len(members) == 0 {
		return fsstr.JoinFunc(tb.orderMembers, ",", func(m *S_Member) string {
			return this.getMemberName(m, alias, withTable)
		})
	}

	dbkeys := []string{}
	// 排除指定成员外的所有成员
	if exclude {
		for _, m := range tb.orderMembers {
			if !fscollection.SliceHas(members, m.name) {
				dbkeys = append(dbkeys, this.getMemberName(m, alias, withTable))
			}
		}
	} else {
		for _, m := range members {
			if member := tb.Member(m); member == nil {
				this.errorf(`object type %v of argument %d has no member named %q in exporession %q`, tb.tobj, argOrder, m, exp)
				return exp
			} else {
				dbkeys = append(dbkeys, this.getMemberName(member, alias, withTable))
			}
		}
	}
	return strings.Join(dbkeys, ",")
}

// 将传入对象的值作为构建 SQL 的传入值
// 格式化字符串：
//   1、%TV{} 或 %[index]TV{}
//     表示将对应参数对象的所有成员值作为 SQL 的传入值
//   2、%TV{mname1, mname2, ...} 或 %[index]TV{mname1, mname2, ...}
//     表示将对应参数对象的成员(括号中指定的成员)值作为 SQL 的传入值
//   3、%TV-{manme1, mname2, ...} 或 %[index]TV-{mname1, mname2, ...}
//     表示将对应参数对象的成员(括号中指定的成员)值作为 SQL 的传入值
//
// 注意：
//   对应参数必须为数据库记录映射对象
func (this *S_SQL) parse_TV(exp string, argOrder int, exclude bool, walias string, ithTable bool, content string, arg any) string {
	members := []string{}
	for _, m := range strings.Split(content, ",") {
		m = strings.TrimSpace(m)
		if m == "" { continue }
		members = append(members, m)
	}
	tb, err := getObjTable(arg)
	if err != nil {
		this.errorf("argument %d must a db map object for expression %q", argOrder, exp)
		return exp
	}

	dollers := []string{}
	addMember := func(m *S_Member) bool {
		value, err := m.value(arg)
		if err != nil {
			this.errorf("value of member %q error in argument %d, %v", m.name, argOrder, err)
		}
		subsql, err := this.addInput(value)
		if err != nil {
			this.errorf("value of member %q error in argument %d, %v", m.name, argOrder, err)
			return false
		}
		dollers = append(dollers, subsql)
		return true
	}

	if len(members) == 0 {
		for _, m := range tb.orderMembers {
			if !addMember(m) { return exp }
		}
	} else if exclude {
		for _, m := range tb.orderMembers {
			if fscollection.SliceHas(members, m.name) { continue }
			if !addMember(m)                          { return exp }
		}
	} else {
		for _, name := range members {
			m := tb.members[name]
			if m == nil {
				this.errorf("argument %d has no member named %q for expression %q", argOrder, name, exp)
				return exp
			}
			if !addMember(m) { return exp }
		}
	}
	return strings.Join(dollers, ",")
}

// 将传入对象的值作为构建 SQL 的设置值
// 格式化字符串：
//   1、%TE{} 或 %[index]TE{}
//     表示将对应参数对象的所有成员值作为 SQL 的传入值
//   2、%TE{mname1, mname2, ...} 或 %[index]TE{mname1, mname2, ...}
//     表示将对应参数对象的成员(大括号中指定的成员)值作为 SQL 的传入值
//   3、%TE-{manme1, mname2, ...} 或 %[index]TE-{mname1, mname2, ...}
//     表示将对应参数对象的成员(括号中指定的成员)值作为 SQL 的设置传入值
//
// 注意：
//   对应参数必须为数据库记录映射对象
func (this *S_SQL) parse_TE(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	members := []string{}
	for _, m := range strings.Split(content, ",") {
		m = strings.TrimSpace(m)
		if m == "" { continue }
		members = append(members, m)
	}
	tb, err := getObjTable(arg)
	if err != nil {
		this.errorf("argument %d must a db map object for expression %q", argOrder, exp)
		return exp
	}

	eqs := []string{}
	addMember := func(m *S_Member) bool {
		value, err := m.value(arg)
		if err != nil {
			this.errorf("value of member %q error in argument %d, %v", m.name, argOrder, err)
		}
		subsql, err := this.addInput(value)
		if err != nil {
			this.errorf("value of member %q error in argument %d, %v", m.name, argOrder, err)
			return false
		}
		eqs = append(eqs, fmt.Sprintf("%s=%s", this.getMemberName(m, alias, withTable), subsql))
		return true
	}

	if len(members) == 0 {
		for _, m := range tb.orderMembers {
			if !addMember(m) { return exp }
		}
	} else if exclude {
		for _, m := range tb.orderMembers {
			if fscollection.SliceHas(members, m.name) { continue }
			if !addMember(m)                          { return exp }
		}
	} else {
		for _, name := range members {
			m := tb.members[name]
			if m == nil {
				this.errorf("argument %d has no member named %q for expression %q", argOrder, name, exp)
				return exp
			}
			if !addMember(m) { return exp }
		}
	}
	return strings.Join(eqs, ",")
}

// 该表达式主要用于插入数据时，如果表中已经存在则更新语句，SET 后面的语句：ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name;
// 如：%TU{M1, M2} 则生成的 SQL 为：m1=EXCLUDED.m1, m2=EXCLUDED.m2
func (this *S_SQL) parse_TU(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	members := []string{}
	for _, m := range strings.Split(content, ",") {
		m = strings.TrimSpace(m)
		if m == "" { continue }
		members = append(members, m)
	}
	tb, err := getObjTable(arg)
	if err != nil {
		this.errorf("argument %d must a db map object for expression %q", argOrder, exp)
		return exp
	}

	eqs := []string{}
	addMember := func(m *S_Member) bool {
		col := this.getMemberName(m, alias, withTable)
		eqs = append(eqs, fmt.Sprintf("%s=EXCLUDED.%s", col, col))
		return true
	}

	if len(members) == 0 {
		for _, m := range tb.orderMembers {
			if !addMember(m) { return exp }
		}
	} else if exclude {
		for _, m := range tb.orderMembers {
			if fscollection.SliceHas(members, m.name) { continue }
			if !addMember(m)                          { return exp }
		}
	} else {
		for _, name := range members {
			m := tb.members[name]
			if m == nil {
				this.errorf("argument %d has no member named %q for expression %q", argOrder, name, exp)
				return exp
			}
			if !addMember(m) { return exp }
		}
	}
	return strings.Join(eqs, ",")
}

// 将传入对象的成员名称作为逗号分隔字符串，并将对象的成员制作存放到 SQL 的传出值列表中
// 格式化字符串：
//   1、%TO{} 或 %[index]TO{}
//     表示将对应参数对象的所有成员名称，以逗号分隔开拼接成一个 sql 字符串，并将所有成员指针放进 sql 对象的传出列表中
//   2、%TO{mname1, mname2, ...} 或 %[index]TO{mname1, mname2, ...}
//     表示将对应参数对象的大括号中指定的成员名称，以逗号分隔开拼接成一个 sql 字符串，并将对应的成员指针放进 sql 对象的传出列表中
//   3、%TO-{mname1, mname2, ...} 或 %[index]TO-{mname1, mname2}
//     表示将对应参数对象的除了大括号中指定的成员以外的所有成员名称，以逗号分隔开拼接成一个 sql 字符串，并将同样对应的成员指针放进 sql 对象的传出列表中
func (this *S_SQL) parse_TO(exp string, argOrder int, exclude bool, alias string, withTable bool, content string, arg any) string {
	members := []string{}
	for _, m := range strings.Split(content, ",") {
		m = strings.TrimSpace(m)
		if m == "" { continue }
		members = append(members, m)
	}
	tb, err := getObjTable(arg)
	if err != nil {
		this.errorf("argument %d must a db map object for expression %q", argOrder, exp)
		return exp
	}

	dbkeys := []string{}
	if len(members) == 0 {
		for _, m := range tb.orderMembers {
			pvalue, err := m.valuePtr(arg)
			if err != nil {
				this.errorf("parse output object for member %q fail, %v", m.name, err)
				return exp
			}
			this.Outputs = append(this.Outputs, pvalue)
			dbkeys = append(dbkeys, this.getOutMemberName(m, alias, withTable))
		}
	} else if exclude {
		for _, m := range tb.orderMembers {
			if fscollection.SliceHas(members, m.name) { continue }
			pvalue, err := m.valuePtr(arg)
			if err != nil {
				this.errorf("parse output object fro member %q fail, %v", m.name, err)
				return exp
			}
			this.Outputs = append(this.Outputs, pvalue)
			dbkeys = append(dbkeys, this.getOutMemberName(m, alias, withTable))
		}
	} else {
		for _, name := range members {
			m := tb.members[name]
			if m == nil {
				this.errorf("argument %d has no member named %q for expression %q", argOrder, name, exp)
				return exp
			}
			pvalue, err := m.valuePtr(arg)
			if err != nil {
				this.errorf("parse output object fro member %q fail, %v", m.name, err)
				return exp
			}
			this.Outputs = append(this.Outputs, pvalue)
			dbkeys = append(dbkeys, this.getOutMemberName(m, alias, withTable))
		}
	}
	return strings.Join(dbkeys, ",")
}

func (this *S_SQL) parseArg(exp string, args []any) string {
	for _, re := range this.resubs {
		group := re.FindStringSubmatch(exp)
		if len(group) == 0 { continue }

		if group[1] != "" { return group[0][1:] } // 引导字符开头的忽略
		// 参数索引
		order := this.argOrder
		if group[2] != "" {
			order, _ = strconv.Atoi(group[2])
		}
		// 参数不够
		if order > len(args) {
			this.errorf("arguments are not enough for %q", group[0])
			return exp
		}
		// 转义符
		esc := group[3]

		// 引用表的别名
		alias := ""
		if len(group) > 4 {
			alias = group[4]
		}

		// 成员列表前是否有 - 号
		exclude := len(group) > 4 && group[5] == "-"

		// 表字段前是否带表名
		withTable := len(group) > 5 && group[6] == "."

		// 括号内容
		content := ""
		if len(group) > 6 {
			content = strings.TrimSpace(group[7])
		}

		// 根据转义符替换为实际内容
		txt := map[string]func(string, int, bool, string, bool, string, any) string{
			"s":  this.parse_s,
			"q":  this.parse_q,
			"v":  this.parse_v,
			"o":  this.parse_o,
			"TN": this.parse_TN,
			"TM": this.parse_TM,
			"TE": this.parse_TE,
			"TU": this.parse_TU,
			"TV": this.parse_TV,
			"TO": this.parse_TO,
		}[esc](group[0], order, exclude, alias, withTable, content, args[order-1])
		this.argOrder++
		return txt
	}
	return ""
}

// -------------------------------------------------------------------
// SQL public
// -------------------------------------------------------------------
func (this *S_SQL) Fmt() string {
	if this.Error != nil { return this.Error.Error() }
	index := 0
	return fmt.Sprintf("%s\n\t%s", this.SQLTxt,
		fsstr.JoinFunc(this.Inputs, ", ", func(e any) string {
			index++
			te := reflect.TypeOf(e)
			if te != nil && te.Kind() == reflect.String {
				return fmt.Sprintf(`$%d="%v"`, index, e)
			}
			return fmt.Sprintf("$%d=%v", index, e)
		}))
}

func (this *S_SQL) SQL(tx string, args ...any) *S_SQL {
	this.argOrder = 1
	this.SQLTxt += this.recombine.ReplaceAllStringFunc(tx, func(x string) string {
		if this.Error != nil {
			return x
		}
		return this.parseArg(x, args)
	})
	return this
}

// -----------------------------------------------------------------------------
// public
// -----------------------------------------------------------------------------
// 构建 sql 语句
// tx 可以包含以下转义符：
//  1、%[参数索引(可选)]s
//      将指定参数直接填补到 sql 语句的 %s 表达式中
//  2、%[参数索引(可选)]q
//      将指定参数直接填补到 sql 语句的 %s 表达式中，并在参数两边加上双引号
//
//  3、%[参数索引(可选)]v
//      普通 sql 传入值参数
//  4、%[参数索引(可选)]o{<sql表达式>}
//      普通传出参数（注意，必须是指针类型）
//
//  5、%[参数索引(可选)]TN
//      S_Table 在数据库中的名称
//
//  6、%[参数索引(可选)]TM{} 或 %[参数索引(可选)]TM.{}
//      指定 S_Table 的所有成员对应的数据库字段名称，以逗号分隔（对应位置的参数必须是 S_Table 对象）
//  7、%[参数索引(可选)]TM{成员名称1, 成员名称2, ...} 或 %[参数索引(可选)]TM.{成员名称1, 成员名称2, ...}
//      指定 S_Table 的成员对应的数据库字段名称，以逗号分隔（对应位置的参数必须是 S_Table 对象）
//  8、%[参数索引(可选)]TM-{成员名称1, 成员名称2, ...} 或 %[参数索引(可选)]TM-.{成员名称1, 成员名称2, ...}
//      指定 S_Table 排除掉大括号中指定的成员后，剩余成员对应的数据库字段名称，以逗号分隔（对应位置的参数必须是 S_Table 对象）
//
//  9、%[参数索引(可选)]TE{} 或 %[参数索引(可选)]TE.{}
//      将数据库映射对象的所有成员对应的数据库字段组合成 <dbkey1>=$n1,<dbkey2>=$n2 的形式
//  10、%[参数索引(可选)]TE{成员名称1, 成员名称2, ...} 或 %[参数索引(可选)]TE.{成员名称1, 成员名称2, ...}
//      将数据库映射对象，在大括号中列出的成员对应的数据库字段组合成 <dbkey1>=$n1,<dbkey2>=$n2 的形式
//  11、%[参数索引(可选)]TE-{成员名称1, 成员名称2, ...} 或 %[参数索引(可选)]TE-.{成员名称1, 成员名称2, ...}
//      将数据库映射对象，除了在大括号中列出的成员以外的所有成员对应的数据库字段组合成 <dbkey1>=$n1,<dbkey2>=$n2 的形式
//
//  12、%[参数索引(可选)]TV{}
//      将指定对象的成员值作为构建 SQL 语句的传入值，包括所有成员值（对应位置的参数必须是一个数据库表记录映射对象）
//  13、%[参数索引(可选)]TV{成员名称1, 成员名称2, ...}
//      将指定对象的成员值作为构建 SQL 语句的传入值，包括大括号中列出的成员值（对应位置的参数必须是一个数据库表记录映射对象）
//  14、%[参数索引(可选)]TV-{成员名称1, 成员名称2, ...}
//      将指定对象的成员值作为构建 SQL 语句的传入值，包括括号中指定成员以外的所有成员的值（对应位置的参数必须是一个数据库表记录映射对象）
//
//  15、%[参数索引(可选)]TO{}
//      指定对象的成员指针，作为构建 SQL 的传出参数，包括对象的所有成员（对应位置的参数必须是一个数据库表记录映射对象）
//  16、%[参数索引(可选)]TO{成员名称1, 成员名称2, ...}
//      指定对象的指定成员指针，作为构建 SQL 的传出参数，指定对象的成员名称（对应位置的参数必须是一个数据库表记录映射对象）
//  17、%[参数索引(可选)]TO-{成员名称1, 成员名称2, ...}
//      指定对象除了指出成员以外的所有成员指针，作为构建 SQL 的传出参数，包含除了大括号指定的成员名称以外的所有成员（对应位置的参数必须是一个数据库表记录映射对象）
//
//  18、%[参数索引(可选)]TU{}
//      将对象所有成员构建成以下等式列表：m1=EXCLUDED.m1, m2=EXCLUDED.m2
//      如：%TU{}，假设对象有两个成员，对应字段为 m1、m2，则生成的 SQL 子句为：m1=EXCLUDED.m1, m2=EXCLUDED.m2
//  19、%[参数索引(可选)]TU{成员名称1, 成员名称2}
//      将大括号中指定的对象名称所以对应的数据字段，生成以下形式的 SQL 子句：m1=EXCLUDED.m1, m2=EXCLUDED.m2
//      如：%TU{M1, M2} 则生成的 SQL 为：m1=EXCLUDED.m1, m2=EXCLUDED.m2
//  20、%[参数索引(可选)]TUi-{成员名称1, 成员名称2}
//      将除了大括号中指定成员名称以外的成员，其对应的数据库字段，生成以下格式的 SQL 子句：m1=EXCLUDED.m1, m2=EXCLUDED.m2
//      如：假设对象有四个成员 M1/M2/M3/M4，则，%TU-{M1, M2} 则生成的 SQL 子句为：m2=EXCLUDED.m2, m3=EXCLUDED.m3
//
//  注意：
//      大括号前面如果有 “.” 号，则表示构建 SQL 语句中，对于表字段的引用，前面加上表名。例如，假设表 table 中有字段 col，则如果加上点去引用，则 SQL 语句类似如下：
//        SELECT "table"."col" FROM "table"
//      如果不加点，则类似如下：
//        SELECT "col" FROM "table"
// ---------------------------------------
// 默认引导字符为 %
func SQL(tx string, args ...any) *S_SQL {
	return newSQL('%').SQL(tx, args...)
}

// 自定义引导字符，并构建 SQL 语句
// 引导字符只能是下面字符之一：!@#$%^&*_-+?|:
func LeadSQL(lead byte) *S_SQL {
	if !strings.Contains(leads, string(lead)) {
		return &S_SQL{Error: fmt.Errorf("sql lead char must be one of %q", leads)}
	}
	return newSQL(lead)
}
