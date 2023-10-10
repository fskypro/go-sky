/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 搜索条件
@author: fanky
@version: 1.0
@date: 2023-09-29
**/

package objsearch

import (
	"errors"
	"fmt"
	"strings"

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
	Key      string
	Match    string
	Value    string
	comparer func(any, string, string) (bool, error)
	conf i_Config
}

func newCnd(conf i_Config) *s_Cnd {
	return &s_Cnd{
		conf: conf,
	}
}

func (this *s_Cnd) parse(cnd map[string]any) (i_Cnd, error) {
	if v, ok := cnd["col"]; !ok {
		return nil, errors.New("condition key is not indicated")
	} else if fstype.IsType[string](v) {
		this.Key = v.(string)
	} else {
		return nil, fmt.Errorf("%q in search condition must be a string", "col")
	}

	if v, ok := cnd["match"]; !ok {
		return nil, errors.New("condition match is not indicated")
	} else if fstype.IsType[string](v) {
		this.Match = v.(string)
	} else {
		return nil, fmt.Errorf("%q in search condition must be a string", "match")
	}

	if v, ok := cnd["value"]; !ok {
		return nil, errors.New("condition value is not indicated")
	} else if fstype.IsType[string](v) {
		this.Value = v.(string)
	} else {
		return nil, fmt.Errorf("%q in search condition must be a string", "value")
	}

	comparer := this.conf.GetMatchHandler(this.Match)
	this.comparer = cmpHandlers.handlers[comparer]
	if this.comparer == nil {
		return nil, fmt.Errorf("compare handler(match=%q, comparer=%q) is not exists", this.Match, comparer)
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

func (this *s_AndCnds) parse(conf i_Config, anyCnds []any) (i_Cnd, error) {
	this.cnds = make([]i_Cnd, 0)
	for _, anyCnd := range anyCnds {
		cnd, err := parseCnd(conf, anyCnd)
		if err != nil { return nil, err }
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

func (this *s_OrCnds) parse(conf i_Config, anyCnds []any) (i_Cnd, error) {
	this.cnds = make([]i_Cnd, 0)
	for _, anyCnd := range anyCnds {
		cnd, err := parseCnd(conf, anyCnd)
		if err != nil { return nil, err }
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
func parseCnd(conf i_Config, anyCnd any) (i_Cnd, error) {
	switch anyCnd.(type) {
	case map[string]any:
		// 单个配置条件
		return newCnd(conf).parse(anyCnd.(map[string]any))
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
			return new(s_OrCnds).parse(conf, anyCnds)
		}
		return new(s_AndCnds).parse(conf, anyCnds)
	}
	return nil, fmt.Errorf(`invalid search condition("%#v")`, anyCnd)
}

