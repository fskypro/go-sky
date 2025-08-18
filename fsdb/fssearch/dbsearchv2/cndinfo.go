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

	"fsky.pro/fssearch"
)

// 匹配类型
type T_MatchType = fssearch.T_MatchType

type S_CndInfo struct {
	MatchType T_MatchType `json:"match"` // 条件匹配器（不同数据库有各自不同的匹配器，需要数据库引擎自己定义）
	FieldName string      `json:"field"` // 数据库字段对外映射的名称
	Value     any         `json:"value"` // 条件判断值
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
	matchType, ok := match.(string)
	if !ok {
		return fmt.Errorf("value of sql condition key %q must be a string", "match")
	}
	this.MatchType = T_MatchType(matchType)

	// Field
	field, ok := mcnd["field"]
	if !ok {
		return fmt.Errorf("sql condition must contains key %q", "field")
	}
	this.FieldName, ok = field.(string)
	if !ok {
		return fmt.Errorf("value of sql condition key %q or %q must be a string", "field", "col")
	}

	// Value
	value, ok := mcnd["value"]
	if !ok {
		return fmt.Errorf("sql condition must contains %q key", "value")
	}
	this.Value = value
	return nil
}
