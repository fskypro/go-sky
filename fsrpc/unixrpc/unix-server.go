/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: unix-rpc server
@author: fanky
@version: 1.0
@date: 2025-05-28
**/

package unixrpc

import (
	"context"
	"errors"
	"io"
	"net"
	"net/rpc"
	"os"
)

type F_ServerCodec func(io.ReadWriteCloser) rpc.ServerCodec

// -------------------------------------------------------------------
// S_UnixServer
// -------------------------------------------------------------------
type S_UnixServer struct {
	base  *rpc.Server
	Addr  string
	codec F_ServerCodec
	ln    net.Listener
}

// 如果 codec == nil，则使用 gob 编码
func NewUnixServer(addr string, codec F_ServerCodec) (*S_UnixServer, error) {
	ln, err := createListener(addr)
	if err != nil {
		return nil, err
	}
	return &S_UnixServer{
		base:  rpc.NewServer(),
		Addr:  addr,
		codec: codec,
		ln:    ln,
	}, nil
}

func (this *S_UnixServer) Run(ctx context.Context) {
	go func() {
		for {
			conn, err := this.ln.Accept()
			if errors.Is(err, net.ErrClosed) {
				return
			}
			if err != nil {
				continue
			}
			if this.codec == nil {
				go this.base.ServeConn(conn)
			} else {
				go this.base.ServeCodec(this.codec(conn))
			}
		}
	}()

	<-ctx.Done()
	this.ln.Close()
	os.Remove(this.Addr)
}

func (this *S_UnixServer) Close() error {
	return this.ln.Close()
}

func (this *S_UnixServer) Register(method string, recv any) {
	this.base.RegisterName(method, recv)
}
