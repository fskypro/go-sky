/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 对象查找
@author: fanky
@version: 1.0
@date: 2023-10-01
**/

package objsearch

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"fsky.pro/fsreflect"
)

// -------------------------------------------------------------------
// Compare Handlers
// -------------------------------------------------------------------
type s_CompareHandlers struct {
	sync.RWMutex
	handlers map[string]func(any, string, string) (bool, error)
}

func (this *s_CompareHandlers) getMemberValue(obj any, key string) (v reflect.Value, err error) {
	hasKey := false
	fsreflect.TrivalStructMembers(obj, func(info *fsreflect.S_TrivalStructInfo) bool {
		tag := info.Field.Tag.Get("json")
		if tag != key { return true }
		rv := info.FieldValue
		hasKey = true
		if !rv.IsValid() {
			err = fmt.Errorf("key %q's value in type %v object is not valid", key, info.Field.Type)
			return false
		}
		v = reflect.NewAt(info.Field.Type, unsafe.Pointer(rv.UnsafeAddr())).Elem()
		return false
	})
	if !hasKey {
		err = fmt.Errorf("no member key %q in object %q", key, reflect.TypeOf(obj))
	}
	return
}

// ---------------------------------------------------------
// 包含
func (this *s_CompareHandlers) contain(obj any, key, value string) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil { return false, err }
	if rv.Type().Kind() != reflect.String {
		return false, fmt.Errorf("the key(%q)'s type %v in %v is not support contain method", key, rv.Type(), reflect.TypeOf(obj))
	}
	mvalue := rv.Convert(reflect.TypeOf("")).Interface().(string)
	return strings.Contains(mvalue, value), nil
}

// 不包含
func (this *s_CompareHandlers) notContain(obj any, key, value string) (bool, error) {
	ok, err := this.contain(obj, key, value)
	return !ok, err
}

// 匹配
func (this *s_CompareHandlers) match(obj any, key, value string) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil { return false, err }
	if rv.Type().Kind() != reflect.String {
		return false, fmt.Errorf("the key(%q)'s type %v in %v is not support contain method", key, rv.Type(), reflect.TypeOf(obj))
	}
	mvalue := rv.Convert(reflect.TypeOf("")).Interface().(string)
	return mvalue == value, nil
}

// 不匹配
func (this *s_CompareHandlers) notMatch(obj any, key, value string) (bool, error) {
	ok, err := this.match(obj, key, value)
	return !ok, err
}

// 等于
func (this *s_CompareHandlers) equal(obj any, key, value string) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil { return false, err }
	if rv.Type().Kind() == reflect.String {
		mvalue := rv.Convert(reflect.TypeOf("")).Interface().(string)
		return mvalue == value, nil
	}
	if rv.CanConvert(reflect.TypeOf(0.1)) {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil { return false, nil }
		return rv.Convert(reflect.TypeOf(0.1)).Equal(reflect.ValueOf(v)), nil
	}
	if rv.CanConvert(reflect.TypeOf(true)) {
		v, err := strconv.ParseBool(value)
		if err != nil { return false, nil }
		b := rv.Convert(reflect.TypeOf(true)).Interface().(bool)
		return b == v, nil
	}
	if rv.CanConvert(reflect.TypeOf(time.Now())) {
		t, err := time.Parse("2006-01-02 15:04-05", value)
		if err != nil {
			return false, fmt.Errorf("value of key %q must be a time format string just likes '2023-10-12 14:23:12'", key)
		}
		return rv.Interface().(time.Time).Equal(t), nil
	}
	return false, nil
}

// 不等于
func (this *s_CompareHandlers) notEqual(obj any, key, value string) (bool, error) {
	ok, err := this.equal(obj, key, value)
	return !ok, err
}

// 大于
func (this *s_CompareHandlers) largeThan(obj any, key, value string) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil { return false, err }
	if rv.CanConvert(reflect.TypeOf(0)) {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false, nil
		}
		return rv.Interface().(float64) > v, nil
	}
	if rv.CanConvert(reflect.TypeOf(time.Now())) {
		t, err := time.Parse("2006-01-02 15:04-05", value)
		if err != nil {
			return false, fmt.Errorf("value of key %q must be a time format string just likes '2023-10-12 14:23:12'", key)
		}
		return rv.Interface().(time.Time).After(t), nil
	}
	if rv.Type().Kind() == reflect.String {
		mvalue := rv.Convert(reflect.TypeOf("")).Interface().(string)
		return strings.Compare(mvalue, value) > 0, nil
	}
	return false, nil
}

// 小于
func (this *s_CompareHandlers) lessThan(obj any, key, value string) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil { return false, err }
	if rv.CanConvert(reflect.TypeOf(0)) {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false, nil
		}
		return rv.Interface().(float64) < v, nil
	}
	if rv.CanConvert(reflect.TypeOf(time.Now())) {
		t, err := time.Parse("2006-01-02 15:04-05", value)
		if err != nil {
			return false, fmt.Errorf("value of key %q must be a time format string just likes '2023-10-12 14:23:12'", key)
		}
		return rv.Interface().(time.Time).Before(t), nil
	}
	if rv.Type().Kind() == reflect.String {
		mvalue := rv.Convert(reflect.TypeOf("")).Interface().(string)
		return strings.Compare(mvalue, value) < 0, nil
	}
	return false, nil
}

// 大于等于
func (this *s_CompareHandlers) largeEqual(obj any, key, value string) (bool, error) {
	ok, err := this.largeThan(obj, key, value)
	if ok || err != nil { return ok, err }
	return this.equal(obj, key, value)
}

// 小于等于
func (this *s_CompareHandlers) lessEqual(obj any, key, value string) (bool, error) {
	ok, err := this.lessThan(obj, key, value)
	if ok || err != nil { return ok, err }
	return this.equal(obj, key, value)
}

var cmpHandlers *s_CompareHandlers

func init() {
	cmpHandlers = &s_CompareHandlers{
		handlers: map[string]func(any, string, string) (bool, error){},
	}
	cmpHandlers.handlers["match"] = cmpHandlers.match
	cmpHandlers.handlers["not_match"] = cmpHandlers.notMatch
	cmpHandlers.handlers["contain"] = cmpHandlers.contain
	cmpHandlers.handlers["not_contain"] = cmpHandlers.notContain
	cmpHandlers.handlers["equal"] = cmpHandlers.equal
	cmpHandlers.handlers["not_equal"] = cmpHandlers.notEqual
	cmpHandlers.handlers["large_than"] = cmpHandlers.largeThan
	cmpHandlers.handlers["less_than"] = cmpHandlers.lessThan
	cmpHandlers.handlers["large_equal"] = cmpHandlers.largeEqual
	cmpHandlers.handlers["less_equal"] = cmpHandlers.lessEqual
}
