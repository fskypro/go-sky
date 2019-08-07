/**
* @copyright: 2016 fantasysky
* @brief: 公共数据结构
* @author: fanky
* @version: 1.0
* @date: 2018-08-31
**/

// 服务器和客户端共同模块
package fsrpc

import "reflect"

// -------------------------------------------------------------------
// 全局变量
// -------------------------------------------------------------------
const (
	DefaultHTTPPath = "/fsrpc"                    // HTTP 默认访问路径
	ConnectedText   = "200 Connected to Go FSRPC" // 链接成功消息串
)

// -------------------------------------------------------------------
// 消息结构
// -------------------------------------------------------------------
// 请求头
type S_ReqHeader struct {
	ServiceName string // 请求服务对象名
	MethodName  string // 请求的方法名称
	ReqID       uint64 // 请求序号，用于匹配请求与回复之间的对应关系，只在客户端用到
}

// 回复头
type S_RspHeader struct {
	ServiceName string // 请求服务对象名
	MethodName  string // 请求的方法名称
	ReqID       uint64 // 回复处理序列号，用于匹配请求与回复之间的对应关系，只在客户端用到
	Error       string // 远程调用函数返回的错误信息
	Fail        string // 调用失败错误
}

// -------------------------------------------------------------------
// 空参数
// -------------------------------------------------------------------
// 空传入参数类型
// 如果一个远程调用不需要传入参数，则服务器定义为该类型
type EArg struct{}

// 空传出参数类型
// 如果一个远程调用不需要传出参数，则服务器定义为该类型
type EReply struct{}

// 判断 arg 是否是 EArg 或 *EArg 类型
func IsEmptyArg(arg interface{}) bool {
	if arg == nil {
		return false
	}
	eargt := reflect.TypeOf(&EArg{})
	argt := reflect.TypeOf(arg)
	if argt.Kind() == reflect.Ptr {
		return argt == eargt
	}
	return argt == eargt.Elem()
}

// 判断 reply 是否是 EReply 或 *EReply 类型
func IsEmptyReply(reply interface{}) bool {
	if reply == nil {
		return false
	}
	ereplyt := reflect.TypeOf(&EReply{})
	replyt := reflect.TypeOf(reply)
	if replyt.Kind() == reflect.Ptr {
		return replyt == ereplyt
	}
	return replyt == ereplyt.Elem()
}
