/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: errors
@author: fanky
@version: 1.0
@date: 2024-05-31
**/

package objsearch

import (
	"fmt"
	"reflect"
	"strings"
)

// -----------------------------------------------------------------------------
// 对象比较时的错误
// -----------------------------------------------------------------------------
// 对象中不存在条件表达式中指定的字段
type ErrNoCndMember struct {
	error
	Type reflect.Type
	Key  string
}

func makeErrNoCndMember(t reflect.Type, key string) ErrNoCndMember {
	return ErrNoCndMember{
		error: fmt.Errorf("no member named %q in object type %v", key, t),
		Type:  t,
		Key:   key,
	}
}

// ---------------------------------------------------------
// 不是合法的匹配器
type ErrLegalMatcher struct {
	error
	Key     string
	Matcher string
}

func makeErrLegalMatcher(key, matcher string) ErrLegalMatcher {
	return ErrLegalMatcher{
		error:   fmt.Errorf("error condition matcher of key %q", key),
		Key:     key,
		Matcher: matcher,
	}
}

// ---------------------------------------------------------
// 字段不支持匹配符号
type ErrUnsupportMatcher struct {
	error
	Key    string   // 字段名称
	Matchs []string // 匹配符号
}

func makeErrUnsupportMatcher(key string, matchs string) ErrUnsupportMatcher {
	return ErrUnsupportMatcher{
		error:  fmt.Errorf("object member %q unsupport match compare %q", key, matchs),
		Key:    key,
		Matchs: strings.Split(matchs, "|"),
	}
}

// ---------------------------------------------------------
// 传入的条件值类型不正确，无法比较
type ErrCndValue struct {
	error
	Key      string       // 字段名称
	KType    reflect.Type // 字段实际类型
	CndVType reflect.Type // 传入的条件值类型
}

func makeErrCndValue(key string, ktype, cndVType reflect.Type) ErrCndValue {
	return ErrCndValue{
		error:    fmt.Errorf("error condition value type of key %q, it must be a %v, but not a %v", key, ktype, cndVType),
		Key:      key,
		KType:    ktype,
		CndVType: cndVType,
	}
}

// ---------------------------------------------------------
// 传入的日期/时间条件值类型不正确，无法比较
type ErrCndTimeValue struct {
	error
	Key      string       // 字段名称
	KType    reflect.Type // 字段实际类型
	CndVType reflect.Type // 传入的条件值类型
	Samples  []string     // 正确的格式列表
}

func makeErrCndTimeValue(key string, ktype, cndVType reflect.Type) ErrCndTimeValue {
	samples := []string{"2006-01-02", "2006-01-02 15:04:05"}
	return ErrCndTimeValue{
		error: fmt.Errorf("condition value of key %q must be a time format string just likes %s",
			key, `"`+strings.Join(samples, `" or "`)+`"`),
		Key:      key,
		KType:    ktype,
		CndVType: cndVType,
		Samples:  samples,
	}
}

// ---------------------------------------------------------
// 条件表达式中的正则表达式写法不正确
type ErrRePattern struct {
	error
	Key     string // 字段名称
	Pattern string // 正则表达式
}

func makeErrRePattern(key string, pattern string) ErrRePattern {
	return ErrRePattern{
		error:   fmt.Errorf("regexp pattern %q assigned to condition value for key %q is invalid", pattern, key),
		Key:     key,
		Pattern: pattern,
	}
}


// -----------------------------------------------------------------------------
// 筛选分页时的错误
// -----------------------------------------------------------------------------
// 排序字段不存在
type ErrNoOrderByKey struct {
	error
	OrderBy string
}

func makeErrNoOrderBy(key string) ErrNoOrderByKey {
	return ErrNoOrderByKey{
		error: fmt.Errorf("search object has no member named %q", key),
		OrderBy: key,
	}
}
