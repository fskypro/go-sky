/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: request wrapper
@author: fanky
@version: 1.0
@date: 2022-01-29
**/

package fshttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// -------------------------------------------------------------------
// inner inits
// -------------------------------------------------------------------
var _timeKind reflect.Kind

func init() {
	var t time.Time
	_timeKind = reflect.TypeOf(t).Kind()
}

// -------------------------------------------------------------------
// S_Request
// -------------------------------------------------------------------
type S_Request struct {
	Service   I_Service
	MatchPath string
	W         http.ResponseWriter
	R         *http.Request
	buff      strings.Builder
	jsonObj   any
	cancel    bool
}

func newRequest(service I_Service, mpath string, w http.ResponseWriter, r *http.Request) *S_Request {
	return &S_Request{
		Service:   service,
		MatchPath: mpath,
		W:         w,
		R:         r,
		cancel:    false,
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 远程地址
func (this *S_Request) RemoteHost() string {
	xxf := this.XForwardedFors()
	for len(xxf) > 0 {
		first := xxf[0]
		if first == "127.0.0.1" || first == "::1" {
			xxf = xxf[1:]
		} else {
			break
		}
	}
	if len(xxf) > 0 {
		return xxf[0]
	}
	addr, _, _ := net.SplitHostPort(this.R.RemoteAddr)
	return addr
}

// 获取最前端地址（要求转发服务器添加 X-Forwarded-For 头）
func (this *S_Request) XForwardedFor() string {
	return strings.Join(this.XForwardedFors(), "->")
}

func (this *S_Request) XForwardedFors() []string {
	if xff := this.R.Header.Get("X-Forwarded-For"); xff != "" {
		addrs := strings.Split(xff, ",")
		for i, addr := range addrs {
			addrs[i] = strings.TrimSpace(addr)
		}
		return addrs
	}

	// 备用方案：X-Real-IP
	if xri := this.R.Header.Get("X-Real-IP"); xri != "" {
		return []string{strings.TrimSpace(xri)}
	}

	// 默认 fallback：直接从连接获取
	ip, _, _ := net.SplitHostPort(this.R.RemoteAddr)
	return []string{ip}
}

// 远程端口
func (this *S_Request) RemotePort() int {
	port, _ := strconv.Atoi(this.R.Header.Get("RemotePort"))
	return port
}

// ---------------------------------------------------------
// 将 get 请求的参数反序列化到 obj 对象
// 对象成员的 tag 标记为：urlkey
func (this *S_Request) UnmarshalQuerys(obj any) error {
	values := this.R.URL.Query()
	return UnmarshalRequestArgs(obj, func(key string) (string, bool) {
		return values.Get(key), values.Has(key)
	})
}

// 将表单参数反序列化为 obj 对象
// 对象成员的 tag 标记为：urlkey
func (this *S_Request) UnmarshalForms(obj any) error {
	return UnmarshalRequestArgs(obj, func(key string) (string, bool) {
		return this.R.Form.Get(key), this.R.Form.Has(key)
	})
}

func (this *S_Request) UnmarshalPostForms(obj any) error {
	return UnmarshalRequestArgs(obj, func(key string) (string, bool) {
		return this.R.PostForm.Get(key), this.R.PostForm.Has(key)
	})
}

// 解释 POST json 参数到结构体对象，如果传入参数中不存在结构体中的成员则返回错误，定义了默认值的例外
// 如：
//
//	args := &struct {
//	   Value1 string `json:"value1"`
//	   Value2 string `json:"value2"`
//	   Value1__ strring
//	}{
//	   Value1__: "default value",
//	}
//
// 用以下 json 字符串解码时，会提示缺少参数 value1：{"value2": "xxxx"}
// 用以下 json 字符串解码时，不会有问题，并且 args.Value2 == "default value"：{"value": "xxxx"}
func (this *S_Request) UnmarshalPostJsonBody(obj any) error {
	// 获取请求参数内容
	body, err := io.ReadAll(this.R.Body)
	if err != nil {
		return err
	}
	defer this.R.Body.Close()
	return UnmarshalPostJsonBody(body, obj)
}

// ---------------------------------------------------------
func (this *S_Request) ReadBody() ([]byte, error) {
	defer this.R.Body.Close()
	return io.ReadAll(this.R.Body)
}

// ---------------------------------------------------------
// 写入响应头
func (this *S_Request) WriteResponseCode(code int) {
	this.W.WriteHeader(code)
}

// 直接写入回复字符串
func (this *S_Request) WriteRspString(str string) {
	this.buff.WriteString(str)
}

func (this *S_Request) WriteRspStringf(str string, args ...any) {
	this.buff.WriteString(fmt.Sprintf(str, args...))
}

// 直接写入回复字符
func (this *S_Request) WriteRspByte(c byte) {
	this.buff.WriteByte(c)
}

// 直接写入回复字节数据
func (this *S_Request) WriteRspBytes(bs []byte) {
	this.buff.Write(bs)
}

// 写入对象转换为 json 数据
func (this *S_Request) WriteRspJsonObject(obj any) {
	this.jsonObj = obj
}

// ---------------------------------------------------------
// 取消返回
func (this *S_Request) Cancel() {
	this.cancel = true
}

func (this *S_Request) CancelStringf(str string, args ...any) {
	this.cancel = true
	fmt.Fprintf(this.W, str, args...)
}

func (this *S_Request) CancelBytes(bs []byte) {
	this.cancel = true
	this.W.Write(bs)
}

func (this *S_Request) CancelJsonObject(obj any) {
	this.cancel = true
	bs, _ := json.Marshal(obj)
	this.W.Write(bs)
}

// ---------------------------------------------------------
// 回复客户端
func (this *S_Request) Response() error {
	defer func() {
		this.jsonObj = nil // 释放引用
	}()
	if this.cancel {
		return nil
	}
	if this.jsonObj != nil {
		this.W.Header().Set("Content-Type", "application/json")
		bs, err := json.Marshal(this.jsonObj)
		if err != nil {
			return fmt.Errorf("marshal response obj(type of %v) fail, %v", reflect.TypeOf(this.jsonObj), err)
		}
		_, err = this.W.Write(bs)
		if err != nil {
			return fmt.Errorf("wirte data to response stream fail, %v", err)
		}
		return nil
	}
	_, err := this.W.Write([]byte(this.buff.String()))
	return err
}

func (this *S_Request) ResponseCrossDomain() bool {
	defer func() {
		this.jsonObj = nil // 释放引用
	}()
	if this.cancel {
		return false
	}
	this.W.Header().Set("Access-Control-Allow-Origin", "*")
	if this.jsonObj != nil {
		this.W.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		this.W.Header().Set("Content-Type", "application/json")
		bs, _ := json.Marshal(this.jsonObj)
		this.W.Write(bs)
		return true
	}
	this.W.Write([]byte(this.buff.String()))
	return true
}

func (this *S_Request) ResponseFile(webroot string, file string) {
	this.cancel = true
	mimeType := GetMimeType(path.Ext(file))
	if mimeType != "" {
		this.W.Header().Set("Content-Type", mimeType)
	}
	fullPath := filepath.Join(webroot, filepath.Clean(file))
	http.ServeFile(this.W, this.R, fullPath)
}

func (this *S_Request) Redirect(url string, code int) {
	this.cancel = true
	http.Redirect(this.W, this.R, url, code)
}

// 找到资源
func (this *S_Request) RedirectFound(url string) {
	this.cancel = true
	http.Redirect(this.W, this.R, url, http.StatusSeeOther)
}

// 找不到资源
func (this *S_Request) ResponseNotFound() {
	// this.W.WriteHeader(404)
	http.NotFound(this.W, this.R)
}
