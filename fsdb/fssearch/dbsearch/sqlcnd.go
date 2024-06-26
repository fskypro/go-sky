/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: json sql condition definations
@author: fanky
@version: 1.0
@date: 2023-06-06
**/

package dbsearch

import (
	"errors"
	"fmt"
	"strings"

	"fsky.pro/fsky"
)

// -------------------------------------------------------------------
// 表达式类型
// -------------------------------------------------------------------
type E_CndElem int

const (
	ECndLQuote E_CndElem = iota // 左括号
	ECndRQuote                  // 右括号
	ECndAnd                     // AND 连接
	ECndOr                      // OR 连接
	ECndObj                     // 条件对象
)

var cndElems = map[E_CndElem]string{
	ECndLQuote: "(",
	ECndRQuote: ")",
	ECndAnd:    "and",
	ECndOr:     "or",
	ECndObj:    "cnd",
}

// -------------------------------------------------------------------
// 条件类型
// -------------------------------------------------------------------
type E_CndType string

const (
	CndObj   E_CndType = "obj"   // 配置类型条件
	CndStr   E_CndType = "str"   // 字符串类型条件
	CndEmpty E_CndType = "empty" // 空条件（允许条件为：[]）
)

// -------------------------------------------------------------------
// SQL condition interface
// -------------------------------------------------------------------
type i_Cnd interface {
	tran(func(E_CndElem, *S_Cnd) error) error
}

func _parseAnyCnd(anyCnd any) (i_Cnd, error) {
	switch anyCnd.(type) {
	case map[string]any:
		// 单个配置条件
		cnd := newCnd(CndObj)
		if err := cnd.parse(anyCnd.(map[string]any)); err != nil {
			return nil, fmt.Errorf("invalid sql condition, %v", err)
		}
		return cnd, nil
	case string:
		cnd := newCnd(CndStr)
		cnd.String = anyCnd.(string)
		return cnd, nil
	case []any:
		// 复合条件
		cndList := newCndList()
		err := cndList.parse(anyCnd.([]any))
		if err != nil {
			return nil, fmt.Errorf("invalid sql condition, %v", err)
		}
		return cndList, nil
	}
	return nil, fmt.Errorf(`invalid sql condition("%#v")`, anyCnd)
}

// -------------------------------------------------------------------
// SQL condition
// -------------------------------------------------------------------
type S_Cnd struct {
	i_Cnd
	Match string // 条件匹配方式（如：匹配、包含、等于、小于、大于等）
	Key   string // 参与条件判断的字段映射名称
	Value string // 条件判断值（如：like '%包含值%'）

	Type   E_CndType `json:"-"` // 条件类型
	String string    `json:"-"` // 字符串形式的条件(是否正确，由用户在 tran 函数中自行判断)
}

func newCnd(ctype E_CndType) *S_Cnd {
	return &S_Cnd{Type: ctype}
}

func (this *S_Cnd) parse(mcnd map[string]any) error {
	if mcnd == nil {
		return errors.New("sql condition is not allow to be a nil value")
	}
	match, ok := mcnd["match"]
	if !ok {
		return fmt.Errorf("sql condition must contains %q key", "match")
	}
	this.Match, ok = match.(string)
	if !ok {
		return fmt.Errorf("value of sql condition key %q must be a string", "match")
	}

	key, ok := mcnd["key"]
	if !ok {
		key, ok = mcnd["col"]
		if !ok {
			return fmt.Errorf("sql condition must contains key %q or %q", "key", "col")
		}
	}

	this.Key, ok = key.(string)
	if !ok {
		return fmt.Errorf("value of sql condition key %q or %q must be a string", "key", "col")
	}

	value, ok := mcnd["value"]
	if !ok {
		return fmt.Errorf("sql condition must contains %q key", "value")
	}
	this.Value = fmt.Sprintf("%v", value)
	return nil
}

func (this *S_Cnd) tran(f func(E_CndElem, *S_Cnd) error) error {
	return f(ECndObj, this)
}

// -------------------------------------------------------------------
// AND/OR Conditions
// -------------------------------------------------------------------
type s_CndList struct {
	i_Cnd
	and  bool
	cnds []i_Cnd
}

func newCndList() *s_CndList {
	return &s_CndList{
		and:  true,
		cnds: make([]i_Cnd, 0),
	}
}

func (this *s_CndList) parse(lcnds []any) error {
	if lcnds == nil {
		return errors.New("sql condition is not allow to be a nil value")
	}
	if len(lcnds) == 0 {
		this.cnds = append(this.cnds, newCnd(CndEmpty))
		return nil
	}
	first := lcnds[0]
	join, ok := first.(string)
	if ok {
		join := strings.ToLower(join)
		if join == "or" {
			this.and = false
		} else if join == "and" {
			this.and = true
		} else {
			return fmt.Errorf("invalid sql condition join mark %q", first)
		}
		lcnds = lcnds[1:]
	}
	for _, anyCnd := range lcnds {
		cnd, err := _parseAnyCnd(anyCnd)
		if err != nil { return err } else {
			this.cnds = append(this.cnds, cnd)
		}
	}
	return nil
}

func (this *s_CndList) tran(f func(E_CndElem, *S_Cnd) error) error {
	if len(this.cnds) == 0 { return nil }
	f(ECndLQuote, nil)
	for idx, cnd := range this.cnds {
		if idx > 0 {
			f(fsky.IfElse(this.and, ECndAnd, ECndOr), nil)
		}
		if err := cnd.tran(f); err != nil {
			return err
		}
	}
	return f(ECndRQuote, nil)
}

// -------------------------------------------------------------------
// SQL coditions
// 假设 A、B、C、D 都是 S_Cnd 对象
// 以下写法都是正确的：
//   1、只有 A 条件时传入        ：A
//   2、(A || B) && (C || D) 传入：[["OR", A, B], ["OR", C, D]]
//   3、A && ((B || C) && D) 传入：["AND", [A, ["OR", B || C], D]]
// 即：
//   如果只有一个条件，可以直接写单个条件，即：{"match": <匹配方式>, "key": <数据库字段映射名>, "value": <条件判断字符串>}
//   如果是多个条件或连接：["OR", 条件1， 条件2, ...]
//   如果是多个条件与连接：["AND", 条件1， 条件2, ...]，或：["AND", 条件1， 条件2, ...]
// -------------------------------------------------------------------
type S_Cnds struct {
	cnd i_Cnd
}

func NewCnds(anyCnd any) (cnds *S_Cnds, err error) {
	cnd, e := _parseAnyCnd(anyCnd)
	if e != nil {
		err = fmt.Errorf("parse sql condition fail, %v", e)
	} else {
		cnds = new(S_Cnds)
		cnds.cnd = cnd
	}
	return
}

func (this *S_Cnds) Tran(f func(E_CndElem, *S_Cnd) error) error {
	if this.cnd == nil {
		return nil
	}
	return this.cnd.tran(f)
}
