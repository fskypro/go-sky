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
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
)

type I_Handler interface {
	OnRequest(*S_Request)
}

// -------------------------------------------------------------------
// handler
// -------------------------------------------------------------------
type s_Handler struct {
	path string
	h1   func(http.ResponseWriter, *http.Request)
	h2   func(*S_Request)
	h3   I_Handler
}

func (this *s_Handler) call(service *S_Service, w http.ResponseWriter, r *http.Request) {
	if this.h1 != nil {
		this.h1(w, r)
	}
	if this.h2 != nil {
		if service.I_Service != nil {
			this.h2(newRequest(service.I_Service, this.path, w, r))
		} else {
			this.h2(newRequest(service, this.path, w, r))
		}
	}
	if this.h3 != nil {
		this.h3.OnRequest(newRequest(service, this.path, w, r))
	}
}

type s_ReHandler struct {
	s_Handler
	repath *regexp.Regexp
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
	locker     sync.RWMutex
	handlers   map[string]*s_Handler
	rehandlers map[string]*s_ReHandler
	OnError    func(string, *http.Request)
	OnPanic    func(string, *http.Request)
}

func NewService() *S_Service {
	service := &S_Service{
		handlers:   make(map[string]*s_Handler),
		rehandlers: make(map[string]*s_ReHandler),
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
// 添加 url 路径请求处理函数，函数格式为：
//	func(http.ResponseWriter, *http.Request)
// 如果 path 以 [re] 开头，则表示该路径为正则表达式
func (this *S_Service) AddHandler(path string, f func(http.ResponseWriter, *http.Request)) error {
	this.locker.Lock()
	defer this.locker.Unlock()
	if !strings.HasPrefix(path, "[re]") {
		if _, ok := this.handlers[path]; ok {
			return fmt.Errorf("http path %q has been exists.", path)
		}
		this.handlers[path] = &s_Handler{path: path, h1: f}
		return nil
	}
	if _, ok := this.rehandlers[path]; ok {
		return fmt.Errorf("regexp http path %q has been exists", path)
	}
	re, err := regexp.Compile(path[4:])
	if err != nil {
		return fmt.Errorf("error regexp pattern path %q", path)
	}
	this.rehandlers[path] = &s_ReHandler{
		s_Handler: s_Handler{path: path, h1: f},
		repath:    re,
	}
	return nil
}

// 添加 url 路径请求处理函数，函数格式为：
//	func(*S_Request)
// 如果 path 以 [re] 开头，则表示该路径为正则表达式
func (this *S_Service) AddHandler2(path string, f func(*S_Request)) error {
	this.locker.Lock()
	defer this.locker.Unlock()
	if !strings.HasPrefix(path, "[re]") {
		if _, ok := this.handlers[path]; ok {
			return fmt.Errorf("http path %q has been exists.", path)
		}
		this.handlers[path] = &s_Handler{path: path, h2: f}
		return nil
	}
	if _, ok := this.rehandlers[path]; ok {
		return fmt.Errorf("regexp http path %q has been exists", path)
	}
	re, err := regexp.Compile(path[4:])
	if err != nil {
		return fmt.Errorf("error regexp pattern path %q", path)
	}
	this.rehandlers[path] = &s_ReHandler{
		s_Handler: s_Handler{path: path, h2: f},
		repath:    re,
	}
	return nil
}

// 添加 url 路径请求处理对象，处理对象必须实现 I_Handler 接口
// 如果 path 以 [re] 开头，则表示该路径为正则表达式
func (this *S_Service) AddHandler3(path string, h I_Handler) error {
	this.locker.Lock()
	defer this.locker.Unlock()
	if !strings.HasPrefix(path, "[re]") {
		if _, ok := this.handlers[path]; ok {
			return fmt.Errorf("http path %q has been exists.", path)
		}
		this.handlers[path] = &s_Handler{path: path, h3: h}
		return nil
	}
	if _, ok := this.rehandlers[path]; ok {
		return fmt.Errorf("regexp http path %q has been exists", path)
	}
	re, err := regexp.Compile(path[4:])
	if err != nil {
		return fmt.Errorf("error regexp pattern path %q", path)
	}
	this.rehandlers[path] = &s_ReHandler{
		s_Handler: s_Handler{path: path, h3: h},
		repath:    re,
	}
	return nil
}

// ---------------------------------------------------------
func (this *S_Service) RemoveHandler(path string) bool {
	this.locker.Lock()
	defer this.locker.Unlock()
	_, ok1 := this.handlers[path]
	_, ok2 := this.rehandlers[path]
	delete(this.handlers, path)
	delete(this.rehandlers, path)
	return ok1 || ok2
}

// ---------------------------------------------------------
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

	this.locker.RLock()
	defer this.locker.RUnlock()

	// 精确匹配
	handler, ok := this.handlers[r.URL.Path]
	if ok {
		handler.call(this, w, r)
		return
	}

	// 正则匹配
	for _, hdl := range this.rehandlers {
		if hdl.repath.MatchString(r.URL.Path) {
			hdl.call(this, w, r)
			return
		}
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
