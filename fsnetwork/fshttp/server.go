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
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type S_Server struct {
	*http.Server

	host string
	port int
}

func NewServer(host string, port int) *S_Server {
	return &S_Server{
		Server: &http.Server{Addr: fmt.Sprintf("%s:%d", host, port)},
		host:   host,
		port:   port,
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
		Server: &http.Server{Addr: addr},
		host:   hp[0],
		port:   port,
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

// -------------------------------------------------------------------
// 添加证书文件，domain 为域名，可以不指定(传入空字符串)
func (this *S_Server) AddTLSFiles(domain string, pem, key string) error {
	cert, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		return err
	}
	if this.TLSConfig == nil {
		this.TLSConfig = new(tls.Config)
	}
	if this.TLSConfig.Certificates == nil {
		this.TLSConfig.Certificates = []tls.Certificate{}
	}
	this.TLSConfig.Certificates = append(this.TLSConfig.Certificates, cert)
	if domain != "" {
		if this.TLSConfig.NameToCertificate == nil {
			this.TLSConfig.NameToCertificate = map[string]*tls.Certificate{}
		}
		this.TLSConfig.NameToCertificate[domain] = &cert
	}
	return nil
}

// -------------------------------------------------------------------
func (this *S_Server) ServeTLS(ln net.Listener) error {
	return this.Server.ServeTLS(ln, "", "")
}

func (this *S_Server) ListenAndServe(service I_Service) error {
	this.Handler = service
	return this.Server.ListenAndServe()
}

func (this *S_Server) ListenAndServeTLS(service I_Service) error {
	addr := this.Addr()
	if addr == "" {
		addr = ":https"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	this.Handler = service

	defer ln.Close()
	return this.ServeTLS(ln)
}

func (this *S_Server) ServeAuto(service I_Service) error {
	if this.TLSConfig != nil &&
		(len(this.TLSConfig.Certificates) > 0 ||
			this.TLSConfig.GetCertificate != nil) {
		return this.ListenAndServeTLS(service)
	}
	return this.ListenAndServe(service)
}
