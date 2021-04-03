/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: tcp erver
@author: fanky
@version: 1.0
@date: 2021-03-20
**/

package tcprpc

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"fsky.pro/fsrpcs"
)

type S_Server struct {
	*rpc.Server
	listener  *net.TCPListener
	heartbeat *s_HeartbeatService
}

func NewServer() *S_Server {
	server := &S_Server{
		Server:    rpc.NewServer(),
		heartbeat: newHeartbeatService(),
	}
	server.Register(server.heartbeat)
	return server
}

func (this *S_Server) SetHeartbeatReply(reply []byte) {
	this.heartbeat.Reply = reply
}

func (this *S_Server) Register(service fsrpcs.I_Service) {
	// 因为 rpc 包对 service 成员方法做了限制：
	//   所有共有成员方法必须要符合被客户端调用的规格，否则会 log 输出警告
	// 因此，这里将 log 警告输出暂时忽略掉
	bs := bytes.NewBuffer([]byte{})
	log.SetOutput(bs)
	this.Server.RegisterName(service.Name(), service)
	log.SetOutput(os.Stderr)
}

// 如果 addr 的端口字段为0，函数将选择一个当前可用的端口，可以用 Addr 方法获得该端口。
func (this *S_Server) Listen(netStr string, addr *net.TCPAddr) (*net.TCPListener, error) {
	listener, err := net.ListenTCP(netStr, addr)
	if err != nil {
		return nil, fmt.Errorf("start tcp rpc server fail: %v", err)
	}
	this.listener = listener
	return listener, nil
}

// Close 停止监听 TCP 地址
func (this *S_Server) Close() {
	if this.listener != nil {
		this.listener.Close()
	}
}

// 使用 gob 编码
func (this *S_Server) Serve(ch chan *net.TCPConn) {
	for {
		conn, err := this.listener.AcceptTCP()
		if err == nil {
			if ch != nil {
				ch <- conn
			}
			go this.Server.ServeConn(conn)
		} else {
			break
		}
	}
}

// 使用自定义编码器
func (this *S_Server) ServeCodec(codec I_TcpServerCodec, ch chan *net.TCPConn) {
	for {
		conn, err := this.listener.AcceptTCP()
		if err == nil {
			if ch != nil {
				ch <- conn
			}
			codec.SetConn(conn)
			go this.Server.ServeCodec(codec)
		} else {
			break
		}
	}
}
