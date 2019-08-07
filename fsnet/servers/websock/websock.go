/**
@copyright: fantasysky 2016
@brief: 实现 websocket 服务器框架
@author: fanky
@version: 1.0
@date: 2019-03-19
**/

package websock

import (
	"fmt"
	"net/http"
	"time"
)

import (
	"github.com/gorilla/websocket"
)
import (
	"fsky.pro/fslog"
	"fsky.pro/fsnet/servers/base"
)

type S_WSServer struct {
	MaxConns     int           // 最大允许登录人数，默认 0，表示无上限
	ExitWaitTime time.Duration // 等待服务结束时间

	host     string
	port     uint16
	upgrader *websocket.Upgrader
	tlsCert  string // 证书 证书
	tlsKey   string // SSL 密匙

	clients map[base.T_ConnID]*S_WSConnInfo // 客户端列表

	registerCh   chan *S_WSConnInfo            // 客户端注册信息通道
	unregisterCh chan *S_WSConnInfo            // 客户端离线通道
	sendCh       chan map[base.T_ConnID][]byte // 发送消息通道
	broadcastCh  chan []byte                   // 待发送广播消息内容通道
	kickoutCh    chan map[base.T_ConnID][]byte // 关闭服务消息通道

	onlineCb  base.F_OnlineCb  // 上线通知消息内容
	offlineCb base.F_OfflineCb // 下线通知消息内容
	receiveCb base.F_ReceiveCb // 接收消息的回调

	isClosing bool // 正在关闭
	isClosed  bool // 标记服务器是否已经关闭
}

// NewWSServer 新建一个 WebSocket 服务
func NewServer(host string, port uint16) *S_WSServer {
	return &S_WSServer{
		MaxConns:     0,
		ExitWaitTime: time.Second * 2,

		host:     host,
		port:     port,
		upgrader: &websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }},
		tlsCert:  "",
		tlsKey:   "",

		clients: make(map[base.T_ConnID]*S_WSConnInfo),

		registerCh:   make(chan *S_WSConnInfo),
		unregisterCh: make(chan *S_WSConnInfo),
		sendCh:       make(chan map[base.T_ConnID][]byte),
		broadcastCh:  make(chan []byte),
		kickoutCh:    make(chan map[base.T_ConnID][]byte),

		onlineCb:  nil,
		offlineCb: nil,
		receiveCb: nil,

		isClosing: false,
		isClosed:  false,
	}
}

func (this *S_WSServer) String() string {
	return fmt.Sprintf("WebSocket Server(%s:%d)", this.host, this.port)
}

// NewWSServerSSL 新建一个 SSL websocket 服务
func NewServerTLS(host string, port uint16, tlsCert, tlsKey string) *S_WSServer {
	wss := NewServer(host, port)
	wss.tlsCert = tlsCert
	wss.tlsKey = tlsKey
	return wss
}

// -------------------------------------------------------------------
// inners
// -------------------------------------------------------------------
// 客户端登录或离线检测
func (this *S_WSServer) _run() {
	for {
		if this.isClosed {
			break
		}

		select {
		// 登入
		case client := <-this.registerCh:
			this.clients[client.ConnID] = client
			fslog.Infof("S_WSServer: a client(%v) has connected.", client)
			if this.onlineCb != nil {
				go this.onlineCb(&client.S_ConnInfo)
			}

		// 登出
		case client := <-this.unregisterCh:
			fslog.Infof("S_WSServer: a client(%v) will be deleted form websocket server!", client)
			_, ok := this.clients[client.ConnID]
			if ok {
				delete(this.clients, client.ConnID)
				fslog.Debugf("S_WSServer: a client(%v) has been deleted from wssocket server!", client)
			} else {
				fslog.Errorf("S_WSServer: a client(%v) has lost, but it is not in websocket list!", client)
			}
			if this.offlineCb != nil {
				go this.offlineCb(&client.S_ConnInfo)
			}

		// 发送消息
		case info := <-this.sendCh:
			for connid, msg := range info {
				client, ok := this.clients[connid]
				if ok {
					fslog.Debugf("S_WSServer: begin send message to client(%v), msg: %s", client, msg)
					go client.send(msg)
				} else {
					fslog.Errorf("S_WSServer: send message to client fail, connid(%d) is not exist in websocket list", connid)
				}
				break
			}

		// 广播消息
		case msg := <-this.broadcastCh:
			fslog.Info("S_WSServer: broadcast message:", msg)
			for _, client := range this.clients {
				go client.send(msg)
			}

		// 踢下客户端
		case info := <-this.kickoutCh:
			for connid, msg := range info {
				client, ok := this.clients[connid]
				if ok {
					fslog.Debugf("S_WSServer: begin kickout client(%v)", client)
					go client.kickout(msg)
				} else {
					fslog.Errorf("S_WSServer: kickout client fail, connid(%d) is not exist in websocket list", connid)
				}
				break
			}
		}
	}

	fslog.Infof(`S_WSServer: webscocket server(addr="%s:%p")'s mission loop has ended!`, this.host, this.port)
}

// websock 服务处理函数
func (this *S_WSServer) _handleRequest(rw http.ResponseWriter, req *http.Request) {
	if this.MaxConns > 0 && len(this.clients) >= this.MaxConns {
		fslog.Warnf("S_WSServer: one client request connect, but it is out of max connections(%d)", this.MaxConns)
		rw.Write([]byte("overload"))
		return
	}

	// 严格来说，这个地方需要锁，但是这里允许小几率错误，不影响整体环境
	if this.isClosing || this.isClosed {
		fslog.Infof(`S_WSServer: server(%v) has been closed, it can't accept connection already!`, this)
		return
	}

	conn, err := this.upgrader.Upgrade(rw, req, nil)
	if err != nil {
		fslog.Error("S_WSServer: create websocket upgrade fail: ", err.Error())
		http.NotFound(rw, req)
	} else {
		fslog.Infof("S_WSServer: a client(%v) request for connection.", conn.RemoteAddr())
		this.registerCh <- newConnInfo(this, conn)
	}
}

// -------------------------------------------------------------------
// package inners
// -------------------------------------------------------------------
// 接收到客户端消息时调用
func (this *S_WSServer) onReceiveMessage(client *S_WSConnInfo, msg []byte) {
	if this.receiveCb != nil {
		go this.receiveCb(&client.S_ConnInfo, msg)
	}
}

// 有连接断开时调用
func (this *S_WSServer) onConnectLost(client *S_WSConnInfo) {
	this.unregisterCh <- client
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// Serve 启动服务
func (this *S_WSServer) Serve(arg interface{}) error {
	addr := fmt.Sprintf("%s:%d", this.host, this.port)
	go this._run()

	mux := http.NewServeMux()
	mux.HandleFunc(arg.(string), this._handleRequest)
	var err error
	if this.tlsCert != "" {
		err = http.ListenAndServeTLS(addr, this.tlsCert, this.tlsKey, mux)
	} else {
		err = http.ListenAndServe(addr, mux)
	}
	this.Close(nil)
	return err
}

// Close 关闭服务
func (this *S_WSServer) Close(msg []byte) {
	if this.isClosed || this.isClosing {
		fslog.Warnf(`S_WSServer: server(%v) has been closed`, this)
		return
	}

	fslog.Infof(`S_WSServer: begin close server(%v)`, this)
	this.isClosing = true
	for _, client := range this.clients {
		go client.kickout(msg)
	}

	// 等待一段时间，期望所有客户端将离线数据处理完毕
	exitTime := time.Now().Add(this.ExitWaitTime)
	for len(this.clients) > 0 {
		// 如果所有客户端还没处理完毕，则直接退出
		if exitTime.Before(time.Now()) {
			fslog.Warnf(`S_WSServer: server(%v") begin close, and not wait all clients are closed`, this)
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	this.clients = make(map[base.T_ConnID]*S_WSConnInfo)
	this.isClosed = true
	fslog.Infof(`S_WSServer: server(%v") has been closed`, this)
}

// ---------------------------------------------------------
// 获取在线数量
func (this *S_WSServer) GetOnlineCount() int {
	return len(this.clients)
}

// 遍历所有连接信息
// 如果 handler 返回 false，则停止遍历
func (this *S_WSServer) IterConnInfos(handler base.F_IterConnsHandler) {
	for _, client := range this.clients {
		if !handler(&client.S_ConnInfo) {
			break
		}
	}
}

// ---------------------------------------------------------
// BindOnlineCb 设置客户端上线回调
func (this *S_WSServer) BindOnlineCb(cb base.F_OnlineCb) {
	this.onlineCb = cb
}

// BindOfflineCb 设置客户端离线通知回调
func (this *S_WSServer) BindOfflineCb(cb base.F_OfflineCb) {
	this.offlineCb = cb
}

// SetCBReceive 设置接收消息回调
func (this *S_WSServer) BindReceiveCb(cb base.F_ReceiveCb) {
	this.receiveCb = cb
}

// ---------------------------------------------------------
// Send 发送消息
func (this *S_WSServer) Send(connid base.T_ConnID, msg []byte) {
	if this.isClosed || this.isClosing {
		fslog.Errorf(`"S_WSServer: send message fail, server(%v) has been closed!`, this)
		return
	}
	if msg == nil {
		fslog.Error(`S_WSServer: send nil message via server(%v) is not to be allow`, this)
	} else {
		// 在循环中发送消息，目的是使得操控 this.clients 保证在同一个协程里
		this.sendCh <- map[base.T_ConnID][]byte{connid: msg}
	}
}

// Broadcast 广播消息
func (this *S_WSServer) Broadcast(msg []byte) {
	if this.isClosed || this.isClosing {
		fslog.Errorf(`S_WSServer: broadcast message fail, server(%v) has been closed!`, this)
		return
	}
	// 在循环中发送消息，目的是使得操控 this.clients 保证在同一个协程里
	this.broadcastCh <- msg
}

// Kickout 踢下线
//  如果 msg 为 nil，则不发送剔出信息到客户端
func (this *S_WSServer) Kickout(connid base.T_ConnID, msg []byte) {
	if this.isClosed || this.isClosing {
		fslog.Errorf(`"S_WSServer: kickout client fail, server(%v) has been closed!`, this)
		return
	}
	// 在循环中发送消息，目的是使得操控 this.clients 保证在同一个协程里
	this.kickoutCh <- map[base.T_ConnID][]byte{connid: msg}
}
