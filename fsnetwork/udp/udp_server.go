/**
@copyright: fantasysky 2016
@brief: implement udp server
@author: fanky
@version: 1.0
@date: 2019-12-21
**/

package udp

import (
	"errors"
	"fmt"
	"net"
)

type S_UDPServer struct {
	udpInfo    *S_UDPInfo
	conn       *net.UDPConn
	OnReceived F_Receiver
	OnClosed   F_Closer
}

// 新建 UDP 服务
func NewServer(udpInfo *S_UDPInfo) (svr *S_UDPServer, err error) {
	conn, err := net.ListenUDP("udp", &udpInfo.UDPAddr)
	if err != nil {
		err = errors.New(fmt.Sprintf("listen udp addr(%s) fail: %s", udpInfo.String(), err.Error()))
		return
	}

	svr = &S_UDPServer{
		udpInfo: udpInfo,
		conn:    conn,
	}
	return
}

// -----------------------------------------------------------------------------
// private
// -----------------------------------------------------------------------------
func (this *S_UDPServer) onReceived(err error, raddr *net.UDPAddr, data []byte) {
	if this.OnReceived == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	go this.OnReceived(err, raddr, data)
}

// -----------------------------------------------------------------------------
// package public
// -----------------------------------------------------------------------------
func (this *S_UDPServer) Serve() {
	defer this.Close()
	for {
		data := make([]byte, this.udpInfo.BuffSize)
		n, raddr, err := this.conn.ReadFromUDP(data)
		if raddr == nil { // conn 已经被关闭
			if this.OnClosed != nil {
				this.OnClosed()
			}
			break
		}
		this.onReceived(err, raddr, data[:n])
	}
}

// 发送消息
func (this *S_UDPServer) Send(udpAddr *net.UDPAddr, data []byte) (int, error) {
	if this.conn == nil {
		return 0, errors.New(fmt.Sprintf("send message to client(%s) fail, udp listener has been down", udpAddr.String()))
	}
	return this.conn.WriteToUDP(data, udpAddr)
}

func (this *S_UDPServer) Close() {
	if this.conn != nil {
		this.conn.Close()
	}
}
