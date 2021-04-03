/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: service proxy for client
@author: fanky
@version: 1.0
@date: 2021-03-20
**/

package tcprpc

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/rpc"
	"time"

	"fsky.pro/fsrpcs"
)

// -------------------------------------------------------------------
// 连接状态（只有调用 ServerProy 的 Run 方法时，才起作用）
// -------------------------------------------------------------------
// 连接状态标记
type E_RunState int

const (
	TRY_DIAL  E_RunState = 1 // 准备连接服务器
	DIAL_FAIL            = 2 // 连接服务器失败
	DIAL_SUCC            = 3 // 连接服务器成功
	LOSE_CONN            = 4 // 与服务器失去连接
)

// 连接状态信息
type S_RunState struct {
	RemoteAddr *net.TCPAddr // 服务器地址
	Conn       *net.TCPConn //
	State      E_RunState   // 状态标记
	Error      error
}

func _newRunState(raddr *net.TCPAddr) *S_RunState {
	return &S_RunState{RemoteAddr: raddr}
}

func (this *S_RunState) reset(state E_RunState, conn *net.TCPConn, err error) *S_RunState {
	this.State = state
	this.Error = err
	return this
}

// -------------------------------------------------------------------
// ServerProxy
// -------------------------------------------------------------------
// base struct for client to call seerves method
type S_ServerProxy struct {
	fsrpcs.I_ServiceProxy
	*rpc.Client
	codec     I_TcpClientCodec
	closed    bool
	connected bool

	netstr    string
	raddr     *net.TCPAddr
	heartbeat *s_HeartbeatClient
}

func NewServerProxy(netstr string, addr *net.TCPAddr) *S_ServerProxy {
	return &S_ServerProxy{
		closed:    false,
		connected: false,
		netstr:    netstr,
		raddr:     addr,
		heartbeat: nil,
	}
}

func NewServerProxWithCodec(codec I_TcpClientCodec, netstr string, addr *net.TCPAddr) *S_ServerProxy {
	proxy := NewServerProxy(netstr, addr)
	proxy.codec = codec
	return proxy
}

// ---------------------------------------------------------
// 设置心跳参数，注意：
//   必须调用 Run() 方法拨号才能触发心跳
func (this *S_ServerProxy) SetHeartbeatInfo(info *S_HeartbeatInfo) {
	if this.heartbeat != nil {
		this.heartbeat.updateInfo(info)
	}
}

// ---------------------------------------------------------
func (this *S_ServerProxy) Call(method string, arg interface{}, reply interface{}) error {
	if !this.connected || this.Client == nil {
		return fmt.Errorf("call service method %q fail, client hasn't connect to server yet.", method)
	}
	return this.Client.Call(method, arg, reply)
}

func (this *S_ServerProxy) Go(method string, arg interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	if !this.connected || this.Client == nil {
		call := new(rpc.Call)
		call.ServiceMethod = method
		call.Args = arg
		call.Error = fmt.Errorf("call service method %q fail, client hasn't connect to server yet.", method)
		call.Done = make(chan *rpc.Call, 2)
		call.Done <- call
		return call
	}
	call := this.Client.Go(method, arg, reply, done)
	return call
}

// ---------------------------------------------------------
func (this *S_ServerProxy) Dial() (*net.TCPConn, error) {
	conn, err := net.DialTCP(this.netstr, nil, this.raddr)
	if err != nil {
		return nil, fmt.Errorf("dial server fail: %v", err)
	}
	if this.codec != nil {
		this.codec.SetConn(conn)
		this.Client = rpc.NewClientWithCodec(this.codec)
	} else {
		this.Client = rpc.NewClient(conn)
	}
	return conn, nil
}

// 会创建心跳，自动重连
// <-ch 为 1 表示连接成功；为 0 表示断开连接
func (this *S_ServerProxy) Run(ctx context.Context, ch chan *S_RunState) {
	this.closed = false
	defer this.Close()
	this.heartbeat = newHeartbeatClient(this)

	runState := _newRunState(this.raddr)
	const maxInterval = 10
	interval := 1
	for {
		select {
		case <-ctx.Done():
			this.Close()
			return
		default:
			if this.closed {
				return
			}
			if ch != nil {
				ch <- runState.reset(TRY_DIAL, nil, nil)
			}
			conn, err := this.Dial()
			if err != nil {
				if ch != nil {
					ch <- runState.reset(DIAL_FAIL, conn, err)
				}
				time.Sleep(time.Second * time.Duration(interval))
				interval = int(math.Min(float64(interval+1), float64(maxInterval)))
			} else {
				this.connected = true
				interval = 1
				if ch != nil {
					ch <- runState.reset(DIAL_SUCC, conn, nil)
				}
				this.heartbeat.loopSend(ctx, conn)
				this.connected = false
				if ch != nil {
					ch <- runState.reset(LOSE_CONN, conn, nil)
				}
			}
		}
	}
}

func (this *S_ServerProxy) Close() error {
	this.closed = true
	this.connected = false
	return this.Client.Close()
}
