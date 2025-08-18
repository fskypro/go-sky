/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: definations
@author: fanky
@version: 1.0
@date: 2024-12-18
**/

package fssearch

import (
	"reflect"
	"slices"
	"time"
)

// -------------------------------------------------------------------
// 条件匹配器类型
// -------------------------------------------------------------------
type T_MatchType string

const (
	MT_Contain    T_MatchType = "contain"     // 包含
	MT_NotContain             = "not_contain" // 不包含
	MT_Match                  = "match"       // 匹配
	MT_NotMatch               = "not_match"   // 不匹配
	MT_Equal                  = "equal"       // 相等
	MT_NotEqual               = "not_equal"   // 不相等
	MT_Less                   = "less"        // 小于
	MT_LessEqual              = "less_equal"  // 小于等于
	MT_Large                  = "large"       // 大于
	MT_LargeEqual             = "large_equal" // 大于等于
	MT_ReMatch                = "re_match"    // 正则匹配

	MT_InArray          = "in_array"           // 在指定数组中
	MT_ArrLenEqual      = "arrlen_equal"       // 数组长度等于
	MT_ArrLenLess       = "arrlen_less"        // 数组长度小于
	MT_ArrLenLessEqual  = "arrlen_less_equal"  // 数组长度小于等于
	MT_ArrLenLarge      = "arrlen_large"       // 数组长度大于
	MT_ArrLenLargeEqual = "arrlen_large_equal" // 数组长度大于等于
)

func (self T_MatchType) Valid() bool {
	return map[T_MatchType]bool{
		MT_Contain:    true,
		MT_NotContain: true,
		MT_Match:      true,
		MT_NotMatch:   true,
		MT_Equal:      true,
		MT_NotEqual:   true,
		MT_Less:       true,
		MT_LessEqual:  true,
		MT_Large:      true,
		MT_LargeEqual: true,
		MT_ReMatch:    true,

		MT_InArray:          true,
		MT_ArrLenEqual:      true,
		MT_ArrLenLess:       true,
		MT_ArrLenLessEqual:  true,
		MT_ArrLenLarge:      true,
		MT_ArrLenLargeEqual: true,
	}[self]
}

// -------------------------------------------------------------------
// 判断指定匹配器是否适合应用与GO值类型
// -------------------------------------------------------------------
var timeKind = reflect.TypeOf(time.Time{}).Kind()
var _numberMatchs = []T_MatchType{MT_Equal, MT_NotEqual, MT_Less, MT_LessEqual, MT_Large, MT_LargeEqual, MT_LargeEqual}
var _stringMatchs = []T_MatchType{MT_Contain, MT_NotContain, MT_Match, MT_NotMatch, MT_Equal, MT_NotEqual, MT_ReMatch}
var _arrayMatchs = []T_MatchType{MT_InArray, MT_ArrLenEqual, MT_ArrLenLess, MT_ArrLenLessEqual, MT_ArrLenLarge, MT_ArrLenLargeEqual}

func IsTypeMatch(t reflect.Type, match T_MatchType) bool {
	switch t.Kind() {
	case reflect.String:
		return slices.Contains(_stringMatchs, match)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return slices.Contains(_numberMatchs, match)
	case reflect.Array:
		return slices.Contains(_arrayMatchs, match)
	case timeKind:
		return slices.Contains(_numberMatchs, match)
	}
	return false
}

// -------------------------------------------------------------------
// 条件搜索参数
// -------------------------------------------------------------------
type S_SearchArg struct {
	Cnds     []any  `json:"cnds"`     // 用户自定义的原始搜索条件
	Page     int    `json:"page"`     // 要查询的页码（第一页为 1）
	PageSize int    `json:"pageSize"` // 每页最大数量
	OrderBy  string `json:"orderBy"`  // 排序字段
	Desc     int8   `json:"desc"`     // 是否倒序排列
}

func NewSearchArg() *S_SearchArg {
	return &S_SearchArg{
		Cnds:     make([]any, 0),
		Page:     1,
		PageSize: 50,
	}
}
