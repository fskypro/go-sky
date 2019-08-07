/**
* @copyright: fantasysky 2016
* @brief: 实现 RPC 客户端
* @author: fanky
* @version: 1.0
* @date: 2018-09-09
**/

package client

import "io"
import "bufio"
import "fmt"
import "sync"
import "net"
import "time"
import "strings"
import "errors"
import "reflect"
import "net/http"
import "fsky.pro/fslog"
import "fsky.pro/fserror"
import "fsky.pro/fsrpc"

// -----------------------------------------------------------------------------
// Error type
// -----------------------------------------------------------------------------
// 远程请求失败错误
type S_ServerError string

func (e S_ServerError) Error() string {
	return string(e)
}

// 调用服务函数返回错误
type T_ServiceError string

func (e T_ServiceError) Error() string {
	return string(e)
}

// -----------------------------------------------------------------------------
// request info
// -----------------------------------------------------------------------------
// 请求信息
type S_ReqInfo struct {
	ServiceName string          // 请求的服务对象名称
	MethodName  string          // 请求的方法
	Arg         interface{}     // 请求参数
	Reply       interface{}     // 请求回复参数
	Error       error           // 请求失败或调用服务器服务函数返回错误
	Done        chan *S_ReqInfo //请求返回后数据放入的通道
}

// 结束远程调用，将结构写入用户通道
func (req *S_ReqInfo) _done() {
	select {
	case req.Done <- req:
	default:
		// 如果 req.Done 阻塞（已满），则程序会跑这里来，
		// 这时，等待 100 毫秒再尝试往通道发送
		go func() {
			time.Sleep(time.Millisecond * 100)
			select {
			case req.Done <- req:
			default:
				// 如果 100 毫秒后，还是无法发送，则丢弃该请求
				// 这样做的目的是防止程序被无限期阻塞，因此 req.Done 必须设置足够的缓冲数量
				fslog.Errorf("fsrpc: a reqiest(%s.%s) has been discarded!", req.ServiceName, req.MethodName)
			}
		}()
	}
}

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
var _errShutdown = errors.New("connection is shut down")

// -------------------------------------------------------------------
// S_Client inner methods
// -------------------------------------------------------------------
// 发送请求
func (c *S_Client) _send(reqInfo *S_ReqInfo) {
	c.mutex.Lock()
	// 链接已关闭
	if c.shutdown || c.closing {
		reqInfo.Error = _errShutdown
		c.mutex.Unlock()
		reqInfo._done()
		return
	}
	reqID := c.currReqID
	c.currReqID++
	c.pending[reqID] = reqInfo
	c.mutex.Unlock()

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
	c.mutex.Lock()
	reqInfo = c.pending[reqID]
	delete(c.pending, reqID)
	c.mutex.Unlock()
	reqInfo.Error = err
	reqInfo._done()
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
		reqInfo := c.pending[reqID]
		delete(c.pending, reqID)
		c.mutex.Unlock()

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
			reqInfo._done()
		case header.Error != "":
			// 调用返回错误
			reqInfo.Error = T_ServiceError(header.Error)
			c.codec.ReadResponseReply(nil)
			reqInfo._done()
		default:
			err := c.codec.ReadResponseReply(reqInfo.Reply)
			if err != nil {
				reqInfo.Error = errors.New("read reply error: " + err.Error())
				fslog.Errorf("fsrsp: read request(%s.%s)'s reply error: %s", header.ServiceName, header.MethodName, err.Error())
			}
			reqInfo._done()
		}
	}

	// 清理所有请求队列
	c.mutex.Lock()
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
		reqInfo._done()
	}
	c.mutex.Unlock()
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
		fslog.Error("fsrpc: " + err.Error())
		return err
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
		fslog.Errorf("fsrpc: http dial server(%s) fail: %s\n", addr, err.Error())
		return err
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
			err = fserror.StrErrorf("http server(%s) verify fail!\n", addr)
			fslog.Error("fsrpc: " + err.Error())
		}
	} else {
		fslog.Errorf("fsrpc: get http server(%s) verify response fail: %s\n", addr, err.Error())
	}
	conn.Close()
	return &net.OpError{
		Op:   "dial-http",
		Net:  "tcp " + addr,
		Addr: nil,
		Err:  err,
	}
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
	reqInfo.Done = chReq
	c._send(reqInfo)
	return reqInfo
}

// 同步调用服务器方法
// 注意：对于同一个 client，如果使用了 Call，就不能同时又使用 Go
func (c *S_Client) Call(svrc string, arg interface{}, reply interface{}) error {
	reqInfo := <-c.Go(svrc, arg, reply, make(chan *S_ReqInfo, 1)).Done
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
		return fserror.StrErrorf("cb function must contains 3 arguments(error, type of pointer, any tpye), but not %d", cbt.NumIn())
	}

	// 回调的第一个参数
	if cbt.In(0) != fserror.RTypeStdError {
		return fserror.StrErrorf("the first argument type of cb must be error, but not %q", cbt.In(0))
	}

	// 回调的第二个参数
	replyType := cbt.In(1)
	if replyType.Kind() != reflect.Ptr {
		return fserror.StrErrorf("the second argument type of cb must be a pointer, but not %q", replyType)
	}

	// 回调的第三个参数类型必须是 interface{}
	if cbt.In(2).Kind() != reflect.Interface {
		return fserror.StrErrorf("the third argument type of cb must be interface{}, but not %q", cbt.In(2))
	}

	// 创建返回值
	reply := reflect.New(replyType.Elem()).Interface()
	reqInfo := <-c.Go(svrc, arg, reply, nil).Done

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
	if c.closing {
		c.mutex.Unlock()
		return _errShutdown
	}
	c.closing = true
	c.mutex.Unlock()
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
