/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 搜索参数
@author: fanky
@version: 1.0
@date: 2023-10-01
**/

package objsearch

import (
	"errors"
	"fmt"
	"sort"

	"fsky.pro/fsreflect"
)

type i_Config interface {
	// 获取匹配处理器：
	// 匹配(字符串)    ："match"
	// 不匹配(字符串)  ："not_match"
	// 包含(字符串)    ："contain"
	// 不包含(字符串)  ："not_contain"
	// 正则匹配(字符串): "re_match"
	// 等于            ："equal"
	// 不等于          ："not_equal"
	// 大于            ："large_than"
	// 小于            ："less_than"
	// 大于等于        ："large_equal"
	// 小于等于        ："less_equal"
	// 在集合内        ："in"
	// 数组长度等于    ："arrlen_equal"
	// 数组长度小于    ："arrlen_less_than"
	// 数组长度小于等于："arrlen_less_equal"
	// 数组长度大于    ："arrlen_large_than"
	// 数组长度大于等于："arrlen_large_equal"
	GetMatchHandler(match string) string

	// 获取默认单页数量
	DefaultPageSize() int
}

// -------------------------------------------------------------------
// SearchArg
// -------------------------------------------------------------------
type S_SearchArg[T any] struct {
	Cnds     []any  `json:"cnds"`     // 用户自定义的原始搜索条件
	Page     int    `json:"page"`     // 要查询的页码（第一页为 1）
	PageSize int    `json:"pageSize"` // 每页最大数量
	OrderBy  string `json:"orderBy"`  // 排序字段
	Desc     int8   `json:"desc"`     // 是否倒序排列

	cnd    i_Cnd
	parsed bool
	conf   i_Config
}

func NewSearchArg[T any](conf i_Config) *S_SearchArg[T] {
	return &S_SearchArg[T]{
		Cnds:     make([]any, 0),
		Page:     1,
		PageSize: conf.DefaultPageSize(),
		conf:     conf,
	}
}

// 可能返回错误：
//
//	ErrNoOrderByKey: 搜索对象没有 orderby 字段
//	error          ：其他错误
func (this *S_SearchArg[T]) Parse() error {
	if this.parsed {
		return nil
	}
	this.parsed = true

	if this.Page < 1 {
		this.Page = 1
	}

	if this.Cnds == nil {
		return nil
	}
	cnd, err := parseCnd(this.conf, this.Cnds)
	if err != nil {
		return fmt.Errorf("parse serch conditions fail, %v", err)
	}
	this.cnd = cnd
	if this.OrderBy == "" {
		return nil
	}

	var obj T
	hasKey := false
	fsreflect.TrivalStructMembers(obj, false, func(info *fsreflect.S_TrivalStructInfo) bool {
		if info.IsBase {
			return true
		}
		tag := info.Field.Tag.Get("json")
		if tag == this.OrderBy {
			hasKey = true
			return false
		}
		return true
	})
	if !hasKey {
		return makeErrNoOrderBy(this.OrderBy)
	}
	return nil
}

// 检查对象是否符合当前搜索条件，可能会产生以下错误：
//
//	ErrUnsupportMatch ：字段不支持条件表达式中指定的匹配方式
//	ErrCndValue       ：条件表达式的值类型错误，不能与字段值进行比较
//	ErrCndTimeValue   ：条件表达式中，比较值使用的时间格式不正确
//	ErrRePattern      ：条件表达式中，要参与比较的正则表达式不正确
func (this *S_SearchArg[T]) Check(obj any) (bool, error) {
	if !this.parsed {
		return false, errors.New("search argument is not parsed")
	}
	if this.cnd == nil {
		return true, nil
	}
	return this.cnd.compare(obj)
}

func (this *S_SearchArg[T]) Filter(items []T) (*S_PageInfo, []T) {
	if this.OrderBy != "" {
		sort.Sort(newSorter(this, items))
	}
	total := len(items)
	pageInfo := NewPageInfo(total, this.Page, this.PageSize)
	start := pageInfo.Offset
	if start >= total {
		return pageInfo, []T{}
	}

	if pageInfo.PageSize > 0 {
		end := start + pageInfo.PageSize
		if end < total {
			return pageInfo, items[start:end]
		}
	}
	return pageInfo, items[start:]
}
