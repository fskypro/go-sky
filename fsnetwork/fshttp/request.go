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
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
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

func (this *S_Request) unmarshalValues(obj interface{}, get func(string) (string, bool)) error {
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
		if sfield.Anonymous {
			continue
		}
		if strings.HasSuffix(sfield.Name, "__") {
			// __ 结尾的为默认值成员
			continue
		}

		vfield := vobj.Field(i)
		tfield := vfield.Type()

		tag := sfield.Tag.Get("urlkey")
		if tag == "" {
			tag = sfield.Tag.Get("json")
			if tag == "" {
				tag = sfield.Name
			}
		}
		var value any
		svalue, ok := get(tag)
		if !ok {
			// 没有传入值
			defValue := vobj.FieldByName(sfield.Name + "__") // 原成员名称后面加 “__” 表示默认值
			if !defValue.IsValid() {                         // 默认值也不设置，则表示该参数是必传参数
				return newNoReqArgError(tag)
			}
			if defValue.Type() != tfield {
				return fmt.Errorf("the type of default value for member is not match, name=%q, default-name=%q", sfield.Name, sfield.Name+"__")
			}
			reflect.NewAt(tfield, unsafe.Pointer(vobj.UnsafeAddr()+sfield.Offset)).Elem().Set(defValue)
			continue L
		} else if svalue == "" {
			continue L
		}
		switch ft := tfield.Kind(); ft {
		case reflect.String:
			value = svalue
		case reflect.Int:
			if v, err := strconv.Atoi(svalue); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = v
			}
		case reflect.Int8:
			if v, err := strconv.Atoi(svalue); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = int8(v)
			}
		case reflect.Int16:
			if v, err := strconv.Atoi(svalue); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = int16(v)
			}
		case reflect.Int32:
			if v, err := strconv.ParseInt(svalue, 10, 32); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = int32(v)
			}
		case reflect.Int64:
			if v, err := strconv.ParseInt(svalue, 10, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = int64(v)
			}
		case reflect.Uint:
			if v, err := strconv.ParseUint(svalue, 10, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint(v)
			}
		case reflect.Uint8:
			if v, err := strconv.ParseUint(svalue, 10, 8); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint8(v)
			}
		case reflect.Uint16:
			if v, err := strconv.ParseUint(svalue, 10, 16); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint16(v)
			}
		case reflect.Uint32:
			if v, err := strconv.ParseUint(svalue, 10, 32); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = uint32(v)
			}
		case reflect.Uint64:
			if v, err := strconv.ParseUint(svalue, 10, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
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
				return newReqArgTypeError(tag, tfield)
			}
		case reflect.Float32:
			if v, err := strconv.ParseFloat(svalue, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = float32(v)
			}
		case reflect.Float64:
			if v, err := strconv.ParseFloat(svalue, 64); err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = v
			}
		case _timeKind:
			v, err := time.ParseInLocation(time.DateTime, svalue, time.Local)
			if err != nil {
				v, err = time.ParseInLocation(time.DateOnly, svalue, time.Local)
			}
			if err != nil {
				return newReqArgTypeError(tag, tfield)
			} else {
				value = v
			}
		default:
			return fmt.Errorf("unsupport value type %v", tfield)
		}
		vvalue := reflect.ValueOf(value).Convert(tfield)
		reflect.NewAt(tfield, unsafe.Pointer(vobj.UnsafeAddr()+sfield.Offset)).Elem().Set(vvalue)
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

// 获取最前端地址（要求转发服务器添加 X-Forwarded-For 头）
func (this *S_Request) XForwardedFor() string {
	addr := this.R.Header.Get("X-Forwarded-For")
	if addr == "" {
		return this.RemoteHost()
	}
	return addr
}

// 远程端口
func (this *S_Request) RemotePort() int {
	port, _ := strconv.Atoi(this.R.Header.Get("RemotePort"))
	return port
}

// ---------------------------------------------------------
// 将 get 请求的参数反序列化到 obj 对象
// 对象成员的 tag 标记为：urlkey
func (this *S_Request) UnmarshalQuerys(obj interface{}) error {
	values := this.R.URL.Query()
	return this.unmarshalValues(obj, func(key string) (string, bool) {
		return values.Get(key), values.Has(key)
	})
}

// 将表单参数反序列化为 obj 对象
// 对象成员的 tag 标记为：urlkey
func (this *S_Request) UnmarshalForms(obj interface{}) error {
	return this.unmarshalValues(obj, func(key string) (string, bool) {
		return this.R.Form.Get(key), this.R.Form.Has(key)
	})
}

func (this *S_Request) UnmarshalPostForms(obj interface{}) error {
	return this.unmarshalValues(obj, func(key string) (string, bool) {
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
	// 对象类型判断
	var tobj = reflect.TypeOf(obj)
	if tobj == nil {
		return fmt.Errorf("output object mustn't be nil value")
	}
	if tobj.Kind() != reflect.Ptr {
		return fmt.Errorf("output object must be a pointer")
	}
	var vobj = reflect.ValueOf(obj)

	// 定位到终极类型
	for {
		if tobj.Kind() == reflect.Ptr {
			if vobj.IsNil() {
				return fmt.Errorf("output object's pointer mustn't be nil")
			}
			tobj = tobj.Elem()
			vobj = vobj.Elem()
		} else {
			break
		}
	}

	// 获取请求参数内容
	body, err := ioutil.ReadAll(this.R.Body)
	if err != nil {
		return err
	}
	defer this.R.Body.Close()

	// 如果是纯 map 或 slice，则直接反序列
	if tobj.Kind() == reflect.Map || tobj.Kind() == reflect.Slice {
		err := json.Unmarshal(body, obj)
		if err != nil {
			return fmt.Errorf("unmarshal json data to %v fail", tobj)
		}
		return nil
	}

	// 查找出传入 json 的所有顶层 key
	jmap := map[string]json.RawMessage{}
	if err := json.Unmarshal(body, &jmap); err != nil {
		return fmt.Errorf("parse json data fail, %v", err)
	}
	if err := json.Unmarshal(body, obj); err != nil {
		return fmt.Errorf("unmarshal json to object fail, %v", err)
	}

	// 遍历结构体所有成员与成员值
	type temp struct {
		rv     reflect.Value
		field  reflect.StructField
		vfield reflect.Value
	}
	var trivalStruct func(reflect.Value, map[string]*temp)
	trivalStruct = func(rv reflect.Value, members map[string]*temp) {
		if !rv.IsValid() {
			return
		}
		rt := rv.Type()
		if rt.Kind() != reflect.Struct {
			return
		}
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			if field.Tag.Get("json") == "-" && !strings.HasSuffix(field.Name, "__") {
				continue
			}
			vfield := rv.Field(i)
			if !field.Anonymous {
				tag := field.Name
				if unicode.IsLower(rune(tag[0])) {
					continue
				}
				if strings.HasSuffix(tag, "__") {
					members[tag] = &temp{rv, field, vfield}
					continue
				} else if v, ok := field.Tag.Lookup("json"); ok {
					tag = v
				}
				members[tag] = &temp{rv, field, vfield}
				continue
			}
			// 匿名结构体
			tfield := field.Type
			for tfield.Kind() == reflect.Ptr {
				if !vfield.IsValid() {
					break
				}
				if vfield.IsNil() {
					break
				}
				tfield = tfield.Elem()
				vfield = vfield.Elem()
			}
			// 继承结构体，继续往上层遍历
			if tfield.Kind() == reflect.Struct {
				trivalStruct(vfield, members)
			}
		}
	}
	members := map[string]*temp{}
	trivalStruct(vobj, members)

	// 检查传入 json 数据中，是否有缺少的成员，
	// 如果有，则查找同名并以 __ 结尾的成员值，作为其默认值
	// 如果找不到默认值，则提示传入 json 中缺少字段
	for k, m := range members {
		if strings.HasSuffix(m.field.Name, "__") {
			// 排除默认值成员
			continue
		}
		if _, ok := jmap[k]; ok {
			// json 提供了该成员的值
			continue
		}
		// json 中没有传入值，则查找默认值
		var def = members[m.field.Name+"__"]
		if def == nil || !def.vfield.IsValid() {
			// 传入 json 中没有对应的字段，并且缺少默认值
			return newNoReqArgError(k)
		}
		if def.field.Type != m.field.Type {
			// 默认值与 json 映射字段值类型不一致
			return fmt.Errorf("the type of default member %[1]q is not match the type of member %[2]q, "+
				"typeof(%[1]s)=%[3]v, typeof(%[2]s)=%[4]v",
				def.field.Name, m.field.Name, def.field.Type, m.field.Type)
		}
		// 将默认值拷贝给对应字段成员
		reflect.NewAt(m.field.Type, unsafe.Pointer(m.rv.UnsafeAddr()+m.field.Offset)).Elem().Set(def.vfield)
	}
	return nil
}

// ---------------------------------------------------------
func (this *S_Request) ReadBody() ([]byte, error) {
	defer this.R.Body.Close()
	return ioutil.ReadAll(this.R.Body)
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
func (this *S_Request) Response() error {
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
