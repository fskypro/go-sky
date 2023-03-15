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
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	. "fsky.pro/fsky"

	"fsky.pro/fsky"
)

var _regexp = regexp.MustCompile(`[#\$\?]\[\d+\]`)

func getTableMemberDBKey(table *S_Table, isLink bool, mname string) string {
	member := table.Member(mname)
	if member != nil {
		return IfElse(isLink, member.quoteWithTable(), member.quote())
	}
	for _, tb := range table.linkTables {
		if dbkey := getTableMemberDBKey(tb, isLink, mname); dbkey != "" {
			return dbkey
		}
	}
	return ""
}

// 将 exp 中的标记替换为 args 中对应的数据库列名
// args 的类型是 string 或 *S_Member
// exp 中，表示列名用：$[参数索引]；表示传入值用：?[参数索引]
// 如，假设 table 对应的对象定义为：
//    type Object struct {
//	     Value string `db:"value"`
//    }
//    tbObject := NewTable("table_object", new(Object))
// 则：
//    调用：explainExp(tbObject, "#[1] join #[1]", tbObject)							// 返回：`table_object` join `table_object`
//	  调用：explainExp(tbObject, "$[1] like '%?[2]%'", "Value", "xxxx")					// 返回：`value` like '%xxxx%'
//    调用：explainExp(tbObject, "$[1] like '%?[2]%'", tbObject.M("Value"), "xxxx")		// 返回：`table_object`.`value` like '%xxxx%'
//    调用：explainExp(nil, "$[1] like '%?[2]%'", "value", "xxxx")						// 返回：`value` like '%xxxx%'
func explainExp(table *S_Table, exp string, args ...any) (newExp string, inValues []any, err error) {
	inValues = make([]any, 0)
	getArg := func(index int) interface{} {
		if index < 0 {
			err = fmt.Errorf(`argument in expression %q is not allow to less than 1.`, exp)
			return nil
		}
		if index >= len(args) {
			err = fmt.Errorf(`arguments are not enougth for expression %q`, exp)
			return nil
		}
		arg := args[index]
		if fsky.IsNil(arg) {
			err = fmt.Errorf("argument %d is not allow to be nil for expression %q", index+1, exp)
			return nil
		}
		return arg
	}

	// invalue
	repInValue := func(index int, arg any) string {
		targ := reflect.TypeOf(arg)
		if targ.Kind() == reflect.Array || targ.Kind() == reflect.Slice {
			varg := reflect.ValueOf(arg)
			es := []string{}
			for i := 0; i < varg.Len(); i++ {
				e := varg.Index(i)
				if !e.IsValid() {
					continue
				}
				if e.Type().Kind() == reflect.String {
					es = append(es, fmt.Sprintf("'%s'", e.Interface()))
				} else {
					es = append(es, fmt.Sprintf("%v", e.Interface()))
				}
			}
			return "(" + strings.Join(es, ",") + ")"
		}
		inValues = append(inValues, arg)
		return "?"
	}

	// member
	repMember := func(index int, arg any) string {
		switch arg.(type) {
		case *S_Member:
			return arg.(*S_Member).quoteWithTable()
		case string:
			mname := arg.(string)
			if table == nil {
				return mname
			}
			dbkey := getTableMemberDBKey(table, table.IsLink(), mname)
			if dbkey == "" {
				err = fmt.Errorf("table %s has no member named %q", table, mname)
				return mname
			}
			return dbkey
		}
		if table == nil {
			err = fmt.Errorf("argument %d must be a table member.", index)
		} else {
			err = fmt.Errorf("argument %d must be a table(%s) member or member name.", index, table)
		}
		return ""
	}

	// table
	repTable := func(index int, arg any) string {
		switch arg.(type) {
		case *S_Table:
			table := arg.(*S_Table)
			inValues = append(inValues, table.linkInValues...)
			return table.quote()
		}
		err = fmt.Errorf("argument %d must be a %v", index, tableType)
		return ""
	}

	replacers := map[string]func(int, any) string{
		"?": repInValue,
		"$": repMember,
		"#": repTable,
	}

	newExp = _regexp.ReplaceAllStringFunc(exp,
		func(e string) string {
			if err != nil {
				return e
			}
			index, _ := strconv.Atoi(e[2 : len(e)-1])
			arg := getArg(index - 1)
			if err != nil {
				return e
			}
			return replacers[e[0:1]](index, arg)
		})
	return
}

// -------------------------------------------------------------------
// SQL builder base
// -------------------------------------------------------------------
type s_SQL struct {
	err      error
	inValues []interface{} // 传入参数
	sqlText  string
}

func (this *s_SQL) errorf(msg string, args ...interface{}) {
	this.err = fmt.Errorf(msg, args...)
}

func (this *s_SQL) notOK() bool {
	return this.err != nil
}

func (this *s_SQL) addInValues(values ...interface{}) {
	if this.inValues == nil {
		this.inValues = make([]interface{}, 0)
	}
	this.inValues = append(this.inValues, values...)
}

func (this *s_SQL) createSQLInfo() *S_SQLInfo {
	sqlInfo := newSQLInfo(this.err, this.sqlText)
	if this.err == nil && this.inValues != nil {
		sqlInfo.InValues = this.inValues
	}
	return sqlInfo
}

// 将 exp 中的标记替换为 args 中对应的数据库列名
// args 的类型是 string 或 *S_Member
// exp 中，表示列名用：$[参数索引]；表示传入值用：?[参数索引]
// 如，假设 table 对应的对象定义为：
//    type Object struct {
//	     Value string `db:"value"`
//    }
//    tbObject := NewTable("table_object", new(Object))
// 则：
//    调用：explainExp(tbObject, "#[1] join #[1]", tbObject)							// 返回：`table_object` join `table_object`
//	  调用：explainExp(tbObject, "$[1] like '%?[2]%'", "Value", "xxxx")					// 返回：`value` like '%xxxx%'
//    调用：explainExp(tbObject, "$[1] like '%?[2]%'", tbObject.M("Value"), "xxxx")		// 返回：`table_object`.`value` like '%xxxx%'
//    调用：explainExp(nil, "$[1] like '%?[2]%'", "value", "xxxx")						// 返回：`value` like '%xxxx%'
func (this *s_SQL) explainExp(table *S_Table, exp string, args ...interface{}) string {
	exp, inValues, err := explainExp(table, exp, args...)
	if err != nil {
		this.err = err
	} else {
		this.addInValues(inValues...)
	}
	return exp
}
