/**
@copyright: fantasysky 2016
@brief: 定义
@author: fanky
@version: 1.0
@date: 2019-05-30
**/

package freejson

type JType int

// 数据类型
const (
	TNull    = 0 // null
	TObject  = 1 // Json 对象
	TList    = 2 // 列表
	TString  = 3 // 字符串类型
	TBool    = 4 // 布尔类型
	TInt64   = 5 // int64 类型
	TUInt64  = 6 // uint64 类型
	TFloat64 = 7 // foat64 类型
)

var typeNames = map[JType]string{
	TNull:    "null",
	TObject:  "object",
	TList:    "list",
	TString:  "string",
	TBool:    "bool",
	TInt64:   "int",
	TUInt64:  "uint",
	TFloat64: "float",
}

// 遍历 Object 函数
// 如果返回 false 则退出遍历
type F_KeyValue func(string, I_Value) bool

// 遍历 List 函数
// 如果返回 false 则退出遍历
type F_Elem func(int, I_Value) bool
