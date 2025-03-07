/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: search definations
@author: fanky
@version: 1.0
@date: 2024-12-15
**/

package dbsearchv2

import (
	"errors"
	"fmt"

	"fsky.pro/fstype"
)

type S_CndInfo struct {
	MatchName string `json:"match"` // 条件匹配器名称（如：匹配、包含、等于、小于、大于等）
	FieldName string `json:"field"` // 数据库字段对外映射的名称
	FieldType string `json:"ftype"` // 字段值类型
	Value     any    `json:"value"` // 条件判断值
}

// 把 map 形式的条件转换为对象
func (this *S_CndInfo) parse(mcnd map[string]any) error {
	if mcnd == nil {
		return errors.New("sql condition is not allow to be a nil value")
	}

	// Match
	match, ok := mcnd["match"]
	if !ok {
		return fmt.Errorf("sql condition must contains %q key", "match")
	}
	this.MatchName, ok = match.(string)
	if !ok {
		return fmt.Errorf("value of sql condition key %q must be a string", "match")
	}

	// Field
	field, ok := mcnd["field"]
	if !ok {
		return fmt.Errorf("sql condition must contains key %q", "field")
	}
	this.FieldName, ok = field.(string)
	if !ok {
		return fmt.Errorf("value of sql condition key %q or %q must be a string", "field", "col")
	}

	// Quote
	if _, ok := mcnd["ftype"]; !ok {
		return fmt.Errorf("sql condition must contains key %q", "ftype")
	}
	ftype, err := fstype.Convert[string](mcnd["ftype"])
	if err != nil {
		return fmt.Errorf("value of sql condition key %q must be an int value", "quote")
	}
	this.FieldType = ftype

	// Value
	value, ok := mcnd["value"]
	if !ok {
		return fmt.Errorf("sql condition must contains %q key", "value")
	}
	this.Value = value
	return nil
}
