/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: errors
@author: fanky
@version: 1.0
@date: 2025-07-20
**/

package fssearch

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// -------------------------------------------------------------------
// 解释搜索参数错误
// -------------------------------------------------------------------
type ParseCndError struct {
	Cnds any
}

func (self ParseCndError) Error() string {
	data, _ := json.Marshal(self.Cnds)
	return fmt.Sprintf("invalid search conditions: %s", string(data))
}

// -------------------------------------------------------------------
// 搜索条件参数在对象中不存在
// -------------------------------------------------------------------
type InvalidCndFieldError struct {
	FieldName string
}

func (self InvalidCndFieldError) Error() string {
	return fmt.Sprintf("invalid search condition field %q", self.FieldName)
}

// -------------------------------------------------------------------
// 字段不支持指定的匹配表达式
// -------------------------------------------------------------------
type InvalidMatchError struct {
	FieldName string
	Match     T_MatchType
}

func (self InvalidMatchError) Error() string {
	return fmt.Sprintf("field %q in conditions is not support match expression %q", self.FieldName, self.Match)
}

// -------------------------------------------------------------------
// 条件值类型错误
// -------------------------------------------------------------------
type CndValueTypeError struct {
	FieldName string
	ValidType reflect.Type
	InputType reflect.Type
}

func (self CndValueTypeError) Error() string {
	return fmt.Sprintf("value type of condition field %q must be a %v but not %v", self.FieldName, self.ValidType, self.InputType)
}

// -------------------------------------------------------------------
// 比较正则表达式错误
// -------------------------------------------------------------------
type RePatternError struct {
	FieldName string
	Pattern   string
}

func (self RePatternError) Error() string {
	return fmt.Sprintf("invalid regexp pattern %q of condition filed %q", self.Pattern, self.FieldName)
}

// -------------------------------------------------------------------
// OrderBy 字段不存在
// -------------------------------------------------------------------
type InvalidOrderByError struct {
	FieldName string
}

func (self InvalidOrderByError) Error() string {
	return fmt.Sprintf("invalid orderby field %q", self.FieldName)
}
