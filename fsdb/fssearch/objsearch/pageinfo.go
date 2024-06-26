/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 分页信息
@author: fanky
@version: 1.0
@date: 2023-10-01
**/

package objsearch

import (
	"math"
)

type S_PageInfo struct {
	Page     int `json:"page"`     // 当前页码
	PageSize int `json:"pageSize"` // 单页记录数
	Offset   int `json:"offset"`   // 当前页的起始索引
	Pages    int `json:"pages"`    // 总页数
	Total    int `json:"total"`    // 总记录
}

func NewPageInfo(total, page, pageSize int) *S_PageInfo {
	offset := 0
	pages := 1
	if pageSize > 0 {
		offset = (page - 1) * pageSize
		pages = int(math.Ceil(float64(total) / float64(pageSize)))
	} else {
		// pages 小于 0 时，返回全部记录
		page = 1
	}
	return &S_PageInfo{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
		Pages:    pages,
		Total:    total,
	}
}

func (this *S_PageInfo) IsFirstPage() bool {
	if this.PageSize > 0 {
		return this.Page == 1
	}
	return true
}

func (this *S_PageInfo) IsLastPage() bool {
	if this.PageSize > 0 {
		return this.Page == this.Pages
	}
	return true
}
