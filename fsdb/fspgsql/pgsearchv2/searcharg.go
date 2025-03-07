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
	"strings"

	"fsky.pro/fspgsql/fssql"

	"fsky.pro/fsky"
	"fsky.pro/fssearch/dbsearchv2"
)

// 根据查询条件对象，返回对应的匹配表达式
// 如，可能返回："name" like "%tom"
type F_GetMatchExp func(*dbsearchv2.S_CndInfo) (S_BaseCnd, error)

// -------------------------------------------------------------------
// SearchArg
// -------------------------------------------------------------------
// 通用搜索参数
type S_SearchArg struct {
	Cnds     []any  `json:"cnds"`     // 用户自定义的原始搜索条件
	Page     int    `json:"page"`     // 要查询的页码（第一页为 1）
	PageSize int    `json:"pageSize"` // 每页最大数量
	OrderBy  string `json:"orderBy"`  // 排序关键字
	Desc     int8   `json:"desc"`     // 是否倒序排列

	cndExp string
	parsed bool
	Inputs []any `json:"-"` // 传入条件参数
}

func NewSearchArg() *S_SearchArg {
	return &S_SearchArg{
		Cnds:     make([]any, 0),
		Page:     1,
		PageSize: 10,
		Inputs:   make([]any, 0),
	}
}

// 解释搜索条件
func (this *S_SearchArg) ParseCnd(fun F_GetMatchExp) error {
	this.parsed = true
	if this.Cnds == nil {
		return nil
	}
	cnds, err := dbsearchv2.NewCnds(this.Cnds)
	if err != nil {
		return fmt.Errorf("parse search condition(%#v) fail, %v", this.Cnds, err)
	}
	if cnds == nil {
		return errors.New("no conditions in search arg")
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
				baseCnd, err := fun(&cnd.S_CndInfo)
				if err != nil {
					return fmt.Errorf("format condition string fail, %v", err)
				}
				exp, value, err := getExp(&cnd.S_CndInfo, baseCnd)
				if exp == "" {
					return fmt.Errorf("explain search condition fail, %v", err)
				}
				this.Inputs = append(this.Inputs, value)
				sb.WriteString(exp)
				return nil
			}
			str := strings.ToLower(cnd.String)
			if str == "true" || str == "false" {
				sb.WriteString(strings.ToUpper(cnd.String))
				return nil
			}
			return fmt.Errorf("string %q is not a sql condition expression", cnd.String)
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
	return nil
}

func (this *S_SearchArg) PageIndex() int {
	if this.PageSize < 1 {
		return 0
	}
	return (this.Page - 1) * this.PageSize
}

func (this *S_SearchArg) Order() string {
	return fsky.IfElse(this.Desc > 0, "DESC", "ASC")
}

// ---------------------------------------------------------
// startOrder 表示参数起始序号
func (this *S_SearchArg) Where(sqlInfo *fssql.S_SQL) string {
	if !this.parsed {
		return "<unparsed search argument>"
	}
	startOrder := sqlInfo.NextInputOrder()
	segs := strings.Split(this.cndExp, orderSign)
	exp := segs[0]
	for i, seg := range segs[1:] {
		exp += fmt.Sprintf("$%d", i+startOrder) + seg
	}
	sqlInfo.AddInputs(this.Inputs)
	return exp
}

func (this *S_SearchArg) WhereAndTail(sqlInfo *fssql.S_SQL) string {
	if !this.parsed {
		return "<unparsed search argument>"
	}
	startOrder := sqlInfo.NextInputOrder()
	segs := strings.Split(this.cndExp, orderSign)
	exp := segs[0]
	for i, seg := range segs[1:] {
		exp += fmt.Sprintf("$%d", i+startOrder) + seg
	}

	sqlInfo.AddInputs(this.Inputs)
	if this.PageSize < 1 {
		if this.OrderBy == "" {
			return exp
		}
		return fmt.Sprintf(`%s ORDER BY "%s" %s`, exp, this.OrderBy, this.Order())
	}
	if this.OrderBy == "" {
		return fmt.Sprintf("%s LIMIT %d OFFSET %d", exp, this.PageSize, this.PageIndex())
	}
	return fmt.Sprintf(`%s ORDER BY "%s" %s LIMIT %d OFFSET %d`, exp, this.OrderBy, this.Order(), this.PageSize, this.PageIndex())
}
