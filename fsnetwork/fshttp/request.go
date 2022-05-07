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
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// 引用私有函数，用于请求文件
//go:linkname serveFile net/http.serveFile
func serveFile(http.ResponseWriter, *http.Request, http.FileSystem, string, bool)

// -------------------------------------------------------------------
// S_Request
// -------------------------------------------------------------------
type S_Request struct {
	Service   I_Service
	MatchPath string
	W         http.ResponseWriter
	R         *http.Request
	buff      strings.Builder
	jsonObj   interface{}
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

func (this *S_Request) unmarshalValues(obj interface{}, get func(string) string) error {
	tobj := reflect.TypeOf(obj)
	if tobj == nil || tobj.Kind() != reflect.Ptr {
		return fmt.Errorf("output obj argument must be a un-nil struct object pointer")
	}
	vobj := reflect.ValueOf(obj)
	if vobj.IsNil() {
		return fmt.Errorf("output obj argument must be a un-nil struct object pointer")
	}
	tobj = tobj.Elem()
	vobj = vobj.Elem()
L:
	for i := 0; i < tobj.NumField(); i++ {
		sfield := tobj.Field(i)
		vfield := vobj.Field(i)
		tfield := vfield.Type()
		tag := sfield.Tag.Get("urlkey")
		if tag == "" {
			tag = sfield.Name
		}
		var value interface{}
		svalue := get(tag)
		if svalue == "" {
			continue L
		}
		switch ft := tfield.Kind(); ft {
		case reflect.String:
			value = svalue
		case reflect.Int:
			if v, err := strconv.Atoi(svalue); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = v
			}
		case reflect.Int8:
			if v, err := strconv.Atoi(svalue); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = int8(v)
			}
		case reflect.Int16:
			if v, err := strconv.Atoi(svalue); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = int16(v)
			}
		case reflect.Int32:
			if v, err := strconv.ParseInt(svalue, 10, 32); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = int32(v)
			}
		case reflect.Int64:
			if v, err := strconv.ParseInt(svalue, 10, 64); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = int64(v)
			}
		case reflect.Uint:
			if v, err := strconv.ParseUint(svalue, 10, 64); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = uint(v)
			}
		case reflect.Uint8:
			if v, err := strconv.ParseUint(svalue, 10, 8); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = uint8(v)
			}
		case reflect.Uint16:
			if v, err := strconv.ParseUint(svalue, 10, 16); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = uint16(v)
			}
		case reflect.Uint32:
			if v, err := strconv.ParseUint(svalue, 10, 32); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = uint32(v)
			}
		case reflect.Uint64:
			if v, err := strconv.ParseUint(svalue, 10, 64); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = uint64(v)
			}
		case reflect.Bool:
			svalue = strings.ToLower(svalue)
			if svalue == "1" || svalue == "true" {
				value = true
			} else if svalue == "0" || svalue == "false" {
				value = false
			} else {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			}
		case reflect.Float32:
			if v, err := strconv.ParseFloat(svalue, 64); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = float32(v)
			}
		case reflect.Float64:
			if v, err := strconv.ParseFloat(svalue, 64); err != nil {
				return fmt.Errorf("value of form key %q is not a %v", tag, tfield.Kind())
			} else {
				value = v
			}
		default:
			return fmt.Errorf("unsupport value type %v", tfield)
		}
		reflect.NewAt(tfield, unsafe.Pointer(vobj.UnsafeAddr()+sfield.Offset)).Elem().Set(reflect.ValueOf(value))
	}
	return nil
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 远程地址
func (this *S_Request) RemoteHost() string {
	return this.R.Header.Get("RemoteHost")
}

// 远程端口
func (this *S_Request) RemotePort() int {
	port, _ := strconv.Atoi(this.R.Header.Get("RemotePort"))
	return port
}

// 将 get 请求的参数反序列化到 obj 对象
// 对象成员的 tag 标记为：urlkey
func (this *S_Request) UnmarshalQuerys(obj interface{}) error {
	values := this.R.URL.Query()
	return this.unmarshalValues(obj, func(key string) string { return values.Get(key) })
}

// 将表单参数反序列化为 obj 对象
// 对象成员的 tag 标记为：urlkey
func (this *S_Request) UnmarshalForms(obj interface{}) error {
	return this.unmarshalValues(obj, func(key string) string { return this.R.Form.Get(key) })
}

func (this *S_Request) UnmarshalPostForms(obj interface{}) error {
	return this.unmarshalValues(obj, func(key string) string { return this.R.PostForm.Get(key) })
}

func (this *S_Request) UnmarshalPostBody(obj interface{}) error {
	defer this.R.Body.Close()
	body, err := ioutil.ReadAll(this.R.Body)
	if err != nil {
		return err
	}

	keyValues := map[string]string{}
	items := strings.Split(string(body), "&")
	for _, item := range items {
		kv := strings.Split(item, "=")
		if len(kv) == 2 {
			keyValues[kv[0]] = kv[1]
		}
	}
	return this.unmarshalValues(obj, func(key string) string { return keyValues[key] })
}

func (this *S_Request) ReadBody() ([]byte, error) {
	defer this.R.Body.Close()
	return ioutil.ReadAll(this.R.Body)
}

// ---------------------------------------------------------
// 直接写入回复字符串
func (this *S_Request) WriteRspString(str string) {
	this.buff.WriteString(str)
}

func (this *S_Request) WriteRspStringf(str string, args ...interface{}) {
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
func (this *S_Request) WriteRspJsonObject(obj interface{}) {
	this.jsonObj = obj
}

// ---------------------------------------------------------
// 取消返回
func (this *S_Request) Cancel() {
	this.cancel = true
}

func (this *S_Request) CancelStringf(str string, args ...interface{}) {
	this.cancel = true
	fmt.Fprintf(this.W, str, args...)
}

func (this *S_Request) CancelBytes(bs []byte) {
	this.cancel = true
	this.W.Write(bs)
}

func (this *S_Request) CancelJsonObject(obj interface{}) {
	this.cancel = true
	bs, _ := json.Marshal(obj)
	this.W.Write(bs)
}

// ---------------------------------------------------------
// 回复客户端
func (this *S_Request) Response() bool {
	if !this.cancel {
		this.W.Write([]byte(this.buff.String()))
	}
	return !this.cancel
}

func (this *S_Request) ResponseCrossDomain() bool {
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
	serveFile(this.W, this.R, http.Dir(webroot), path.Clean(file), true)
}

func (this *S_Request) Redirect(url string, code int) {
	this.cancel = true
	http.Redirect(this.W, this.R, url, code)
}
