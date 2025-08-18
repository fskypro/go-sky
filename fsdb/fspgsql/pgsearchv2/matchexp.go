/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: search definations
@author: fanky
@version: 1.0
@date: 2024-12-15
**/

package pgsearchv2

import (
	"fmt"
	"reflect"

	"fsky.pro/fspgsql/fssql"

	"fsky.pro/fssearch"
	"fsky.pro/fssearch/dbsearchv2"
	"fsky.pro/fstype"
)

// 输入值占位符
const orderSign = "${\x02}"

// -------------------------------------------------------------------
// 条件值类型
// -------------------------------------------------------------------
type t_ValueType string

const (
	vt_All                = ""
	vt_String t_ValueType = "string" // 字符串类型
	vt_Number             = "number" // 数值类型
	vt_Array              = "array"  // 数组类型
)

// -------------------------------------------------------------------
// 比较运算符与表达式映射
// -------------------------------------------------------------------
type t_TypeExp map[t_ValueType]string

var matchExps = map[fssearch.T_MatchType]t_TypeExp{
	fssearch.MT_Contain:    {vt_All: "%s LIKE %v", vt_Array: "%s @> %v"},            // 包含
	fssearch.MT_NotContain: {vt_String: "%s NOT LIKE %v", vt_Array: "NOT %s @> %v"}, // 不包含
	fssearch.MT_Match:      {vt_String: "%s=%v", vt_Array: "%s=%v"},                 // 匹配
	fssearch.MT_NotMatch:   {vt_String: "%s<>%v", vt_Array: "%s<>%v"},               // 不匹配
	fssearch.MT_Equal:      {vt_All: "%s=%v"},                                       // 等于
	fssearch.MT_NotEqual:   {vt_All: "%s<>%v"},                                      // 不等于
	fssearch.MT_Less:       {vt_All: "%s<%v"},                                       // 小于
	fssearch.MT_LessEqual:  {vt_All: "%s<=%v"},                                      // 小于等于
	fssearch.MT_Large:      {vt_All: "%s>%v"},                                       // 大于
	fssearch.MT_LargeEqual: {vt_All: "%s>=%v"},                                      // 大于等于
	fssearch.MT_ReMatch:    {vt_String: "%s~%v"},                                    // 正则匹配
}

func (self t_TypeExp) get(t t_ValueType) string {
	if self[t] == "" {
		return self[vt_All]
	}
	return self[t]
}

// -------------------------------------------------------------------
// 内部搜索条件信息
// -------------------------------------------------------------------
type S_PGCndInfo struct {
	*dbsearchv2.S_CndInfo
	Table *fssql.S_Table
}

func NewPGCndInfo(table *fssql.S_Table, baseCnd *dbsearchv2.S_CndInfo) *S_PGCndInfo {
	return &S_PGCndInfo{
		S_CndInfo: baseCnd,
		Table:     table,
	}
}

// 返回匹配器对应的比较表达式
func (this *S_PGCndInfo) GetExp(dbObj any) (string, any, error) {
	mb := dbObjCache.getFieldMember(this.Table, dbObj, this.FieldName)
	if mb == nil {
		return "", nil, fssearch.InvalidCndFieldError{FieldName: this.FieldName}
	}

	texp := matchExps[this.MatchType]
	if texp == nil {
		return "", nil, fssearch.InvalidMatchError{FieldName: this.FieldName, Match: this.MatchType}
	}

	if !fssearch.IsTypeMatch(mb.Type(), this.MatchType) {
		return "", nil, fssearch.InvalidMatchError{FieldName: this.FieldName, Match: this.MatchType}
	}

	exp := ""
	value := this.Value
	if fstype.IsNumber(this.Value) {
		exp = texp.get(vt_Number)
	} else if fstype.IsString(this.Value) {
		exp = texp.get(vt_String)
		v, _ := fstype.Convert[string](value)
		if this.MatchType == fssearch.MT_Contain || this.MatchType == fssearch.MT_NotContain {
			value = "%" + v + "%"
		}
	} else {
		t := reflect.TypeOf(this.Value)
		if t != nil && (t.Kind() == reflect.Array || t.Kind() == reflect.Slice) {
			exp = texp.get(vt_Array)
		} else {
			exp = texp.get(vt_All)
		}
	}
	if exp == "" {
		return "", nil, fssearch.CndValueTypeError{FieldName: this.FieldName, ValidType: mb.Type(), InputType: reflect.TypeOf(this.Value)}
	}
	return fmt.Sprintf(exp, mb.AutoQuote(), orderSign), value, nil
}
