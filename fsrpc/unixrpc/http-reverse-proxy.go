/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: http 通过 UNIX SOCKET 转发反向代理客户端
@author: fanky
@version: 1.0
@date: 2025-06-08
**/

package unixrpc

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func NewHTTPReverseProxy(addr string, timeout time.Duration) *httputil.ReverseProxy {
	target, _ := url.Parse("http://unix")
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		DialContext: func(context.Context, string, string) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		},
	}
	return proxy
}
