/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: server
@author: fanky
@version: 1.0
@date: 2021-12-06
**/

package fshttp

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type S_Server struct {
	host string
	port int

	tlsKey   string
	tlsPem   string
	listener net.Listener
}

func NewServer(host string, port int) *S_Server {
	return &S_Server{
		host: host,
		port: port,
	}
}

func NewServerAddr(addr string) (*S_Server, error) {
	hp := strings.Split(addr, ":")
	if len(hp) != 2 {
		return nil, fmt.Errorf("error http address %q", addr)
	}
	port, _ := strconv.Atoi(hp[1])
	if port <= 0 {
		return nil, fmt.Errorf("error http address %q", addr)
	}
	return &S_Server{
		host: hp[0],
		port: port,
	}, nil
}

func (this *S_Server) Host() string {
	return this.host
}

func (this *S_Server) Port() int {
	return this.port
}

func (this *S_Server) Addr() string {
	return fmt.Sprintf("%s:%d", this.host, this.port)
}

func (this *S_Server) SetTLS(pem, key string) {
	this.tlsPem = pem
	this.tlsKey = key
}

func (this *S_Server) Listen() error {
	ln, err := net.Listen("tcp", this.Addr())
	this.listener = ln
	return err
}

func (this *S_Server) Serve(service I_Service) error {
	if this.listener == nil {
		return errors.New("must listen at first.")
	}
	return http.Serve(this.listener, service)
}

func (this *S_Server) ServeTLS(service I_Service) error {
	if this.listener == nil {
		return errors.New("must listen at first.")
	}
	return http.ServeTLS(this.listener, service, this.tlsPem, this.tlsKey)
}

func (this *S_Server) Close() {
	if this.listener != nil {
		this.listener.Close()
	}
}
