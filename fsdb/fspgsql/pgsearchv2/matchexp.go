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

	"fsky.pro/fssearch"
	"fsky.pro/fssearch/dbsearchv2"
	"fsky.pro/fstype"
)

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
// 内部搜索条件信息
// -------------------------------------------------------------------
type S_BaseCnd struct {
	MatchType fssearch.T_MatchType // 匹配器，由 dbsearchv2.S_CndInfo.MatchName 映射
	ColName   string               // 数据库字段，由 dbsearchv2.S_CndInfo.FieldName 映射
}

type t_TypeExp map[t_ValueType]string

var matchExps = map[fssearch.T_MatchType]t_TypeExp{
	fssearch.MT_Contain:    t_TypeExp{vt_All: "%s LIKE %v", vt_Array: "%s @> %v"},            // 包含
	fssearch.MT_NoContain:  t_TypeExp{vt_String: "%s NOT LIKE %v", vt_Array: "NOT %s @> %v"}, // 不包含
	fssearch.MT_Match:      t_TypeExp{vt_String: "%s=%v", vt_Array: "%s=%v"},                 // 匹配
	fssearch.MT_NotMatch:   t_TypeExp{vt_String: "%s<>%v", vt_Array: "%s<>%v"},               // 不匹配
	fssearch.MT_Equal:      t_TypeExp{vt_All: "%s=%v"},                                       // 等于
	fssearch.MT_NoEqual:    t_TypeExp{vt_All: "%s<>%v"},                                      // 不等于
	fssearch.MT_Less:       t_TypeExp{vt_All: "%s<%v"},                                       // 小于
	fssearch.MT_LessEqual:  t_TypeExp{vt_All: "%s<=%v"},                                      // 小于等于
	fssearch.MT_Large:      t_TypeExp{vt_All: "%s>%v"},                                       // 大于
	fssearch.MT_LargeEqual: t_TypeExp{vt_All: "%s>=%v"},                                      // 大于等于
	fssearch.MT_ReMatch:    t_TypeExp{vt_String: "%s~%v"},                                    // 正则匹配
}

func (self t_TypeExp) get(t t_ValueType) string {
	if self[t] == "" {
		return self[vt_All]
	}
	return self[t]
}

// 返回匹配器对应的比较表达式
// 如果不是合法的匹配器，则返回空字符串
func getExp(cndInfo *dbsearchv2.S_CndInfo, baseCnd S_BaseCnd) (string, any, error) {
	texp := matchExps[baseCnd.MatchType]
	if texp == nil {
		return "", nil, fmt.Errorf("unsupport match type %q", baseCnd.MatchType)
	}
	exp := ""
	value := cndInfo.Value
	if fstype.IsNumber(cndInfo.Value) {
		exp = texp.get(vt_Number)
	} else if fstype.IsString(cndInfo.Value) {
		exp = texp.get(vt_String)
		v, _ := fstype.Convert[string](value)
		if baseCnd.MatchType == fssearch.MT_Contain || baseCnd.MatchType == fssearch.MT_NoContain {
			value = "%" + v + "%"
		}
	} else {
		t := reflect.TypeOf(cndInfo.Value)
		if t != nil && (t.Kind() == reflect.Array || t.Kind() == reflect.Slice) {
			exp = texp.get(vt_Array)
		} else {
			exp = texp.get(vt_All)
		}
	}
	if exp == "" {
		return "", nil, fmt.Errorf("field %q unsupprt match type %q", cndInfo.FieldName, baseCnd.MatchType)
	}
	return fmt.Sprintf(exp, baseCnd.ColName, orderSign), value, nil
}
