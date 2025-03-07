/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 对象查找
@author: fanky
@version: 1.0
@date: 2025-01-08
**/

package objsearchv2

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"fsky.pro/fssearch"

	"fsky.pro/fsreflect"
)

type T_MatchType = fssearch.T_MatchType

// -------------------------------------------------------------------
// Compare Handlers
// -------------------------------------------------------------------
type s_CompareHandlers struct {
	sync.RWMutex
	handlers map[fssearch.T_MatchType]func(any, string, *s_CmpValue) (bool, error)
}

func (this *s_CompareHandlers) getMemberValue(obj any, key string) (v reflect.Value, err error) {
	hasKey := false
	// 遍历时，优先选择子结构体的成员
	fsreflect.TrivalStructMembers(obj, false, func(info *fsreflect.S_TrivalStructInfo) bool {
		if info.IsBase {
			return true
		}
		tag := info.Field.Tag.Get("json")
		if tag != key {
			return true
		}
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
		err = makeErrNoCndMember(fsreflect.BaseRefType(reflect.TypeOf(obj)), key)
	}
	return
}

// ---------------------------------------------------------
// 包含
func (this *s_CompareHandlers) contain(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}
	if rv.Type().Kind() != reflect.String {
		// 只有字符串才能使用“包含”
		return false, makeErrUnsupportMatcher(key, "contain|not_contain")
	}
	// 这是为了解决 obj 的成员 key 是以下这种定义方法的情况：
	//   type T string
	//   var obj.key T
	mvalue := rv.Convert(refString).Interface().(string)
	v, ok := value.asString()
	if !ok {
		return false, makeErrCndValue(key, rv.Type(), value.vtype)
	}
	return strings.Contains(mvalue, v), nil
}

// 不包含
func (this *s_CompareHandlers) notContain(obj any, key string, value *s_CmpValue) (bool, error) {
	ok, err := this.contain(obj, key, value)
	return !ok, err
}

// 匹配
func (this *s_CompareHandlers) match(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}
	if rv.Type().Kind() != reflect.String {
		return false, makeErrUnsupportMatcher(key, "match|not_match")
	}
	mvalue := rv.Convert(refString).Interface().(string)
	v, ok := value.asString()
	if !ok {
		return false, makeErrCndValue(key, rv.Type(), value.vtype)
	}
	return mvalue == v, nil
}

// 不匹配
func (this *s_CompareHandlers) notMatch(obj any, key string, value *s_CmpValue) (bool, error) {
	ok, err := this.match(obj, key, value)
	return !ok, err
}

// 正则匹配
func (this *s_CompareHandlers) reMatch(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}
	if rv.Type().Kind() != reflect.String {
		return false, makeErrUnsupportMatcher(key, "re_match")
	}
	mvalue := rv.Convert(refString).Interface().(string)
	v, ok := value.asString()
	if !ok {
		return false, makeErrCndValue(key, rv.Type(), value.vtype)
	}
	re, err := regexp.Compile(v)
	if err != nil {
		return false, makeErrRePattern(key, v)
	}
	return re.Match([]byte(mvalue)), nil
}

// 等于
func (this *s_CompareHandlers) equal(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}
	if rv.Type().Kind() == reflect.String {
		mvalue := rv.Convert(refString).Interface().(string)
		v, ok := value.asString()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return mvalue == v, nil
	}

	switch rv.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, ok := value.asInt64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refInt64).Interface().(int64) == v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, ok := value.asUint64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refUint64).Interface().(uint64) == v, nil
	case reflect.Float32, reflect.Float64:
		v, ok := value.asFloat64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refFloat64).Interface().(float64) == v, nil
	case refTime.Kind():
		v, ok := value.asTime()
		if !ok {
			return false, makeErrCndTimeValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refTime).Interface().(time.Time).Equal(v), nil
	}
	return false, makeErrUnsupportMatcher(key, "equal|not_equal|large_then|large_equal|less_than|less_equal")
}

// 不等于
func (this *s_CompareHandlers) notEqual(obj any, key string, value *s_CmpValue) (bool, error) {
	ok, err := this.equal(obj, key, value)
	return !ok, err
}

// 大于
func (this *s_CompareHandlers) largeThan(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}

	switch rv.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, ok := value.asInt64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refInt64).Interface().(int64) > v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, ok := value.asUint64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refUint64).Interface().(uint64) > v, nil
	case reflect.Float32, reflect.Float64:
		v, ok := value.asFloat64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refFloat64).Interface().(float64) > v, nil
	case refTime.Kind():
		v, ok := value.asTime()
		if !ok {
			return false, makeErrCndTimeValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refTime).Interface().(time.Time).After(v), nil
	}
	return false, makeErrUnsupportMatcher(key, "equal|not_equal|large_then|large_equal|less_than|less_equal")
}

// 小于
func (this *s_CompareHandlers) lessThan(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}

	switch rv.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, ok := value.asInt64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refInt64).Interface().(int64) < v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, ok := value.asUint64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refUint64).Interface().(uint64) < v, nil
	case reflect.Float32, reflect.Float64:
		v, ok := value.asFloat64()
		if !ok {
			return false, makeErrCndValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refFloat64).Interface().(float64) < v, nil
	case refTime.Kind():
		v, ok := value.asTime()
		if !ok {
			return false, makeErrCndTimeValue(key, rv.Type(), value.vtype)
		}
		return rv.Convert(refTime).Interface().(time.Time).Before(v), nil
	}
	return false, makeErrUnsupportMatcher(key, "equal|not_equal|large_then|large_equal|less_than|less_equal")
}

// 大于等于
func (this *s_CompareHandlers) largeEqual(obj any, key string, value *s_CmpValue) (bool, error) {
	ok, err := this.largeThan(obj, key, value)
	if ok || err != nil {
		return ok, err
	}
	return this.equal(obj, key, value)
}

// 小于等于
func (this *s_CompareHandlers) lessEqual(obj any, key string, value *s_CmpValue) (bool, error) {
	ok, err := this.lessThan(obj, key, value)
	if ok || err != nil {
		return ok, err
	}
	return this.equal(obj, key, value)
}

// 在集合中
func (this *s_CompareHandlers) in(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}

	cvs, ok := value.asArray()
	if !ok {
		rt := reflect.MakeSlice(rv.Type(), 0, 0).Type()
		return false, makeErrCndValue(key, rt, value.vtype)
	}

	mvalue := rv.Convert(refInt64).Interface()
	switch rv.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for _, cv := range cvs {
			v, ok := cv.asInt64()
			if !ok {
				continue
			}
			if mvalue.(int64) == v {
				return true, nil
			}
		}
		return false, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		for _, cv := range cvs {
			v, ok := cv.asUint64()
			if !ok {
				continue
			}
			if mvalue.(uint64) == v {
				return true, nil
			}
		}
		return false, nil
	case reflect.Float32, reflect.Float64:
		for _, cv := range cvs {
			v, ok := cv.asFloat64()
			if !ok {
				continue
			}
			if mvalue.(float64) == v {
				return true, nil
			}
		}
		return false, nil
	case refTime.Kind():
		for _, cv := range cvs {
			v, ok := cv.asTime()
			if !ok {
				continue
			}
			if mvalue.(time.Time).Equal(v) {
				return true, nil
			}
		}
		return false, nil
	}
	return false, makeErrUnsupportMatcher(key, "in")
}

// 数组长度等于
func (this *s_CompareHandlers) arrlenEqual(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}

	length, ok := value.asInt64()
	if !ok {
		rt := reflect.MakeSlice(rv.Type(), 0, 0).Type()
		return false, makeErrCndValue(key, rt, value.vtype)
	}

	switch rv.Type().Kind() {
	case reflect.Slice, reflect.Array:
		return rv.Len() == int(length), nil
	}
	return false, makeErrUnsupportMatcher(key, "arrlen_equal|arrlen_less_than|arrlen_less_equal|arrlen_large_than|arrlen_large_equal")
}

// 数组长度小于
func (this *s_CompareHandlers) arrlenLessThan(obj any, key string, value *s_CmpValue) (bool, error) {
	rv, err := this.getMemberValue(obj, key)
	if err != nil {
		return false, err
	}

	length, ok := value.asInt64()
	if !ok {
		rt := reflect.MakeSlice(rv.Type(), 0, 0).Type()
		return false, makeErrCndValue(key, rt, value.vtype)
	}

	switch rv.Type().Kind() {
	case reflect.Slice, reflect.Array:
		return rv.Len() < int(length), nil
	}
	return false, makeErrUnsupportMatcher(key, "arrlen_equal|arrlen_less_than|arrlen_less_equal|arrlen_large_than|arrlen_large_equal")
}

// 数组长度小于等于
func (this *s_CompareHandlers) arrlenLessEqual(obj any, key string, value *s_CmpValue) (bool, error) {
	ok, err := this.arrlenEqual(obj, key, value)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return this.arrlenLessThan(obj, key, value)
}

// 数组长度大于
func (this *s_CompareHandlers) arrlenLargeThan(obj any, key string, value *s_CmpValue) (bool, error) {
	ok, err := this.arrlenLessEqual(obj, key, value)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

// 数组长度大于等于
func (this *s_CompareHandlers) arrlenLargeEqual(obj any, key string, value *s_CmpValue) (bool, error) {
	ok, err := this.arrlenLessThan(obj, key, value)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

var cmpHandlers *s_CompareHandlers

func init() {
	cmpHandlers = &s_CompareHandlers{
		handlers: map[T_MatchType]func(any, string, *s_CmpValue) (bool, error){},
	}
	cmpHandlers.handlers[fssearch.MT_Contain] = cmpHandlers.contain
	cmpHandlers.handlers[fssearch.MT_NoContain] = cmpHandlers.notContain
	cmpHandlers.handlers[fssearch.MT_Match] = cmpHandlers.match
	cmpHandlers.handlers[fssearch.MT_NotMatch] = cmpHandlers.notMatch
	cmpHandlers.handlers[fssearch.MT_Equal] = cmpHandlers.equal
	cmpHandlers.handlers[fssearch.MT_NoEqual] = cmpHandlers.notEqual
	cmpHandlers.handlers[fssearch.MT_Less] = cmpHandlers.lessThan
	cmpHandlers.handlers[fssearch.MT_LessEqual] = cmpHandlers.lessEqual
	cmpHandlers.handlers[fssearch.MT_Large] = cmpHandlers.largeThan
	cmpHandlers.handlers[fssearch.MT_LargeEqual] = cmpHandlers.largeEqual
	cmpHandlers.handlers[fssearch.MT_ReMatch] = cmpHandlers.reMatch
	cmpHandlers.handlers[fssearch.MT_InArray] = cmpHandlers.in
	cmpHandlers.handlers[fssearch.MT_ArrLenEqual] = cmpHandlers.arrlenEqual
	cmpHandlers.handlers[fssearch.MT_ArrLenLess] = cmpHandlers.arrlenLessThan
	cmpHandlers.handlers[fssearch.MT_ArrLenLessEqual] = cmpHandlers.arrlenLessEqual
	cmpHandlers.handlers[fssearch.MT_ArrLenLarge] = cmpHandlers.arrlenLargeThan
	cmpHandlers.handlers[fssearch.MT_ArrLenLargeEqual] = cmpHandlers.arrlenLargeEqual
}
