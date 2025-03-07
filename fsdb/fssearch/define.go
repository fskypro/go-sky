/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: definations
@author: fanky
@version: 1.0
@date: 2024-12-18
**/

package fssearch

// -------------------------------------------------------------------
// 条件匹配器类型
// -------------------------------------------------------------------
type T_MatchType string

const (
	MT_Contain    T_MatchType = "contain"     // 包含
	MT_NoContain              = "not_contain" // 不包含
	MT_Match                  = "match"       // 匹配
	MT_NotMatch               = "not_match"   // 不匹配
	MT_Equal                  = "equal"       // 相等
	MT_NoEqual                = "not_equal"   // 不相等
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
		MT_NoContain:  true,
		MT_Match:      true,
		MT_NotMatch:   true,
		MT_Equal:      true,
		MT_NoEqual:    true,
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
