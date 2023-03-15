/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: fstype utils
@author: fanky
@version: 1.0
@date: 2022-10-20
**/

package fstype

// 判断传入的值是否是模板参数中的类型
func IsType[T any](v any) (istype bool) {
	defer func() { recover() }()
	var tmp T = v.(T)
	func(T) {}(tmp)
	istype = true
	return
}

// 布尔型转换为整型
func BoolToNumber[T T_Number](v bool) T {
	if v {
		return T(1)
	}
	return T(0)
}
