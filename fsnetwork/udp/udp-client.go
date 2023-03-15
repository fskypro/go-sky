/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: udp client
@author: fanky
@version: 1.0
@date: 2021-09-05
**/

package udp

import (
	"fmt"
	"net"
)

// -----------------------------------------------------------------------------
// UDPClient
// -----------------------------------------------------------------------------
type S_UDPClient struct {
	udpInfo    *S_UDPInfo
	conn       *net.UDPConn
	OnReceived func(error, []byte)
}

func NewClient(udpInfo *S_UDPInfo) *S_UDPClient {
	return &S_UDPClient{
		udpInfo: udpInfo,
	}
}

func (this *S_UDPClient) Dial() error {
	conn, err := net.DialUDP("udp", nil, &this.udpInfo.UDPAddr)
	if err != nil {
		return err
	}
	this.conn = conn
	return nil
}

// 如果需要接收返回消息，则需要调用该方法，否则可以不调用
func (this *S_UDPClient) Serve() {
	buff := make([]byte, this.udpInfo.BuffSize)
	for {
		n, raddr, err := this.conn.ReadFromUDP(buff)
		if raddr == nil {
			break
		}
		if this.OnReceived != nil {
			this.OnReceived(err, buff[:n])
		}
	}
}

func (this *S_UDPClient) Send(data []byte) (int, error) {
	return this.conn.Write(data)
}

func (this *S_UDPClient) Close() error {
	if this.conn == nil {
		return fmt.Errorf("udp client has not call dial method")
	}
	return this.conn.Close()
}
