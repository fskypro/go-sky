/**
@copyright: fantasysky 2016
@brief: 实现 websocket 连接信息
@author: fanky
@version: 1.0
@date: 2019-03-19
**/

package websock

import (
	"context"
	"net"
	"strconv"
	"sync"
)
import (
	"github.com/gorilla/websocket"
)
import (
	"fsky.pro/fsenv"
	"fsky.pro/fslog"
	"fsky.pro/fsnet/servers/base"
)

type S_WSConnInfo struct {
	base.S_ConnInfo
	ownerPtr *S_WSServer     // 所属的 WebSocket server
	connPtr  *websocket.Conn // websock 底层连接

	sendLocker  sync.Mutex         // 发送通道锁
	sendCh      chan []byte        // 发送消息通道
	ctxCanceler context.CancelFunc // 发送、接收协程关闭控制器
}

func newConnInfo(owner *S_WSServer, conn *websocket.Conn) *S_WSConnInfo {
	ip, sport, _ := net.SplitHostPort(conn.RemoteAddr().String())
	port, _ := strconv.ParseUint(sport, 10, 16)
	connInfo := &S_WSConnInfo{
		S_ConnInfo: *base.NewConnInfo(ip, uint16(port)),
		ownerPtr:   owner,
		connPtr:    conn,
		sendCh:     make(chan []byte),
	}

	// 发送通道上下文控制器
	ctx, canceler := context.WithCancel(context.Background())
	connInfo.ctxCanceler = canceler
	// 协调发送、接收通道一致性等待组
	wg := sync.WaitGroup{}

	go connInfo._cycRecive(&wg)
	go connInfo._cycSend(&wg, ctx)
	go connInfo._testOffline(&wg)
	return connInfo
}

// 轮询读取客户端内容
func (this *S_WSConnInfo) _cycRecive(wg *sync.WaitGroup) {
	state := base.CONN_STATE_LOST
	// 客户端自己离线时，在这里停止发送协程
	defer func() {
		this._exit(state)
		wg.Done()
		fslog.Debugf("S_WSConnInfo: receive loop of client(%v) has been end", this)
	}()

	wg.Add(1)
	for {
		// 接收客户端消息
		_, msg, err := this.connPtr.ReadMessage()
		if err == nil {
			fslog.Debugf("S_WSConnInfo: receive message from client(%v):%s\t%s", this, fsenv.Endline, msg)
			this.ownerPtr.onReceiveMessage(this, msg)
		} else {
			switch err.(type) {
			case *net.OpError: // 服务器主动关闭时
				state = base.CONN_STATE_KICKOUT
				fslog.Infof("S_WSServer: connection to client(%v) has been closed, and receive groutine weill stop info: %s", this, err.Error())
			case *websocket.CloseError: // 客户端离线时
				fslog.Infof("S_WSConnInfo: a client(%v) has been closed, reveive goroutine will stop info: %s", this, err.Error())
			default:
				fslog.Errorf("S_WSConnInfo: receive message from client(%v) fail: %s", this, err.Error())
			}
			break
		}
	}
}

// 轮询发送消息到客户端
func (this *S_WSConnInfo) _cycSend(wg *sync.WaitGroup, ctx context.Context) {
	// 在这里清理连接
	defer func() {
		this.connPtr.Close() // 服务器踢下线时，在这里主动关闭连接，迫使退出接收信息循环
		wg.Done()
		fslog.Debugf("S_WSConnInfo: send loop of client(%v) has been end", this)
	}()

	wg.Add(1)
L:
	for {
		select {
		case <-ctx.Done():
			fslog.Debugf("S_WSConnInfo: send channel of client(%v) has been closed!", this)
			break L
		case msg := <-this.sendCh:
			fslog.Debugf("S_WSConnInfo: write message to web socket!")
			err := this.connPtr.WriteMessage(websocket.TextMessage, msg)
			if err != nil { // 离线
				fslog.Errorf("S_WSConnInfo: send message to client(%v) fail! err:", this, err.Error())
			}
		}
	}
}

// 离线检测
func (this *S_WSConnInfo) _testOffline(wg *sync.WaitGroup) {
	this.SetState(base.CONN_STATE_ONLINE)
	wg.Wait()
	// 此时，接收和发送循环都已经结束
	// 清空发送通道中的数据
	select {
	case msg := <-this.sendCh:
		fslog.Warnf("S_WSConnInfo: client(%v) will be offline, but a message in send channel will be abandoned, msg=%q", this, msg)
	default:
		break
	}
	fslog.Infof("S_WSConnInfo: client(%v)'s goroutine of receive/send has been end.", this)
	this.ownerPtr.onConnectLost(this)
}

// 退出连接
func (this *S_WSConnInfo) _exit(state base.T_ConnState) {
	if !this.IsState(base.CONN_STATE_ONLINE) {
		fslog.Debugf("S_WSConnInfo: client(%v) has been offline, it is not need to exit", this)
		return
	}
	this.SetState(state)
	this.ctxCanceler() // 友好关闭发送通道
}

// -------------------------------------------------------------------
// 发送消息
func (this *S_WSConnInfo) send(msg []byte) {
	// 用锁的目的是保证 send 通道上只能有一个等待数据
	// 当有一个发送数据在等待时，再有其他更多发送数据到来的话，将会卡在锁上，而不是通道上
	// 这样将会造成死锁
	// 解决办法是：
	//		一旦发送循环结束，则堵在发送通道上的数据最多只有一个（这个先别管他），其他发送队列都堵在锁上。
	//		当发送和接收循环都结束后，会做两件事情：
	//			1、将连接状态设置为离线（即 this.IsState(base.CONN_STATE_ONLINE) == false）
	//			2、会释放掉发送通道上的数据，从而解锁发送通道，发送通道解锁后，sendLocker 也会跟着解锁
	//		此时，之前被卡在锁上的发送数据，将会得到继续。
	//		但是，继续后，将无法通过状态判断（即：this.IsState(base.CONN_STATE_ONLINE) == false）
	//		为此，数据也不会再往 send 通道里发送。
	//		这样就解决了发送通道死锁的问题
	this.sendLocker.Lock()
	defer this.sendLocker.Unlock()
	if this.IsState(base.CONN_STATE_ONLINE) {
		this.sendCh <- msg
	} else {
		fslog.Warnf("S_WSConnInfo: client(%v) has been offline, send message fail! msg=%s", this, msg)
	}
}

// 关闭客户端
func (this *S_WSConnInfo) kickout(msg []byte) {
	if msg != nil {
		err := this.connPtr.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fslog.Errorf("S_WSConnInfo: begin kickout client(%v), but send kickout message to client fail, err: %s", this, err.Error())
		}
	}
	fslog.Debugf("S_WSConnInfo: begin kickout client(%v).", this)
	this._exit(base.CONN_STATE_KICKOUT)
}
