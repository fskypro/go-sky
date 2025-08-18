/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 搜索参数
@author: fanky
@version: 1.0
@date: 2025-01-08
**/

package objsearchv2

import (
	"sort"

	"fsky.pro/fssearch"

	"fsky.pro/fsreflect"
)

// -------------------------------------------------------------------
// SearchArg
// -------------------------------------------------------------------
type S_SearchArg[T any] struct {
	*fssearch.S_SearchArg
	cnd    i_Cnd
	parsed bool
}

func NewSearchArg[T any]() *S_SearchArg[T] {
	return &S_SearchArg[T]{
		S_SearchArg: fssearch.NewSearchArg(),
	}
}

// 可能返回：
// -> fssearch.ParseCndError
// -> fssearch.InvalidMatchError
// -> fssearch.InvalidOrderByError
func (this *S_SearchArg[T]) parse() error {
	if this.parsed {
		return nil
	}

	if this.Page < 1 {
		this.Page = 1
	}

	if this.Cnds == nil {
		return nil
	}
	cnd, err := parseCnd(this.Cnds)
	if err != nil {
		return err
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
		return fssearch.InvalidOrderByError{FieldName: this.OrderBy}
	}
	this.parsed = true
	return nil
}

// 检查对象是否符合当前搜索条件，可能会产生以下错误：
// -> fssearch.ParseCndError       ：解释条件表达式失败
// -> fssearch.InvalidCndFieldError：条件字段不存在
// -> fssearch.InvalidMatchError   ：条件表达式中，某字段采用 了无效匹配器
// -> fssearch.InvalidOrderByError ：order by 指定的字段不存在
// -> fssearch.RePatternError      ：条件中的正则表达式错误
func (this *S_SearchArg[T]) Check(obj *T) (bool, error) {
	if !this.parsed {
		if err := this.parse(); err != nil {
			return false, err
		}
	}
	if this.cnd == nil {
		return true, nil
	}
	return this.cnd.compare(obj)
}

func (this *S_SearchArg[T]) Filter(items []*T) (*fssearch.S_PageInfo, []*T) {
	if this.OrderBy != "" {
		sort.Sort(newSorter(this, items))
	}
	total := len(items)
	pageInfo := fssearch.NewPageInfo(total, this.Page, this.PageSize)
	start := pageInfo.Offset
	if start >= total {
		return pageInfo, []*T{}
	}

	if pageInfo.PageSize > 0 {
		end := start + pageInfo.PageSize
		if end < total {
			return pageInfo, items[start:end]
		}
	}
	return pageInfo, items[start:]
}
