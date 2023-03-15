/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: selector
@author: fanky
@version: 1.0
@date: 2022-01-01
**/

package fssql

import (
	"fmt"
	"strings"

	. "fsky.pro/fsky"
)

// -----------------------------------------------------------------------------
// Select
// 可调用链：
//   操作单个表：
//	   [Select()/SelectAll()/SelectBesides()].From().Where().View().End()
//	   [Select()/SelectAll()/SelectBesides()].From().Where().End()
//     [Select()/SelectAll()/SelectBesides()].From().View().End()
//	   [Select()/SelectAll()/SelectBesides()].From().End()
//
//   对单个表自定义查询表达式：
//	   SelectExp().From().Where().View().End()
//	   SelectExp().From().Where().End()
//	   SelectExp().From().View().End()
//	   SelectExp().From().End()
//
//   自定义查询表达式
//     SelectExp().FromExp().Where().View().End()
//     SelectExp().FromExp().Where().End()
//     SelectExp().FromExp().View().End()
//     SelectExp().FromExp().End()
//
//   查询关联表的方法：
//	   1、创建映射表对象
//		  type S_Table1 struct{
//			Key1   int    `db:"key"`
//		    Value1 string `db:"value"`
//		  }
//		  tb1, _ := fssql.NewTable("table1", new(S_Table1))
//
//		  type S_Table2 struct{
//			Key2   int    `db:"key"`
//		    Value2 string `db:"value"`
//		  }
//		  tb2, _ := fssql.NewTable("table2", new(S_Table2))
//
//	  2、创建关联表对象，下面相当于如下SQL语句：
//		  SELECT table1.key, table1.value, table2.value FROM table1 JOIN table2 ON table1.key=table2.key
//
//        type S_T1_Join_T2 struct {
//			Key    int    `db:"table1.key`
//			Value1 string `db:"table1.value"`
//			Value2 string `db:"table2.value"`
//        }
//		  tb1_tb2, _ := fssql.NewLinkTable(new(S_T1_Join_T2), "#[1] JOIN #[2] ON $[3]=$[4]", tb1, tb2, tb1.M("Key1"), tb2.M("Key2"))
//
//	 3、构建查询语句，如：
//		  SELECT table1.key, table1.value, table2.value FROM table1 JOIN table2 ON table1.key=table2.key WHERE table1.key IN (100, 200, 300)
//
//		  调用方法为：
//			fssql.SelectAll().From(tb1_tb2).Where("$[1] in ?[2]", table1.M("Key"), []int{100, 200, 300})
//		  或：
//			fssql.SelectAll().From(tb1_tb2).Where("$[1] in ?[2]", "Key", []int{100, 200, 300})
//			上面 $[1] 引用的 “Key” 只传入一个字符串，那么该字符串，在构建的 SQL 语句中，将会引用查询表（即 tb1_tb2 ）的 “Key” 成员
//			所对应的 “db” 指定的 tag，即: table1.key
//
// 注意：
//   1、End() 表示查询语句构建结束
//	 2、查询表达式中：
//		$[1] 表示对表对象成员的引用，即引用参数列表中，第二个索引参数，并且指明该参数是标对象成员；
//		?[2] 表示将要构建的 SQL 语句中，有传入参数，并且该参数是调用函数参数列表中的第三个参数
// -----------------------------------------------------------------------------
type s_Select struct {
	s_SQL
	selMNames  []string // 要查询的成员名称表达式
	selBesides []string // 被排除 select 的成员
	selExp     string   // 查询语句
	selExpArgs []any    // 查询表达式参数列表

	selTable *S_Table    // 传出对象所属的 table
	members  []*S_Member // 要查询的对象成员

	whered bool // 是否写了条件关机键
}

func Select(mnames ...string) *s_SelectFrom {
	if len(mnames) == 0 {
		return SelectAll()
	}
	this := &s_Select{
		members: make([]*S_Member, 0),
	}
	this.selMNames = mnames
	return (*s_SelectFrom)(this)
}

func SelectAll() *s_SelectFrom {
	this := &s_Select{
		members: make([]*S_Member, 0),
	}
	this.selExp = ""
	return (*s_SelectFrom)(this)
}

// select 除指定成员以外，的其他所有成员，调用该函数后，必须接着调用 ToObj 方法
func SelectBesides(mnames ...string) *s_SelectFrom {
	this := &s_Select{
		members: make([]*S_Member, 0),
	}
	this.selBesides = mnames
	return (*s_SelectFrom)(this)
}

// 如果 select Max(money) from `order`：
//	SelectExp("max(%s)", "Money")
// 或：
//	SelectExp("max[%[1]s]", order.M("Money"))
func SelectExp(exp string, args ...interface{}) *s_SelectFrom {
	this := &s_Select{
		members: make([]*S_Member, 0),
	}
	this.selExp = exp
	this.selExpArgs = args
	return (*s_SelectFrom)(this)
}

// -----------------------------------------------------------------------------
// From
// -----------------------------------------------------------------------------
type s_SelectFrom s_Select

func (this *s_SelectFrom) From(table *S_Table) *S_SelectWhere {
	this.selTable = table
	if this.selExp != "" {
		exp := this.explainExp(table, this.selExp, this.selExpArgs...)
		if this.notOK() {
			this.errorf("explain select expression %q fail, %v", this.selExp, this.err)
		} else {
			this.sqlText = fmt.Sprintf("SELECT %s FROM %s", exp, table.quote())
		}
		return (*S_SelectWhere)(this)
	}

	dbkeys := []string{}
	if this.selMNames != nil {
		for _, name := range this.selMNames {
			m := table.Member(name)
			if m == nil {
				this.errorf("table %s has no member named %q", table, name)
				return (*S_SelectWhere)(this)
			}
			this.members = append(this.members, m)
			dbkeys = append(dbkeys, IfElse(table.IsLink(), m.quoteWithTable(), m.quote()))
		}
	} else if this.selBesides != nil {
	L:
		for _, m := range table.members {
			for _, name := range this.selBesides {
				if name == m.name {
					continue L
				}
			}
			this.members = append(this.members, m)
			dbkeys = append(dbkeys, IfElse(table.IsLink(), m.quoteWithTable(), m.quote()))
		}
	} else {
		for _, m := range table.members {
			this.members = append(this.members, m)
			dbkeys = append(dbkeys, IfElse(table.IsLink(), m.quoteWithTable(), m.quote()))
		}
	}
	this.sqlText = fmt.Sprintf("SELECT %s FROM %s", strings.Join(dbkeys, ","), table.quote())
	return (*S_SelectWhere)(this)
}

// -----------------------------------------------------------------------------
// Where
// -----------------------------------------------------------------------------
type S_SelectWhere s_Select

func (this *S_SelectWhere) where(exp string, args ...interface{}) (string, bool) {
	if this.notOK() {
		return "", false
	}
	exp = this.explainExp(this.selTable, exp, args...)
	if this.notOK() {
		this.errorf(`error "where" expression %q, %v`, exp, this.err)
		return "", false
	}
	return exp, true
}

func (this *S_SelectWhere) concat(link string, exp string) *S_SelectWhere {
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
func (this *S_SelectWhere) Quote() *S_SelectWhere {
	return this.concat("", "(")
}

// 与前括号
func (this *S_SelectWhere) AndQuote() *S_SelectWhere {
	return this.concat("AND", "(")
}

// 或前括号
func (this *S_SelectWhere) OrQuote() *S_SelectWhere {
	return this.concat("OR", "(")
}

// 后括号
func (this *S_SelectWhere) RQuote() *S_SelectWhere {
	this.sqlText += ")"
	return this
}

// 传入条件子句
// 如果参数为 *S_Member 则 exp 的转义符为 #[参数索引]，如
//   select * from tb where ID=100，则 Where 的调用方式是：
//	   Where("$[1]=?[2]", tb1.M("ID"), 100)
// 注意：
//   传入参数支持 slice 和数组，如：
//     select * from tb where ID in (1,2,3)，则 on 的调用方式是：
//	      Where("$[1] in ?[2]", tb1.M("ID"), []int{1,2,3})
func (this *S_SelectWhere) Where(exp string, args ...interface{}) *S_SelectWhere {
	return this.AndWhere(exp, args...)
}

func (this *S_SelectWhere) AndWhere(exp string, args ...interface{}) *S_SelectWhere {
	exp, ok := this.where(exp, args...)
	if !ok {
		return this
	}
	return this.concat("AND", exp)
}

func (this *S_SelectWhere) OrWhere(exp string, args ...interface{}) *S_SelectWhere {
	exp, ok := this.where(exp, args...)
	if !ok {
		return this
	}
	return this.concat("OR", exp)
}

func (this *S_SelectWhere) EmptyView() *S_SelectView {
	return (*S_SelectView)(this)
}

func (this *S_SelectWhere) View(exp string, args ...interface{}) *S_SelectView {
	return (*S_SelectView)(this).View(exp, args...)
}

func (this *S_SelectWhere) End() *S_SelectInfo {
	return (*s_SelectEnd)(this).End()
}

// -----------------------------------------------------------------------------
// View
// -----------------------------------------------------------------------------
type S_SelectView s_Select

// 输出控制子句
// 如果参数为 *S_Member 则 exp 的转义符为 #[参数索引]，如：
//   select * from tb where ID=100 group by ID，则 View 的调用方式是：
//	   View("group by $[1]", tb1.M("ID"))
// 又如：
//   select * from tb limit 100, 则 View 的调用方式：
//     View("limit ?[1]", 100)
func (this *S_SelectView) View(exp string, args ...interface{}) *S_SelectView {
	if this.notOK() {
		return this
	}
	exp = this.explainExp(this.selTable, exp, args...)
	if this.notOK() {
		this.errorf(`error view expression %q, %v`, exp, this.err)
	}
	this.sqlText += " " + exp
	return this
}

func (this *S_SelectView) End() *S_SelectInfo {
	return (*s_SelectEnd)(this).End()
}

// -----------------------------------------------------------------------------
// End
// -----------------------------------------------------------------------------
type s_SelectEnd s_Select

// 如果只需要构建一个 SQL 语句，则调用 End
func (this *s_SelectEnd) End() *S_SelectInfo {
	return newSelectInfo(this.createSQLInfo(), this.selTable, this.members)
}
