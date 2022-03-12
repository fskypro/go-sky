/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: service
@author: fanky
@version: 1.0
@date: 2021-12-04
**/

package fshttp

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
)

type s_Handler struct {
	h1 func(http.ResponseWriter, *http.Request)
	h2 func(*S_Request)
}

func (this *s_Handler) call(service *S_Service, w http.ResponseWriter, r *http.Request) {
	if this.h1 != nil {
		this.h1(w, r)
	}
	if this.h2 != nil {
		if service.I_Service != nil {
			this.h2(newRequest(service.I_Service, w, r))
		} else {
			this.h2(newRequest(service, w, r))
		}
	}
}

// -------------------------------------------------------------------
// Service interface
// -------------------------------------------------------------------
type I_Service interface {
	http.Handler
}

// -------------------------------------------------------------------
// Service base
// -------------------------------------------------------------------
type S_Service struct {
	I_Service
	handlers map[string]*s_Handler
	OnError  func(string, *http.Request)
	OnPanic  func(string, *http.Request)
}

func NewService() *S_Service {
	service := &S_Service{
		handlers: make(map[string]*s_Handler),
	}
	service.OnError = service.onError
	service.OnPanic = service.onPanic
	return service
}

// ---------------------------------------------------------
// private
// ---------------------------------------------------------
func (this *S_Service) getRemoteAddr(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("XRealIP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("XForwardedFor"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func (this *S_Service) onError(err string, r *http.Request) {
	fmt.Println(err)
}

func (this *S_Service) onPanic(stack string, r *http.Request) {
	fmt.Println(stack)
}

// ---------------------------------------------------------
// public
// ---------------------------------------------------------
func (this *S_Service) AddHandler(path string, f func(http.ResponseWriter, *http.Request)) error {
	if _, ok := this.handlers[path]; ok {
		return fmt.Errorf("http path %q has been exists.", path)
	}
	this.handlers[path] = &s_Handler{h1: f}
	return nil
}

func (this *S_Service) AddHandler2(path string, f func(*S_Request)) error {
	if _, ok := this.handlers[path]; ok {
		return fmt.Errorf("http path %q has been exists.", path)
	}
	this.handlers[path] = &s_Handler{h2: f}
	return nil
}

func (this *S_Service) RemoveHandler(path string) bool {
	_, ok := this.handlers[path]
	if ok {
		delete(this.handlers, path)
	}
	return ok
}

func (this *S_Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if this.OnPanic != nil {
			blines := debug.Stack()
			lines := bytes.Split(blines, []byte("\n"))
			trace := bytes.Join(lines[7:], []byte("\n"))
			this.OnPanic(fmt.Sprintf("panic for handling request(path=%q) from client(%s): %v.\n%s",
				r.URL.Path, r.Header.Get("ClientAddr"), err, trace), r)
		}
	}()
	if r.Header.Get("ClientAddr") != "" {
		r.RemoteAddr = r.Header.Get("ClientAddr")
	}
	host, port, _ := net.SplitHostPort(r.RemoteAddr)
	r.Header.Add("RemoteHost", host)
	r.Header.Add("RemotePort", port)

	// 精确匹配
	handler, ok := this.handlers[r.URL.Path]
	if ok {
		handler.call(this, w, r)
		return
	}

	// 往上匹配
	path := r.URL.Path
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	segs := strings.Split(path, "/")
	idx := len(segs)
	for idx > 1 {
		idx--
		parent := strings.Join(segs[:idx], "/") + "/"
		if handler, ok := this.handlers[parent]; ok {
			handler.call(this, w, r)
			return
		}
	}
	if this.OnError != nil {
		this.OnError(fmt.Sprintf("no handler for handling request(RemoteAddr=%s) url: %s",
			r.RemoteAddr, r.URL.Path), r)
	}
	return
}
