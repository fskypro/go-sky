/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 搜索参数
@author: fanky
@version: 1.0
@date: 2024-12-15
**/

package pgsearchv2

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"fsky.pro/fspgsql/fssql"
	"fsky.pro/fssearch"

	"fsky.pro/fsky"
	"fsky.pro/fssearch/dbsearchv2"
)

type S_CndInfo struct {
	*dbsearchv2.S_CndInfo
	table *fssql.S_Table
}

// -------------------------------------------------------------------
// SearchArg
// -------------------------------------------------------------------
// 通用搜索参数
type S_SearchArg struct {
	*fssearch.S_SearchArg
	cndExp  string // 解释生成的条件表达式
	orderBy string // 排序字段
	parsed  bool   // 是否已经解释过（调用过 ParseCnd 函数）
	inputs  []any  // 传入条件参数
}

func NewSearchArg(base *fssearch.S_SearchArg) *S_SearchArg {
	return &S_SearchArg{
		S_SearchArg: base,
		inputs:      make([]any, 0),
	}
}

// 调用对象自定义条件表达式
func (this *S_SearchArg) callCndParser(dbObj any, pgCnd *S_PGCndInfo) (string, []any, error) {
	method, err := dbObjCache.getCndParser(dbObj, pgCnd)
	if err != nil {
		return "", nil, err
	}
	if !method.IsValid() {
		return "", nil, nil
	}

	// 调用方法
	rets := method.Call([]reflect.Value{
		reflect.ValueOf(orderSign),
		reflect.ValueOf(pgCnd),
	})

	// 提取返回值
	exp := rets[0].Interface().(string)
	values := rets[1].Interface().([]any)
	err, _ = rets[2].Interface().(error)

	return exp, values, err
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 解释搜索条件
// 可能返回错误：
// -> fssearch.ParseCndError        ：解释搜索参数错误
// -> fssearch.InvalidCndFieldError ：无效的搜索条件字段
// -> fssearch.UnsupportMatcherError：搜索字段不支持指定的匹配器
// -> fssearch.CndValueTypeError    ：搜索条件值的类型不正确
func (this *S_SearchArg) ParseCnds(table *fssql.S_Table, dbObj any) error {
	this.parsed = true
	if len(this.Cnds) == 0 {
		return nil
	}
	cnds, err := dbsearchv2.NewCnds(this.Cnds)
	if err != nil {
		return fssearch.ParseCndError{Cnds: this.Cnds}
	}

	// 生成条件表达式
	var sb strings.Builder
	err = cnds.Tran(func(e dbsearchv2.E_CndElem, cnd *dbsearchv2.S_Cnd) error {
		switch e {
		case dbsearchv2.ECndLQuote:
			sb.WriteByte('(')
		case dbsearchv2.ECndRQuote:
			sb.WriteByte(')')
		case dbsearchv2.ECndAnd:
			sb.WriteString(" AND ")
		case dbsearchv2.ECndOr:
			sb.WriteString(" OR ")
		case dbsearchv2.ECndObj:
			if cnd.CndType == dbsearchv2.CndEmpty {
				sb.WriteString("TRUE")
				return nil
			}
			if cnd.CndType == dbsearchv2.CndObj {
				pgCnd := NewPGCndInfo(table, &cnd.S_CndInfo)
				// 优先使用对象自定义的搜索参数解释器
				exp, values, err := this.callCndParser(dbObj, pgCnd)
				if err != nil {
					return errors.Join(errors.New("custom explain search condition fail"), err)
				}
				if exp == "" {
					var value any
					exp, value, err = pgCnd.GetExp(dbObj)
					if err != nil {
						return errors.Join(errors.New("explain search condition fail"), err)
					}
					values = append(values, value)
				}
				this.inputs = append(this.inputs, values...)
				sb.WriteString(exp)
				return nil
			}
			str := strings.ToLower(cnd.String)
			if str == "true" || str == "false" {
				sb.WriteString(strings.ToUpper(cnd.String))
				return nil
			}
			return fssearch.ParseCndError{Cnds: this.Cnds}
		}
		return nil
	})
	if err != nil {
		return err
	}
	this.cndExp = sb.String()
	if this.cndExp == "" {
		this.cndExp = "TRUE"
	}

	if this.OrderBy != "" {
		mb := dbObjCache.getFieldMember(table, dbObj, this.OrderBy)
		if mb == nil {
			return fssearch.InvalidOrderByError{FieldName: this.OrderBy}
		}
		this.orderBy = mb.AutoQuote()
	}

	return nil
}

func (this *S_SearchArg) PageIndex() int {
	if this.PageSize < 1 {
		return 0
	}
	return (this.Page - 1) * this.PageSize
}

func (this *S_SearchArg) order() string {
	return fsky.IfElse(this.Desc > 0, "DESC", "ASC")
}

// ---------------------------------------------------------
// startOrder 表示 SQL 语句的参数起始序号
// 返回值：
// -> 第一个返回值表示 where 条件语句
// -> 第二个返回值为输入条件参数列表
func (this *S_SearchArg) Where(sqlInfo *fssql.S_SQL) string {
	if !this.parsed {
		return "search conditions have not been parsed"
	}
	if this.cndExp == "" {
		return "TRUE"
	}
	startOrder := sqlInfo.NextInputOrder()
	fields := strings.Split(this.cndExp, orderSign)
	wsql := fields[0]
	for i, field := range fields[1:] {
		wsql += fmt.Sprintf("$%d", i+startOrder) + field
	}
	sqlInfo.AddInputs(this.inputs)
	return wsql
}

func (this *S_SearchArg) WhereAndTail(sqlInfo *fssql.S_SQL) string {
	if !this.parsed {
		return "search conditions have not been parsed"
	}
	startOrder := sqlInfo.NextInputOrder()
	fields := strings.Split(this.cndExp, orderSign)
	wsql := fields[0]
	for i, field := range fields[1:] {
		wsql += fmt.Sprintf("$%d", i+startOrder) + field
	}
	if wsql == "" {
		wsql = "TRUE"
	}

	sqlInfo.AddInputs(this.inputs)

	if this.PageSize < 1 {
		if this.orderBy == "" {
			return wsql
		}
		return fmt.Sprintf(`%s ORDER BY %s %s`, wsql, this.orderBy, this.order())
	}
	if this.orderBy == "" {
		return fmt.Sprintf("%s LIMIT %d OFFSET %d", wsql, this.PageSize, this.PageIndex())
	}
	return fmt.Sprintf(`%s ORDER BY %s %s LIMIT %d OFFSET %d`, wsql, this.orderBy, this.order(), this.PageSize, this.PageIndex())
}
