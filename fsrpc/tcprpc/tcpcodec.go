/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: tcp codec base
@author: fanky
@version: 1.0
@date: 2021-03-28
**/

package tcprpc

import (
	"net"
	"net/rpc"
)

// -------------------------------------------------------------------
// ServeCodec
// 如果自定义编码器，则需要对服务器编码器实现该接口，并对 SetConn 中拿到的 conn 进行流操作
// 可以参考：rpc.gobServerCodec
// -------------------------------------------------------------------
type I_TcpServerCodec interface {
	rpc.ServerCodec
	SetConn(*net.TCPConn)
}

// -------------------------------------------------------------------
// ClientCodec
// 如果自定义编码器，则需要对客户端编码器实现该接口，并对 SetConn 中拿到的 conn 进行流操作
// 可以参考：rpc.gobClientCodec
// -------------------------------------------------------------------
type I_TcpClientCodec interface {
	rpc.ClientCodec
	SetConn(*net.TCPConn)
}
