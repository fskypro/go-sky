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

import (
	"fsky.pro/fslog"
	"fsky.pro/fsstr/convert"
)

// -----------------------------------------------------------------------------
// module private
// -----------------------------------------------------------------------------
type S_UDPServer struct {
	udpInfo *S_UDPInfo
	connPtr *net.UDPConn
	recvCb  F_Receive
}

func (this *S_UDPServer) serve() {
	defer this.Close()
	for {
		if this.connPtr == nil {
			break
		}

		data := make([]byte, this.udpInfo.BuffSize)
		n, raddr, err := this.connPtr.ReadFromUDP(data)
		if raddr == nil { // conn 已经被关闭
			break
		}
		if err != nil {
			fslog.Errorf("read data from client(%s) fail, err: %s", raddr.String(), err.Error())
			continue
		}
		if this.recvCb == nil {
			fslog.Info("read data from client(%s): %s.", raddr.String(), convert.Bytes2String(data[:n]))
		} else {
			go this.recvCb(raddr, data[:n])
		}
	}
}

// -----------------------------------------------------------------------------
// package public
// -----------------------------------------------------------------------------
func (this *S_UDPServer) Run() {
	go this.serve()
}

// 发送消息
func (this *S_UDPServer) Send(udpAddr *net.UDPAddr, data []byte) (int, error) {
	if this.connPtr == nil {
		return 0, errors.New(fmt.Sprintf("send message to client(%s) fail, udp listener has been down", udpAddr.String()))
	}
	return this.connPtr.WriteToUDP(data, udpAddr)
}

func (this *S_UDPServer) Close() {
	if this.connPtr != nil {
		this.connPtr.Close()
		this.connPtr = nil
	}
}

// 绑定接收消息函数
func (this *S_UDPServer) BindReceiver(recv F_Receive) {
	this.recvCb = recv
}

// -----------------------------------------------------------------------------
// package public
// -----------------------------------------------------------------------------
// 新建 UDP 服务
func New(udpInfo *S_UDPInfo) (svr *S_UDPServer, err error) {
	conn, err := net.ListenUDP("udp", &udpInfo.UDPAddr)
	if err != nil {
		err = errors.New(fmt.Sprintf("listen udp addr(%s) fail: %s", udpInfo.String(), err.Error()))
		return
	}

	svr = &S_UDPServer{
		udpInfo: udpInfo,
		connPtr: conn,
	}
	return
}
