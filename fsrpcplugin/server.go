/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: plugin server
@author: fanky
@version: 1.0
@date: 2023-05-03
**/

package fsrpcplugin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"

	"fsky.pro/fslog"
)

type S_Server struct {
}

func NewServer() *S_Server {
	return &S_Server{}
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func (this *S_Server) serve(ctx context.Context, conn net.Conn) {
	go rpc.ServeConn(conn)
	<-ctx.Done()
	conn.Close()
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_Server) Register(name string, rcvr any) error {
	return rpc.RegisterName(name, rcvr)
}

func (this *S_Server) Serve(ctx context.Context) error {
	ufile := os.Getenv("GO_PLUGIN_UNIX_SOCKET_FILE")
	if ufile == "" {
		return fmt.Errorf("unix socket file environment %q is not seted", "GO_PLUGIN_UNIX_SOCKET_FILE")
	}
	ln, err := net.Listen("unix", ufile)
	if err != nil {
		return fmt.Errorf("listen on unix file %q fail, %v", ufile, err)
	}

	go func() {
		<-ctx.Done()
		ln.Close()
		if err = os.Remove(ufile); err != nil {
			fslog.Errorf("remove unix socket file %q fail, %v", ufile, err)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			return fmt.Errorf("accept connect fail, %v", err)
		}
		go this.serve(ctx, conn)
	}
	return net.ErrClosed
}
