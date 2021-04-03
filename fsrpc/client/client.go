/**
* @copyright: fantasysky 2016
* @brief: 实现 RPC 客户端
* @author: fanky
* @version: 1.0
* @date: 2018-09-09
**/

package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"fsky.pro/fserror"
	"fsky.pro/fslog"
	"fsky.pro/fsrpc"
)

// -----------------------------------------------------------------------------
// client struct
// -----------------------------------------------------------------------------
// 客户端结构
type S_Client struct {
	sync.Mutex
	codec     I_ClientCodec     // 编码解码器
	reqHeader fsrpc.S_ReqHeader // 请求头，每次发送时临时使用，以免每次都创建又销毁

	mutex     sync.Mutex
	currReqID uint64                // 当前请求ID，给每个请求分配一个唯一ID
	pending   map[uint64]*S_ReqInfo // 已经发送出去的请求队列
	closing   bool
	shutdown  bool
}

// 客户端已关闭错误
var _errShutdown = errors.New("connection is not running")

// -------------------------------------------------------------------
// S_Client inner methods
// -------------------------------------------------------------------
// 发送请求
func (c *S_Client) _send(reqInfo *S_ReqInfo) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 链接已关闭
	if c.shutdown || c.closing {
		reqInfo.Error = _errShutdown
		reqInfo.send()
		return
	}
	reqID := c.currReqID
	c.currReqID++
	c.pending[reqID] = reqInfo

	// 编码数据并发送请求
	c.Lock()
	c.reqHeader.ServiceName = reqInfo.ServiceName
	c.reqHeader.MethodName = reqInfo.MethodName
	c.reqHeader.ReqID = reqID
	err := c.codec.WriteRequest(&c.reqHeader, reqInfo.Arg)
	c.Unlock()
	if err == nil {
		return
	} else {
		fslog.Error("fsrpc: send request fail: " + err.Error())
	}

	// 写入数据失败，则清空缓存数据
	reqInfo = c.pending[reqID]
	delete(c.pending, reqID)
	reqInfo.Error = err
	reqInfo.send()
}

// 接收请求
func (c *S_Client) _receive() {
	var err error
	var header fsrpc.S_RspHeader
	for err == nil {
		header = fsrpc.S_RspHeader{}
		err = c.codec.ReadResponseHeader(&header)
		if err != nil {
			break
		}

		reqID := header.ReqID
		c.mutex.Lock()
		defer c.mutex.Unlock()

		reqInfo := c.pending[reqID]
		delete(c.pending, reqID)

		switch {
		case reqInfo == nil:
			// 通常不会出现这种情况，除非无限低概率地读取到一个错误的头
			c.codec.ReadResponseReply(nil)
			fslog.Errorf("fsrpc: a request('%s.%s') has lost, ReqID=%d", header.ServiceName, header.MethodName, reqID)
		case header.Fail != "":
			// 调用失败
			reqInfo.Error = S_ServerError(header.Fail)
			c.codec.ReadResponseReply(nil)
			fslog.Errorf("fsrpc: server error, request('%s.%s') fail: %s", header.ServiceName, header.MethodName, header.Fail)
			reqInfo.send()
		case header.Error != "":
			// 调用返回错误
			reqInfo.Error = T_ServiceError(header.Error)
			c.codec.ReadResponseReply(nil)
			reqInfo.send()
		default:
			err := c.codec.ReadResponseReply(reqInfo.Reply)
			if err != nil {
				reqInfo.Error = errors.New("read reply error: " + err.Error())
				fslog.Errorf("fsrsp: read request(%s.%s)'s reply error: %s", header.ServiceName, header.MethodName, err.Error())
			}
			reqInfo.send()
		}
	}

	// 清理所有请求队列
	c.shutdown = true
	closing := c.closing
	if err == io.EOF {
		if closing {
			err = _errShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, reqInfo := range c.pending {
		reqInfo.Error = err
		reqInfo.send()
	}
	if err != io.EOF && !closing {
		fslog.Error("fsrpc: client protocol error: " + err.Error())
	}
}

// -------------------------------------------------------------------
// S_Client public methods
// -------------------------------------------------------------------
// TCP 拨号
func (c *S_Client) DialTCP(host string, port uint16) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial tcp(%s:%d) error: %v", host, port, err)
	} else {
		c.codec.Initialize(conn)
		go c._receive()
	}
	return nil
}

// HTTP 拨号
func (c *S_Client) DialHTTP(host string, port uint16) error {
	return c.DialHTTPPath(host, port, fsrpc.DefaultHTTPPath)
}

// HTTP 以指定路径拨号
func (c *S_Client) DialHTTPPath(host string, port uint16, path string) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial http server(%s:%d) error: %v", host, port, err)
	}

	// 发送链接校验
	io.WriteString(conn, "CONNECT "+path+" HTTP/1.0\n\n")

	// 获取校验返回
	rsp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err == nil {
		// 链接成功
		if rsp.Status == fsrpc.ConnectedText {
			// 初始化编码解码器
			c.codec.Initialize(conn)
			go c._receive()
			return nil
		} else {
			// 校验失败（原则上不会出现这种情况）
			err = fmt.Errorf("dial http server(%s:%d) error, verify fail", host, port)
		}
	}
	err = fmt.Errorf("dial http server(%s:%d) error, no respond from server", host, port)
	conn.Close()
	return err
}

// 异步调用服务器方法
// 注意：对于同一个 client，如果使用了 Go，就不能同时又使用 Call
func (c *S_Client) Go(svrc string, arg interface{}, reply interface{}, chReq chan *S_ReqInfo) *S_ReqInfo {
	svrcm := strings.Split(svrc, ".")
	if len(svrcm) != 2 {
		fslog.Panic("fsrpc: error service description, it must be: serviceName.methodName.")
	}

	reqInfo := new(S_ReqInfo)
	reqInfo.ServiceName = svrcm[0]
	reqInfo.MethodName = svrcm[1]
	reqInfo.Arg = arg
	reqInfo.Reply = reply
	if chReq == nil {
		chReq = make(chan *S_ReqInfo, 10)
	} else {
		// 要求通道必须带缓冲
		if cap(chReq) == 0 {
			fslog.Panic("fsrpc: S_ReqInfo channel is unbuffered!")
		}
	}
	reqInfo.ReqCh = chReq
	c._send(reqInfo)
	return reqInfo
}

// 同步调用服务器方法
// 注意：对于同一个 client，如果使用了 Call，就不能同时又使用 Go
func (c *S_Client) Call(svrc string, arg interface{}, reply interface{}) error {
	reqInfo := <-c.Go(svrc, arg, reply, make(chan *S_ReqInfo, 1)).ReqCh
	return reqInfo.Error
}

// 异步调用服务器方法
// cb 必须包含两个参数：
//    error 表示调用远程服务是否有错误
//    reply 表示返回值，必须是指针
//    extra 额外附带参数
func (c *S_Client) AyncCall(svrc string, arg interface{}, cb interface{}, data interface{}) error {
	cbt := reflect.TypeOf(cb)

	// cb 必须是函数
	if cbt.Kind() != reflect.Func {
		return errors.New("argument cb must a function!")
	}

	// 回调参数个数
	if cbt.NumIn() != 3 {
		return fmt.Errorf("cb function must contains 3 arguments(error, type of pointer, any tpye), but not %d", cbt.NumIn())
	}

	// 回调的第一个参数
	if cbt.In(0) != fserror.RTypeStdError {
		return fmt.Errorf("the first argument type of cb must be error, but not %q", cbt.In(0))
	}

	// 回调的第二个参数
	replyType := cbt.In(1)
	if replyType.Kind() != reflect.Ptr {
		return fmt.Errorf("the second argument type of cb must be a pointer, but not %q", replyType)
	}

	// 回调的第三个参数类型必须是 interface{}
	if cbt.In(2).Kind() != reflect.Interface {
		return fmt.Errorf("the third argument type of cb must be interface{}, but not %q", cbt.In(2))
	}

	// 创建返回值
	reply := reflect.New(replyType.Elem()).Interface()
	reqInfo := <-c.Go(svrc, arg, reply, nil).ReqCh

	errv := fserror.NilErrorValue
	if reqInfo.Error != nil {
		errv = reflect.ValueOf(reqInfo.Error)
	}
	datav := fserror.NilErrorValue
	if data != nil {
		datav = reflect.ValueOf(data)
	}
	reflect.ValueOf(cb).Call([]reflect.Value{errv, reflect.ValueOf(reply), datav})
	return nil
}

// 关闭链接
func (c *S_Client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closing {
		return _errShutdown
	}
	c.closing = true
	return c.codec.Close()
}

// -----------------------------------------------------------------------------
// package public
// -----------------------------------------------------------------------------
// 新建一个指定编码器的客户端
func NewClient(codec I_ClientCodec) *S_Client {
	return &S_Client{
		codec:   codec,
		pending: make(map[uint64]*S_ReqInfo),
	}
}
