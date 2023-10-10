/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 排序器
@author: fanky
@version: 1.0
@date: 2023-10-01
**/

package objsearch

import (
	"reflect"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// sorter
// -----------------------------------------------------------------------------
type s_Sorter[T any] struct {
	searchArg *S_SearchArg[T]
	items     []T
}

func newSorter[T any](arg *S_SearchArg[T], items []T) *s_Sorter[T] {
	return &s_Sorter[T]{
		searchArg: arg,
		items:     items,
	}
}

func (self s_Sorter[T]) Len() int {
	return len(self.items)
}

func (self s_Sorter[T]) Less(i, j int) bool {
	v1, _ := cmpHandlers.getMemberValue(self.items[i], self.searchArg.OrderBy)
	v2, _ := cmpHandlers.getMemberValue(self.items[j], self.searchArg.OrderBy)
	if v1.Type().Kind() == reflect.String {
		s1 := v1.Interface().(string)
		s2 := v2.Interface().(string)
		if self.searchArg.Desc > 0 {
			return strings.Compare(s2, s1) < 0
		}
		return strings.Compare(s1, s2) < 0
	} else if v1.CanConvert(reflect.TypeOf(1.0)) {
		n1 := v1.Convert(reflect.TypeOf(1.0)).Interface().(float64)
		n2 := v2.Convert(reflect.TypeOf(1.0)).Interface().(float64)
		if self.searchArg.Desc > 0 {
			return n2 < n1
		}
		return  n1 < n2
	} else if v1.CanConvert(reflect.TypeOf(time.Now())) {
		t1 := v1.Convert(reflect.TypeOf(time.Now())).Interface().(time.Time)
		t2 := v2.Convert(reflect.TypeOf(time.Now())).Interface().(time.Time)
		if self.searchArg.Desc > 0 {
			return t2.Before(t1)
		}
		return t1.Before(t2)
	}
	return false
}

func (self s_Sorter[T]) Swap(i, j int) {
	self.items[i], self.items[j] = self.items[j], self.items[i]
}
