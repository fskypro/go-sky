/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 调用对象基类
@author: fanky
@version: 1.0
@date: 2025-05-28
**/

package unixrpc

import (
	"fmt"
	"reflect"
	"time"
)

// -------------------------------------------------------------------
// 用于限制服务器端接收器
// -------------------------------------------------------------------
type I_Receiver[Arg, Ret any] interface {
	OnCalled(*Arg, *Ret) error
}

// -------------------------------------------------------------------
// 用于匹配 Unix 协议服务器/客户端
// -------------------------------------------------------------------
type I_UnixServer interface {
	Register(string, any)
}

type I_UnixCleint interface {
	Call(time.Duration, string, any, any) error
}

// -------------------------------------------------------------------
// 远程调用包装器
// Recv 为消息接收者类型
// Arg  为消息参数类型
// Ret  为调用返回值类型
// -------------------------------------------------------------------
type S_CallWrap[Proxy, Arg, Ret any] struct {
	recvName string
	method   string
}

// 创建服务器端绑定对象
func NewCallWrap[Proxy, Arg, Ret any]() *S_CallWrap[Proxy, Arg, Ret] {
	var proxy Proxy
	t := reflect.TypeOf(proxy)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return &S_CallWrap[Proxy, Arg, Ret]{
		recvName: fmt.Sprintf("%s::%s", t.PkgPath(), t.Name()),
	}
}

func (this *S_CallWrap[Proxy, Arg, Ret]) RegisterTo(server I_UnixServer, recv I_Receiver[Arg, Ret]) {
	server.Register(this.recvName, recv)
}

func (this *S_CallWrap[Proxy, Arg, Ret]) CallWith(client I_UnixCleint, timeout time.Duration, arg *Arg, ret *Ret) error {
	method := fmt.Sprintf("%s.OnCalled", this.recvName)
	return client.Call(timeout, method, arg, ret)
}
