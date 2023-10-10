/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: reflect utils
@author: fanky
@version: 1.0
@date: 2021-02-19
**/

// ------------------------------------------------------------
// GetFieldValue（获取结构体的私有字段值）
// SetFieldValue（设置结构体的私有字段值）
//
// 获取包的全局变量值、调用包的私有函数、调用结构体的私有函数，不需要
// 反射，可以用 //go:linkname 指令
// ------------------------------------------------------------

package fsreflect

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
type s_FieldPathItem struct {
	name  string
	isKey bool
}

var _remain = regexp.MustCompile(`^[\d\w_]+`)
var _rekey = regexp.MustCompile(`^\[.*?\]`)

// 将深层路径切割成多个分段
func _splitFieldPath(str string) []*s_FieldPathItem {
	str = strings.TrimSpace(str)
	if str == "" {
		return []*s_FieldPathItem{}
	}
	items := make([]*s_FieldPathItem, 0)

	name := _remain.FindString(str)
	if name == "" {
		return nil
	}
	items = append(items, &s_FieldPathItem{name, false})
	str = str[len(name):]

	for len(str) > 0 {
		if str[0] == '.' {
			str = str[1:]
			name := _remain.FindString(str)
			if name == "" {
				return nil
			}
			items = append(items, &s_FieldPathItem{name, false})
			str = str[len(name):]
		} else if str[0] == '[' {
			key := _rekey.FindString(str)
			if key == "" {
				return nil
			}
			count := len(key)
			key = strings.TrimSpace(key[1 : count-1])
			items = append(items, &s_FieldPathItem{key, true})
			str = str[count:]
		} else {
			return nil
		}
	}
	return items
}

// ---------------------------------------------------------
// 获取结构体成员的指针
// vobj 必须是结构体指针
func getFieldPtr(vobj reflect.Value, fname string) (pfield reflect.Value, err error) {
	if vobj.Type().Kind() != reflect.Ptr {
		err = fmt.Errorf("obj must be pointer type")
		return
	}
	if vobj.IsNil() {
		err = fmt.Errorf("obj is nil")
		return
	}
	vobj = vobj.Elem()
	field := vobj.FieldByName(fname)
	if !field.IsValid() {
		err = fmt.Errorf("%v has no member named %q", vobj.Type(), fname)
		return
	}

	// 如果成员本身就是指针，则返回本身
	if field.Type().Kind() == reflect.Ptr {
		pfield = field
	} else if field.CanAddr() {
		pfield = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr()))
	} else {
		err = fmt.Errorf("member %q of %v is unaccessable", fname, vobj.Type())
	}
	return
}

// 根据路径，深度获取结构体成员的指针
// vobj 必须是非 nil 结构体指针，因为是私有函数，这里不作校验
func getDeepFieldPtr(vobj reflect.Value, fpath string) (pv reflect.Value, err error) {
	tobj := vobj.Type().Elem()
	pathItems := _splitFieldPath(fpath)
	if pathItems == nil {
		err = fmt.Errorf("error field path %q", fpath)
		return
	}
	path := ""
	pv = vobj
	for _, pitem := range pathItems {
		if !pv.IsValid() {
			err = fmt.Errorf("field path %q of %v is inavaliable", fpath, tobj)
			return
		}
		tv := pv.Type().Elem().Kind()
		if pv.IsNil() {
			err = fmt.Errorf("field path %q of %v is nil.", fpath, tobj)
			return
		}

		if pitem.isKey {
			path = path + fmt.Sprintf("[%s]", pitem.name)
			tv = pv.Type().Elem().Kind()
			if tv == reflect.Map {
				if pv, err = mapGetValuePtr(pv, pitem.name); err != nil {
					err = fmt.Errorf("get map value(%s) fail, %v", path, err)
					return
				}
			} else if tv == reflect.Slice {
				index, e := strconv.Atoi(pitem.name)
				if e != nil {
					err = fmt.Errorf(`invalid field path "%s" of %v. slice subscript must be an intager.`, path, tobj)
					return
				}
				if pv, err = sliceGetItemPtr(pv, index); err != nil {
					err = fmt.Errorf("get slice element(%s) fail, %v", path, err)
					return
				}
			} else if tv == reflect.Array {
				index, e := strconv.Atoi(pitem.name)
				if e != nil {
					err = fmt.Errorf(`invalid field path "%s" of %v. array subscript must be an intager.`, path, tobj)
					return
				}
				if pv, err = arrayGetItemPtr(pv, index); err != nil {
					err = fmt.Errorf("get array element(%s) fail, %v", path, err)
					return
				}
			} else {
				err = fmt.Errorf(`field path "%s" of %v is not exists`, path, tobj)
				return
			}
		} else {
			path = path + fmt.Sprintf(".%s", pitem.name)
			if tv == reflect.Struct {
				if pv, err = getFieldPtr(pv, pitem.name); err != nil {
					err = fmt.Errorf("get field %q of %v fail, error: %v", fpath, tobj, err)
					return
				}
			}
		}
	}
	if !pv.IsValid() || !pv.Elem().CanInterface() {
		err = fmt.Errorf("field path %q of %v is not accessable", fpath, tobj)
	}
	return
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 获取结构体字段值，包括私有字段
func GetFieldValue(obj interface{}, fname string) (fv interface{}, err error) {
	if obj == nil {
		err = errors.New("obj argument must't be nil")
		return
	}

	v := reflect.ValueOf(obj)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Type().Kind() != reflect.Struct {
		err = errors.New("obj argument must be a struct object or struct pointer")
		return
	}

	field := v.FieldByName(fname)
	if !field.IsValid() {
		err = newFieldError(v.Type(), fname)
		return
	}

	if field.CanInterface() {
		fv = field.Interface()
		return
	}

	// 成员为指针类型
	if field.Type().Kind() == reflect.Ptr {
		if field.IsNil() {
			return
		}
		if field.Elem().CanAddr() {
			up := unsafe.Pointer(field.Elem().UnsafeAddr())
			vv := reflect.NewAt(field.Elem().Type(), up)
			fv = vv.Interface()
			return
		}
	} else if field.CanAddr() {
		up := unsafe.Pointer(field.UnsafeAddr())
		vv := reflect.NewAt(field.Type(), up)
		fv = vv.Elem().Interface()
		return
	}
	err = fmt.Errorf("member %q of %v is unaccessable", fname, v.Type())
	return
}

// 设置结构体字段值，包括私有字段
// hardcvr 表示，如果类型不一致是否试图进行强转
func SetFieldValue(obj interface{}, fname string, fv interface{}, hardcvr bool) error {
	if obj == nil {
		return errors.New("object must't be nil.")
	}

	tobj := reflect.TypeOf(obj)
	vobj := reflect.ValueOf(obj)
	if tobj.Kind() == reflect.Ptr {
		tobj = tobj.Elem()
		vobj = vobj.Elem()
	}
	if tobj.Kind() != reflect.Struct {
		return errors.New("obj argument must be a struct object or struct pointer.")
	}

	// 查找对应名称的成员
	var field reflect.Value
	for i := 0; i < vobj.NumField(); i++ {
		if tobj.Field(i).Name == fname {
			field = vobj.Field(i)
			break
		}
	}

	// 域名不存在
	if !field.IsValid() {
		return newFieldError(tobj, fname)
	}
	tfield := field.Type()

	// 新值
	var v reflect.Value
	if fv == nil {
		v = reflect.Zero(tfield)
	} else {
		v = reflect.ValueOf(fv)
	}
	tv := v.Type()

	// 指定的域类型为指针类型
	if tfield.Kind() == reflect.Ptr {
		if tv != tfield { // 类型不匹配
			return newValueError(tobj, fname, tfield, tv)
		}
		if field.CanSet() {
			field.Set(v)
		} else if field.CanAddr() {
			upp := unsafe.Pointer(field.UnsafeAddr()) // 指针的指针
			ppv := reflect.NewAt(field.Type(), upp)   // 创建指针的指针对象（该指针的指针，指向 field 原来的位置）
			ppv.Elem().Set(v)                         // 将指针的指针，指向的位置（即 field 的内存位置）修改为新的值
		} else {
			return fmt.Errorf("member %q of %v is unaccessable", fname, tobj)
		}
	} else {
		if tv != tfield {
			// 不强转
			if !hardcvr {
				return newValueError(tobj, fname, tfield, tv)
			}
			// 强制类型转换
			if cv, ok := hardConvert(fv, field.Type()); ok {
				v = cv
			} else {
				return newValueError(tobj, fname, tfield, tv)
			}
		}
		if field.CanSet() {
			field.Set(v)
		} else if field.CanAddr() {
			up := unsafe.Pointer(field.UnsafeAddr()) // 获取 field 的指针（注意：不能把 field.UnsafeAddr() 存放到临时变量，否则可能会被 GC）
			pv := reflect.NewAt(field.Type(), up)    // 创建一个新的指针，指向原来 field 的位置
			pv.Elem().Set(v)
		} else {
			return fmt.Errorf("member %q of %v is unaccessable", fname, tobj)
		}
	}
	return nil
}

// -------------------------------------------------------------------
// 获取子孙成员值
// path 用 . 分隔父子成员名称
// obj 必须是结构体指针
func GetDeepFieldValue(obj interface{}, fpath string) (v interface{}, err error) {
	if obj == nil {
		err = errors.New("struct object must not be a nil value")
		return
	}
	tobj := reflect.TypeOf(obj)
	vobj := reflect.ValueOf(obj)
	if tobj.Kind() != reflect.Ptr {
		err = errors.New("obj argument must be a pointer of struct object")
		return
	}
	if tobj.Elem().Kind() != reflect.Struct {
		err = errors.New("obj argument must be a pointer of struct object")
		return
	}

	pv, e := getDeepFieldPtr(vobj, fpath)
	if e == nil {
		v = pv.Elem().Interface()
	} else {
		err = e
	}
	return
}

// 设置子孙成员值
// obj 必须是结构体对象指针
// path 用 . 分隔，map、slice、array 用 [] 作为下标括号
func SetDeepFieldValue(obj interface{}, fpath string, value interface{}) error {
	if obj == nil {
		return errors.New("obj argument is not allow to be a nil value")
	}
	tobj := reflect.TypeOf(obj)
	vobj := reflect.ValueOf(obj)
	if tobj.Kind() != reflect.Ptr {
		return errors.New("obj argument must be a pointer of struct object")
	}
	tobj = tobj.Elem()
	if tobj.Kind() != reflect.Struct {
		return errors.New("obj argument must be a pointer of struct object")
	}

	pv, err := getDeepFieldPtr(vobj, fpath)
	if err != nil {
		return err
	}
	tv := pv.Elem().Type()

	var vvalue reflect.Value
	if value == nil {
		vvalue = reflect.New(tv)
	} else if cv, ok := hardConvert(value, tv); ok {
		vvalue = cv
	} else {
		return fmt.Errorf("value type %v can't convert to the field(%s) type %v of %v", vvalue.Type(), fpath, tv, tobj)
	}
	if !pv.Elem().CanSet() {
		return fmt.Errorf("field %q of %v is unaccessable", fpath, tobj)
	}
	pv.Elem().Set(vvalue)
	return nil
}

// -------------------------------------------------------------------
// 遍历结构体成员，包括父结构体的成员
// 如果参数 f 返回 false，则停止遍历
// 函数参数 f 的参数：
//	S_TrivalStructInfo.StructType：
//		遍历过程中，当前结构体的类型
//  S_TrivalStructInfo::StructValue：
//		历过程中，当前结构体对象
//  S_TrivalStructInfo::Field：
//		遍历过程中，当前成员域
//  S_TrivalStructInfo::FieldValue：
//		遍历过程中，当前成员的值，FieldValue.IsValid()、FieldValue.Type()、FieldValue.IsNil()，都是不确定的
// 提示：
//  参数 v 可以传入任何结构体的 nil 值，但是如果传入 nil，则遍历过程中，f 的参数 a1.IsValid() 和 a3.IsValid() 都是 false
// 示例：
//  type A struct {
//		member1 string
//  }
//  type B struct {
//		member2 int
//	}
//	type C struct {
//		A
//		*B
//		member3 uint64
//  }
//
//  var c *C = nil
//  TrivalFields(c, func(*S_TrivalStructInfo)bool{
//      return true
//  })
// -------------------------------------------------------------------
type S_TrivalStructInfo struct {
	StructType  reflect.Type
	StructValue reflect.Value
	Field       reflect.StructField
	FieldValue  reflect.Value
}

func TrivalStructMembers(v any, f func(*S_TrivalStructInfo) bool) {
	rt := reflect.TypeOf(v)
	if rt == nil { return }
	rv := reflect.ValueOf(v)
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		if rv.IsValid() || !rv.IsNil() {
			rv = rv.Elem()
		} else {
			rv = reflect.Value{}
		}
	}

	var trivalStruct func(reflect.Type, reflect.Value) bool
	trivalStruct = func(rt reflect.Type, rv reflect.Value) bool {
		if rt == nil                   { return true }
		if rt.Kind() != reflect.Struct { return true }
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			vfield := reflect.ValueOf(nil)
			if rv.IsValid() {
				vfield = rv.Field(i)
			}
			if !field.Anonymous {
				info := new(S_TrivalStructInfo)
				info.StructType = rt
				info.StructValue = rv
				info.Field = field
				info.FieldValue = vfield
				if f(info) { continue }
				return false
			}
			// 匿名结构体
			tfield := field.Type
			for tfield.Kind() == reflect.Ptr {
				tfield = tfield.Elem()
				if !vfield.IsValid() || vfield.IsNil() {
					vfield = reflect.ValueOf(nil)
				} else {
					vfield = vfield.Elem()
				}
			}
			// 继承结构体，继续往上层遍历
			if tfield.Kind() == reflect.Struct {
				if !trivalStruct(tfield, vfield) {
					return false
				}
			}
		}
		return true
	}
	trivalStruct(rt, rv)
}

// -------------------------------------------------------------------
// 浅拷贝结构对象，src 必须为结构体指针
func CopyStructObject(src interface{}) (dst interface{}, err error) {
	tsrc := reflect.TypeOf(src)
	if src == nil || tsrc == nil {
		err = errors.New("src object is not allow to be a nil value")
		return
	}
	if tsrc.Kind() != reflect.Ptr {
		err = errors.New("src object must be an object pointer")
		return
	}
	tsrc = tsrc.Elem()
	// 限制类型必须为结构体
	if tsrc.Kind() != reflect.Struct {
		err = fmt.Errorf("type of src object must be a pointer of struct")
		return
	}

	vsrc := reflect.ValueOf(src).Elem()
	vdst := reflect.New(tsrc)
	dst = vdst.Interface()
	vdst = vdst.Elem()

	// 将 src 所有成员复制给 dst
	for i := 0; i < tsrc.NumField(); i++ {
		vfSrc := vsrc.Field(i)
		vfDst := vdst.Field(i)
		if vfDst.CanSet() {
			vfDst.Set(vfSrc)
		} else {
			upSrc := unsafe.Pointer(vfSrc.UnsafeAddr())
			pfSrc := reflect.NewAt(vfSrc.Type(), upSrc)
			v := pfSrc.Elem().Interface()

			upDst := unsafe.Pointer(vfDst.UnsafeAddr())
			pfDst := reflect.NewAt(vfDst.Type(), upDst)
			pfDst.Elem().Set(reflect.ValueOf(v))
		}
	}
	return
}
