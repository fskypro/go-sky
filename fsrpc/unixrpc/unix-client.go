/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: unix rpc client
@author: fanky
@version: 1.0
@date: 2025-05-28
**/

package unixrpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"sync/atomic"
	"time"
)

type F_ClientCodec func(io.ReadWriteCloser) rpc.ClientCodec

// -------------------------------------------------------------------
// RpcClient
// -------------------------------------------------------------------
type s_RpcClient struct {
	*rpc.Client
	addr  string
	codec F_ClientCodec
}

func newRpcClient(addr string, codec F_ClientCodec) *s_RpcClient {
	return &s_RpcClient{
		addr:  addr,
		codec: codec,
	}
}

func (this *s_RpcClient) dial(timeout time.Duration) error {
	conn, err := createDialier(this.addr, timeout)
	if err != nil {
		return fmt.Errorf("dial to %q fail, %v", this.addr, err)
	}
	this.Client = rpc.NewClientWithCodec(this.codec(conn))
	return nil
}

// -------------------------------------------------------------------
// UnixClient
// -------------------------------------------------------------------
type S_UnixClient struct {
	Addr      string
	codec     F_ClientCodec
	idles     chan *s_RpcClient // 空闲连接
	sems      chan struct{}     // 允许新建连接令牌(令牌队列中如果能放入数据，则说明可以创建新连接)
	idleCount atomic.Int32      // 空闲连接数量
}

// 如果 codec == nil，则使用 gob 编码方式
func NewUnixClient(size int, addr string, codec F_ClientCodec) (*S_UnixClient, error) {
	uclient := &S_UnixClient{
		Addr:  addr,
		codec: codec,
		idles: make(chan *s_RpcClient, size),
		sems:  make(chan struct{}, size),
	}
	uclient.idleCount.Add(int32(size))
	return uclient, nil
}

// 获取一个空闲连接
func (this *S_UnixClient) get(timeout time.Duration) (*s_RpcClient, error) {
	select {
	case c := <-this.idles:
		this.idleCount.Add(-1)
		return c, nil
	default:
	}

	// 没有空闲连接，则判断新建令牌是否允许
	select {
	case this.sems <- struct{}{}:
		c := newRpcClient(this.Addr, this.codec)
		if err := c.dial(timeout); err != nil {
			return nil, fmt.Errorf("dial to server %q fail, %v", this.Addr, err)
		}
		this.idleCount.Add(-1)
		return c, nil
	default:
	}

	// 如果新建令牌也没拿到，则等待直到超时
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		// 一直没等到空闲连接，超时返回
		return nil, errors.New("all connections are busy")
	case c := <-this.idles:
		this.idleCount.Add(-1)
		return c, nil
	}
}

// 归还连接
func (this *S_UnixClient) put(c *s_RpcClient, healthy bool) {
	if healthy {
		select {
		case this.idles <- c:
			// 将连接放回空闲连接
			this.idleCount.Add(1)
			return
		default:
		}
	}

	// 如果是非健康连接，或者连接池已满，则关闭连接，释放令牌
	c.Close()
	select {
	case <-this.sems:
	default:
	}
}

// 根据 call 返回错误判断连接是否健康
func (this *S_UnixClient) isHealthyErr(err error) bool {
	if err == nil {
		return true
	}
	// 远端业务 error（rpc.ServerError）或你自定义的业务 error
	// 留给调用方自己判断；连接本身可以复用
	var srvErr rpc.ServerError
	if errors.As(err, &srvErr) {
		return true
	}

	// 这里本来可以直接返回 false 即可
	// 之所以继续编写下面的错误判断，只要是罗列出可能出现的底层错误

	// 1. net/rpc 自己的“我已经崩了”标记
	if errors.Is(err, rpc.ErrShutdown) {
		return false
	}

	// 2. 远端主动关连接 / 读写半途被掐
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return false
	}

	// 3. deadline、reset、no route、connection refused…
	var nerr net.Error
	if errors.As(err, &nerr) {
		// 任何 timeout 都判死；Temporary()==false 也判死
		if nerr.Timeout() || !nerr.Temporary() {
			return false
		}
	}

	// 4. 本地取消 / 上层 context 超时
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	// 5. 任何未知错误：宁可保守判“坏”，避免把隐患连带进池
	return false
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 空闲连接数量
func (this *S_UnixClient) IdleCount() int {
	return int(this.idleCount.Load())
}

func (this *S_UnixClient) Call(timeout time.Duration, method string, arg, ret any) error {
	c, err := this.get(timeout)
	if err != nil {
		return err
	}

	// 执行远程调用
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := c.Go(method, arg, ret, make(chan *rpc.Call, 1)).Done
	select {
	case call := <-done:
		// 等待调用返回
		this.put(c, this.isHealthyErr(call.Error))
		return call.Error
	case <-ctx.Done():
		// 超时退出
	}
	this.put(c, true)
	return ErrTimeout
}

func (this *S_UnixClient) Close() {
	close(this.idles)
	for c := range this.idles {
		_ = c.Close()
	}
	this.idleCount.Store(0)
	return
}
