/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: util
@author: fanky
@version: 1.0
@date: 2021-04-02
**/

package fsky

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

// 判断 interface{} 包装值是否为 nil
// 注意：
//   var a *string
//   var aa any = a
//   fmt.Println(a == nil, aa == nil)  // true false
//   fmt.Println(IsNil(a), IsNil(aa))  // true true
func IsNil(v any) bool {
	rv := reflect.ValueOf(v)
	return !rv.IsValid() || (rv.Type().Kind() == reflect.Ptr && rv.IsNil())
}

// IfElse 模拟三目运算符
func IfElse[T any](b bool, left, right T) T {
	if b {
		return left
	}
	return right
}

// -------------------------------------------------------------------
// 获取函数的第一个返回值
func Ret1[T any](args ...any) T {
	return args[0].(T)
}

// 获取函数的第二个返回值
func Ret2[T any](args ...any) T {
	return args[1].(T)
}

// 获取函数的第三个返回值
func Ret3[T any](args ...any) T {
	return args[2].(T)
}

// -------------------------------------------------------------------
// 深拷贝对象
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return fmt.Errorf("encode src instance error: %v", err)
	}
	reader := bytes.NewReader(buf.Bytes())
	if err := gob.NewDecoder(reader).Decode(dst); err != nil {
		return fmt.Errorf("decode memory buffer error: %v", err)
	}
	return nil
}
