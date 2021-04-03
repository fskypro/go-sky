/**
@copyright: fantasysky 2016
@brief: RPC 服务器
@author: fanky
@version: 1.0
@date: 2018-09-09
**/

package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"sync"

	"fsky.pro/fslog"
	"fsky.pro/fsrpc"
)

// ---------------------------------------------------------------------------------------
// S_Server
// ---------------------------------------------------------------------------------------
// 每个监听对应一个 server
// 每个 server 接受 n 个链接
// 每个链接对应一个 S_ServerCodec
// 每个服务接收器对应一个 service
type S_Server struct {
	codecer F_CodecCreator // 解码器创建函数

	services sync.Map // 服务对象列表

	reqLock sync.Mutex
	freeReq *S_ReqCache // 空闲请求列表

	rspLock sync.Mutex
	freeRsp *S_RspCache // 空闲回复列表
}

// 请求头缓存
type S_ReqCache struct {
	header fsrpc.S_ReqHeader
	next   *S_ReqCache
}

// 回复头缓存
type S_RspCache struct {
	header fsrpc.S_RspHeader
	next   *S_RspCache
}

// -----------------------------------------------------------------------------
// S_Server inner methods
// -----------------------------------------------------------------------------
// 获取一个请求头数据
// server.freeReq 是一个链表，保存空闲的请求头数据
// 每当一个请求过来，先获取 server.freeReq 第一个空闲数据，如果不存在则新建一个
// 等到请求处理完毕，通过 _freeReqHeader 方法回收到 server.freeReq 链表中，以供下一个请求使用
// 这样做的目的是防止频繁创建和销毁 S_ReqCache
func (s *S_Server) _getReqHeader() *S_ReqCache {
	s.reqLock.Lock()
	defer s.reqLock.Unlock()
	req := s.freeReq
	if req == nil {
		req = new(S_ReqCache)
	} else {
		s.freeReq = req.next
		*req = S_ReqCache{}
	}
	return req
}

// 释放并回收请求数据
func (s *S_Server) _freeReqHeader(req *S_ReqCache) {
	s.reqLock.Lock()
	defer s.reqLock.Unlock()
	req.next = s.freeReq
	s.freeReq = req
}

// 获取一个回复信息数据
func (s *S_Server) _getRspHeader() *S_RspCache {
	s.rspLock.Lock()
	defer s.rspLock.Unlock()
	rsp := s.freeRsp
	if rsp == nil {
		rsp = new(S_RspCache)
	} else {
		s.freeRsp = rsp.next
		*rsp = S_RspCache{}
	}
	return rsp
}

// 释放并回收回复数据
func (s *S_Server) _freeRspHeader(rsp *S_RspCache) {
	s.rspLock.Lock()
	defer s.rspLock.Unlock()
	rsp.next = s.freeRsp
	s.freeRsp = rsp
}

// -------------------------------------------------------------------
// 获取请求头信息
// 获取请求对应的 service 和 method
func (s *S_Server) _readReqHeader(codec S_ServerCodec) (req *S_ReqCache, ok bool) {
	req = s._getReqHeader()
	// 读取并解码客户端请求头
	err := codec.ReadRequestHeader(&req.header)
	// 链接已关闭或者链接断开
	if err == io.EOF {
		return
	}
	if err == io.ErrUnexpectedEOF {
		fslog.Error("fsrpc: read request header fail: " + err.Error())
		return
	}

	// 解码错误，通常不会出现这种错误
	// 除非 fsrpc.S_ReqHeader 结构改了，而服务器和客户端版本不一样，不同时改
	if err != nil {
		err = errors.New("fsrpc: server can't decode request: " + err.Error())
	}
	ok = true
	return
}

// 回复客户端请求
// fail 表示远程调用失败
// err 表示远程调用返回的错误值
func (s *S_Server) _sendResponse(sendMutex *sync.Mutex, codec S_ServerCodec, reqHeader *fsrpc.S_ReqHeader, reply interface{}, fail error, err error) {
	rsp := s._getRspHeader()
	rspHeader := &rsp.header
	rspHeader.ReqID = reqHeader.ReqID
	if fail != nil {
		rspHeader.Fail = fail.Error()
	}
	if err != nil {
		rspHeader.Error = err.Error()
	}
	sendMutex.Lock()
	err = codec.WriteResponse(rspHeader, reply)
	sendMutex.Unlock()
	if err != nil {
		fslog.Error("fsrpc: writing response fail: " + err.Error())
	}
	s._freeRspHeader(rsp)
}

// 处理服务请求
func (s *S_Server) _handleService(wg *sync.WaitGroup, sendMutex *sync.Mutex, codec S_ServerCodec, req *S_ReqCache) {
	var argv reflect.Value // 调用服务的传入参数
	var err error = nil    // 调用服务的返回值
	header := &req.header

	// 根据服务名称，加载服务对象
	svrc, ok := s.services.Load(header.ServiceName)
	if !ok {
		err = fmt.Errorf("request service %q is not exists!", header.ServiceName)
		codec.ReadRequestArg(header, nil)
	} else {
		// 获取调用服务的传入参数
		// 注意，将获取参数和调用服务分开两个方法的目的是：
		//     获取方法参数需要访问 codec，所以不能放到另一个协程上，否则要加锁
		//     而调用服务处理函数，则需要放到另一个协程上，以提高效率
		//     因此，这里需要先调用个方法获取到服务的传入参数，然后下面再另起一个协程，以执行服务方法
		argv, err = svrc.(*s_Service).getArgValue(codec, header)
	}
	if err != nil {
		fslog.Error("fsrpc: " + err.Error())
	}

	go func() {
		wg.Done()
		// 只有没有任何错误时，才调用服务
		if err == nil {
			reply, err := svrc.(*s_Service).call(argv, header)
			s._sendResponse(sendMutex, codec, header, reply, nil, err)
		} else {
			// 任何错误都需要回复客户端
			s._sendResponse(sendMutex, codec, header, nil, err, nil)
		}
		s._freeReqHeader(req)
	}()
}

// 对单个链接服务
// 每接受一个链接则启动一个 goroutine
func (s *S_Server) _serveConn(conn io.ReadWriteCloser) {
	codec := s.codecer(conn)
	sendMutex := new(sync.Mutex) // Response 缓冲锁
	wg := new(sync.WaitGroup)
	for {
		// 读取并解码请求头信息
		// 注意：读缓冲不需要锁，因为读取是在同一个 goroutine 中
		req, ok := s._readReqHeader(codec)
		if !ok {
			s._freeReqHeader(req)
			// 如果读取头失败，则意味着服务器和客户端的 fsrpc 版本不一致，链接将会断开
			break
		} else {
			wg.Add(1)
			// 调用请求方法，并将调用结果写入客户端 Response 缓冲
			s._handleService(wg, sendMutex, codec, req)
		}
	}
	wg.Wait()
	codec.Close()
}

// -----------------------------------------------------------------------------
// S_Server public methods
// -----------------------------------------------------------------------------
// 绑定一个处理器（以处理器类型名称作为访问名称）
func (s *S_Server) Register(rcvr interface{}) error {
	return s.RegisterByName(rcvr, "")
}

// 绑定一个处理器（并指定访问名称）
func (s *S_Server) RegisterByName(rcvr interface{}, name string) (err error) {
	service, err := newService(rcvr, name)
	if err != nil {
		fslog.Error("fsrpc: " + err.Error())
	}
	_, loaded := s.services.LoadOrStore(service.name, service)
	if loaded {
		err = fmt.Errorf("fsrpc: %q has been registered!", service.name)
		fslog.Errorf(err.Error())
	}
	return
}

// ------------------------------------------------------------------
// HTTP 服务启动参数
type S_HttpServeArg struct {
	RpcPath   string // HTTP 访问路径（只有 HTTP 协议有）
	DebugPath string // HTTP debug 访问路径（只有 HTTP 协议有）

	// TLS 的 pem 和 key 文件
	certFile string
	keyFile  string
}

// 启动 TCP 协议服务（host 格式为：地址:端口）
func (s *S_Server) ServeTCP(host string, port uint16) {
	addr := fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fslog.Error("fsrpc: can't start tcp serve for " + host)
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			fslog.Errorf("fsrpc: tcp accept connection error: %s\n", err.Error())
		}
		go s._serveConn(conn)
	}
}

// 启动 HTTP 协议服务
func (s *S_Server) ServeHTTP(host string, port uint16, arg *S_HttpServeArg) {
	rpcPath := fsrpc.DefaultHTTPPath
	if arg != nil {
		rpcPath = arg.RpcPath
	}

	// handle rpc
	http.HandleFunc(rpcPath, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "CONNECT" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "fsrpc: request error 405: must CONNECT\n")
			return
		}
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			fslog.Error("fsrpc: rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		} else {
			io.WriteString(conn, "HTTP/1.0 "+fsrpc.ConnectedText+"\n\n")
			s._serveConn(conn)
		}
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	if arg != nil && arg.certFile != "" && arg.keyFile != "" {
		http.ListenAndServeTLS(addr, arg.certFile, arg.keyFile, nil)
	} else {
		http.ListenAndServe(addr, nil)
	}
}

// ---------------------------------------------------------------------------------------
// package public methods
// ---------------------------------------------------------------------------------------
// 新建一个服务
func NewServer(codecer F_CodecCreator) *S_Server {
	return &S_Server{codecer: codecer}
}
