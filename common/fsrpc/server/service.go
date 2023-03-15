/**
* @copyright: fantasysky 2016
* @brief: 服务器可注册处理器解释和调用
* @author: fanky
* @version: 1.0
* @date: 2018-09-04
**/

package server

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"fsky.pro/fserror"
	"fsky.pro/fslog"
	"fsky.pro/fsky"

	. "fsky.pro/fsrpc"
)

// -----------------------------------------------------------------------------
// 可调用方法信息
// 可远程调用的方法：
//   1 方法名称必须以 "_rpc" 结尾
//   2 必须有且只有两个参数（第一个参数为请求参数；第二个参数为作为返回结果）
//   3 方法必须有且只有一个返回值，并且返回类型为 error
// -----------------------------------------------------------------------------
type s_MethodInfo struct {
	sync.Mutex
	method    reflect.Method // 请求方法
	argType   reflect.Type   // 方法参数类型
	replyType reflect.Type   // 回复客户端类型
	numCalls  uint64         // 客户端请求次数
}

// -----------------------------------------------------------------------------
// service inner methods
// -----------------------------------------------------------------------------
// 获取参数值
func (m *s_MethodInfo) _getArgValue(codec S_ServerCodec, req *S_ReqHeader) (argv reflect.Value, err error) {
	argIsPtr := false
	if m.argType.Kind() == reflect.Ptr { // 指针类型
		argv = reflect.New(m.argType.Elem()) // 则将指针类型转换为其指向的值类型
		argIsPtr = true
	} else {
		argv = reflect.New(m.argType) // 非指针类型
	}
	if err = codec.ReadRequestArg(req, argv.Interface()); err != nil {
		err = fmt.Errorf("read request of %s.%s body fail: %s\n", req.ServiceName, req.MethodName, err.Error())
		return
	}
	if !argIsPtr { // 如果是非指针类型，则将获取到的值转换为非指针类型
		argv = argv.Elem() // 相当于还原为指针类型
	}
	return
}

// 创建一个返回值实例
func (m *s_MethodInfo) _newReply() reflect.Value {
	reply := reflect.New(m.replyType.Elem())
	switch m.replyType.Elem().Kind() {
	case reflect.Map:
		reply.Elem().Set(reflect.MakeMap(m.replyType.Elem()))
	case reflect.Slice:
		reply.Elem().Set(reflect.MakeSlice(m.replyType.Elem(), 0, 0))
	}
	return reply
}

// -----------------------------------------------------------------------------
// 请求服务(每个服务对应一个注册)
// -----------------------------------------------------------------------------
type s_Service struct {
	name     string                   // 服务名称（可以在注册处理对象时传入）
	rcvr     reflect.Value            // 消息接收器（处理消息的处理器对象）
	rcvrType reflect.Type             // 消息接收器类（处理消息的处理器结构）
	methods  map[string]*s_MethodInfo // 消息处理器中的所以 public 方法
}

// -------------------------------------------------------------------
// service package methods
// -------------------------------------------------------------------
// 获取服务传入参数
func (svrc *s_Service) getArgValue(codec S_ServerCodec, req *S_ReqHeader) (argv reflect.Value, err error) {
	methodInfo, ok := svrc.methods[req.MethodName]
	// 请求的方法不存在，或者不是 RPC 方法（没有以 _rpc 结尾）
	if !ok {
		err = fmt.Errorf("service(%q) method %q is not a rpc method!", req.ServiceName, req.MethodName)
		codec.ReadRequestArg(req, nil)
		return
	}
	argv, err = methodInfo._getArgValue(codec, req)
	return
}

// 根据请求调用服务
func (svrc *s_Service) call(argv reflect.Value, req *S_ReqHeader) (reply interface{}, err error) {
	// 获取参数时，已经验证过一次，所以这里一定存在，不需要判断第二个返回值
	methodInfo, _ := svrc.methods[req.MethodName]

	// 构建返回值参数
	replyv := methodInfo._newReply()

	// 调用服务器函数
	methodInfo.Lock()
	methodInfo.numCalls++
	methodInfo.Unlock()
	fun := methodInfo.method.Func
	rets := fun.Call([]reflect.Value{svrc.rcvr, argv, replyv})
	reply = replyv.Interface()
	errInter := rets[0].Interface()
	if errInter != nil {
		err = errInter.(error)
	}
	return
}

// -----------------------------------------------------------------------------
// inner functions
// -----------------------------------------------------------------------------
// 获取注册服务的所有远程可调用方法
func _takeMethods(rcvrType reflect.Type) (methods map[string]*s_MethodInfo) {
	for n := 0; n < rcvrType.NumMethod(); n++ {
		method := rcvrType.Method(n)
		mtype := method.Type
		mname := method.Name

		// 排除私有方法
		if method.PkgPath != "" {
			continue
		}

		// 方法名称必须以 “_rpc” 结尾
		if !strings.HasSuffix(mname, "_rpc") {
			fslog.Errorf("fsrpc: rpc method's name(%q) is not ends with '_rpc'!\n", mname)
			continue
		}

		// 必须有且只有三个参数（第一个参数为处理器对象的指针）
		if mtype.NumIn() != 3 {
			fslog.Errorf("fsrpc: rpc method %q must be only contain 3 arguments!\n", mname)
			continue
		}

		// 第二个参数必须是一个可访问或内建类型参数
		argType := mtype.In(1)
		if !fsky.IsExposedOrBuiltinType(argType) {
			fslog.Errorf("fsrpc: argument type of method %q is not exposed: %q\n", mname, argType)
			continue
		}

		// 第三个参数必须是一个指针
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			fslog.Errorf("fsrpc: argument type of method %q is not a pointer but %q\n", mname, replyType)
			continue
		}

		// 返回值必须只有一个
		if mtype.NumOut() != 1 {
			fslog.Errorf("fsrpc: method %q has %d output paramaters, needs exactly 1.", mname, mtype.NumOut())
			continue
		}

		// 返回值必须是 error 类型
		if retType := mtype.Out(0); retType != fserror.RTypeStdError {
			fslog.Errorf("fsrpc: return type of method %p must be error, but not %q\n", mname, retType)
			continue
		}

		methods = make(map[string]*s_MethodInfo)
		methods[mname] = &s_MethodInfo{
			method:    method,
			argType:   argType,
			replyType: replyType,
			numCalls:  0,
		}
	}
	return
}

// -----------------------------------------------------------------------------
// package inner methods
// -----------------------------------------------------------------------------
// 新建一个服务
// 如果 sname 为空串，则以 rcvr 的结构名作为服务名称
// 客户端通过 “服务名称.方法名称” 来调用服务费方法
func newService(rcvr interface{}, sname string) (svr *s_Service, err error) {
	// rcvr 必须是一个结构示例指针
	if reflect.TypeOf(rcvr).Kind() != reflect.Ptr {
		err = fmt.Errorf("service(%v) must be an instance of struct pointer", rcvr)
		return
	}

	svr = new(s_Service)
	svr.rcvr = reflect.ValueOf(rcvr)
	svr.rcvrType = reflect.TypeOf(rcvr)
	svrName := reflect.Indirect(svr.rcvr).Type().Name()

	// 处理器类型必须是 public
	if !fsky.IsExposed(svrName) {
		err = fmt.Errorf("service %q must be public!", svrName)
		return
	}

	// 如果没有显式指定服务名称，则使用处理器结构名字
	if sname == "" {
		sname = svrName
	}
	svr.name = sname

	// 提取所有可远程调用方法
	svr.methods = _takeMethods(svr.rcvrType)
	return
}
