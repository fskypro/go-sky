/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 搜索条件
@author: fanky
@version: 1.0
@date: 2025-01-08
**/

package objsearchv2

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"fsky.pro/fssearch"

	"fsky.pro/fstype"
)

// -----------------------------------------------------------------------------
// cnd
// -----------------------------------------------------------------------------
type i_Cnd interface {
	compare(obj any) (bool, error)
}

// -------------------------------------------------------------------
// Single Condition
// -------------------------------------------------------------------
type s_Cnd struct {
	i_Cnd
	Key       string      // 要搜索的字段，在表达式中可以用以下名称表示：field/key/col
	MatchType T_MatchType // 比较运算符
	Value     *s_CmpValue // 比较值
	comparer  func(any, string, *s_CmpValue) (bool, error)
}

func newCnd() *s_Cnd {
	return &s_Cnd{}
}

// 可能返回：
// -> fssearch.ParseCndError
// -> fssearch.InvalidMatchError
func (this *s_Cnd) parse(cnd map[string]any) (i_Cnd, error) {
	// 解释条件中的 key/col
	v, ok := cnd["key"]
	if !ok {
		v, ok = cnd["col"]
	}
	if !ok {
		v, ok = cnd["field"]
	}

	if !ok {
		err := fmt.Errorf("search condition object must contain a field named %q", "field")
		return nil, errors.Join(fssearch.ParseCndError{Cnds: cnd}, err)
	} else if !fstype.IsType[string](v) {
		err := fmt.Errorf("search condition value of %q must be a string, but not %v", "field", reflect.TypeOf(v))
		return nil, errors.Join(fssearch.ParseCndError{Cnds: cnd}, err)
	} else {
		this.Key = v.(string)
	}
	if this.Key == "" {
		err := fmt.Errorf("search condition must be a unempty string")
		return nil, errors.Join(fssearch.ParseCndError{Cnds: cnd}, err)
	}

	// 解释条件中的 match
	if v, ok := cnd["match"]; !ok {
		err := fmt.Errorf("search condition object must contain a field named %q", "match")
		return nil, errors.Join(fssearch.ParseCndError{Cnds: cnd}, err)
	} else if fstype.IsType[string](v) {
		this.MatchType = T_MatchType(v.(string))
	} else {
		err := fmt.Errorf("search condition value of %q must be a string, but not %v", "match", reflect.TypeOf(v))
		return nil, errors.Join(fssearch.ParseCndError{Cnds: cnd}, err)
	}

	// 解释条件中的 value
	if v, ok := cnd["value"]; !ok {
		err := fmt.Errorf("search condition object must contain a field named %q", "value")
		return nil, errors.Join(fssearch.ParseCndError{Cnds: cnd}, err)
	} else {
		this.Value = newCmpValue(v)
	}

	// 匹配器检查
	if !this.MatchType.Valid() {
		return nil, fssearch.InvalidMatchError{FieldName: this.Key, Match: this.MatchType}
	}
	this.comparer = cmpHandlers.handlers[this.MatchType]
	if this.comparer == nil {
		return nil, fssearch.InvalidMatchError{FieldName: this.Key, Match: this.MatchType}
	}
	return this, nil
}

func (this *s_Cnd) compare(obj any) (bool, error) {
	return this.comparer(obj, this.Key, this.Value)
}

// -------------------------------------------------------------------
// And Conditions
// -------------------------------------------------------------------
type s_AndCnds struct {
	i_Cnd
	cnds []i_Cnd
}

func (this *s_AndCnds) parse(anyCnds []any) (i_Cnd, error) {
	this.cnds = make([]i_Cnd, 0)
	for _, anyCnd := range anyCnds {
		cnd, err := parseCnd(anyCnd)
		if err != nil {
			return nil, err
		}
		this.cnds = append(this.cnds, cnd)
	}
	return this, nil
}

func (this *s_AndCnds) compare(obj any) (bool, error) {
	for _, cnd := range this.cnds {
		ok, err := cnd.compare(obj)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// -------------------------------------------------------------------
// Or Conditions
// -------------------------------------------------------------------
type s_OrCnds struct {
	i_Cnd
	cnds []i_Cnd
}

func (this *s_OrCnds) parse(anyCnds []any) (i_Cnd, error) {
	this.cnds = make([]i_Cnd, 0)
	for _, anyCnd := range anyCnds {
		cnd, err := parseCnd(anyCnd)
		if err != nil {
			return nil, err
		}
		this.cnds = append(this.cnds, cnd)
	}
	return this, nil
}

func (this *s_OrCnds) compare(obj any) (bool, error) {
	for _, cnd := range this.cnds {
		ok, err := cnd.compare(obj)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// -------------------------------------------------------------------
// parser
// -------------------------------------------------------------------
func parseCnd(anyCnd any) (i_Cnd, error) {
	switch anyCnd.(type) {
	case map[string]any:
		// 单个配置条件
		return newCnd().parse(anyCnd.(map[string]any))
	case []any:
		// 复合条件
		cndOR := false
		anyCnds := []any{}
		for _, anyCnd := range anyCnd.([]any) {
			if fstype.IsType[string](anyCnd) {
				if strings.ToLower(anyCnd.(string)) == "or" {
					cndOR = true
				}
				continue
			}
			anyCnds = append(anyCnds, anyCnd)
		}
		if cndOR {
			return new(s_OrCnds).parse(anyCnds)
		}
		return new(s_AndCnds).parse(anyCnds)
	}
	return nil, fssearch.ParseCndError{Cnds: anyCnd}
}
