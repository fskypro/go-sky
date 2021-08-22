/**
@copyright: fantasysky 2016
@brief: 实现格式化一个结构体
@author: fanky
@version: 1.0
@date: 2019-01-08
**/

package fmtex

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unsafe"

	"fsky.pro/fsos"
)

// 获取自定义显示字符串
func _getViewString(v reflect.Value) (ok bool, str string) {
	m := v.MethodByName("String")
	if !m.IsValid() {
		return
	}
	tm := m.Type()
	if tm.NumIn() != 0 || tm.NumOut() != 1 {
		return
	}
	if tm.Out(0).Kind() != reflect.String {
		return
	}
	// v.CanSet() 为 false 时，String 方法将不可被调用
	if v.CanSet() {
		outs := m.Call([]reflect.Value{})
		str = outs[0].Interface().(string)
		ok = true
		return
	}
	if v.CanAddr() {
		pv := reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr()))
		m = pv.Elem().MethodByName("String")
		outs := m.Call([]reflect.Value{})
		str = outs[0].Interface().(string)
		ok = true
		return
	}
	return
}

// -------------------------------------------------------------------
// temp writer
// -------------------------------------------------------------------
type s_Writer struct {
	w         *bufio.Writer
	prefix    string
	ident     string
	fmtcounts map[string]int

	layer  int
	idents string
	path   []string
}

func _newWriter(w io.Writer, opts *S_FmtOpts) *s_Writer {
	return &s_Writer{
		w:         bufio.NewWriter(w),
		prefix:    opts.Prefix,
		ident:     opts.Ident,
		fmtcounts: opts.FmtCounts,
		layer:     0,
		idents:    opts.Prefix,
	}
}

func (this *s_Writer) flush() {
	this.w.Flush()
}

func (this *s_Writer) incLayer() {
	this.layer++
	this.idents = this.prefix + strings.Repeat(this.ident, this.layer)
}

func (this *s_Writer) decLayer() {
	this.layer--
	this.idents = this.prefix + strings.Repeat(this.ident, this.layer)
}

func (this *s_Writer) enterPath(name string) {
	this.path = append(this.path, name)
}

func (this *s_Writer) leavePath() {
	this.path = this.path[0 : len(this.path)-1]
}

// 对于 array/slice/map 展示的元素个数
// -1：全部显示；0：不显示；>0 显示指定个数
func (this *s_Writer) getFmtCount() int {
	path := strings.Join(this.path, ".")
	count, ok := this.fmtcounts[path]
	if ok {
		return count
	}
	return -1
}

// ---------------------------------------------------------
func (this *s_Writer) writeStringf(s string, args ...interface{}) {
	this.w.WriteString(fmt.Sprintf(s, args...))
}

func (this *s_Writer) writeByte(b byte) {
	this.w.WriteByte(b)
}

func (this *s_Writer) writeIdents() {
	this.w.WriteString(this.idents)
}

func (this *s_Writer) writeEndline() {
	this.w.WriteString(fsos.Endline)
}

// ---------------------------------------------------------
func (this *s_Writer) writeValue(v reflect.Value, isTop bool) {
	if !isTop {
		if ok, str := _getViewString(v); ok {
			this.writeStringf(str)
			return
		}
	}

	if v.Type().Kind() != reflect.Ptr {
		printer, ok := _printers[v.Type().Kind()]
		if ok {
			printer(this, v)
		} else {
			_printOther(this, v)
		}
		return
	}

	// 空指针
	if !v.IsValid() || v.IsNil() {
		this.writeStringf("%#v", v)
		return
	}

	elem := v.Elem()
	pprinter, ok := _pprinters[elem.Type().Kind()]
	if ok {
		pprinter(this, elem)
	} else {
		_printPOther(this, elem)
	}
}

// -------------------------------------------------------------------
// module private
// -------------------------------------------------------------------

// 数组/切片
func _printArray(w *s_Writer, v reflect.Value) {
	isArray := v.Type().Kind() == reflect.Array
	ecount := v.Len()
	fcount := w.getFmtCount()

	// 没有元素
	if ecount == 0 {
		w.writeStringf("%#v", v)
		return
	}

	// fcount 指示不显示任何元素
	if fcount == 0 {
		if isArray {
			w.writeStringf("%v{...}", v.Type())
		} else {
			w.writeStringf("%v{...}(len=%d)", v.Type(), v.Len())
		}
		return
	}

	// 显示全部元素
	fmtAll := fcount < 0 || fcount >= ecount
	if fmtAll {
		fcount = ecount
	}

	// 第一个元素
	e := v.Index(0)

	// 元素为基础类型，则不对元素进行换行处理
	if _isBaseType(e.Type()) {
		// 显示全部
		if fmtAll {
			w.writeStringf("%#v", v)
			return
		}
		// 只显示部分
		w.writeStringf("%v{", v.Type())
		for i := 0; i < fcount; i++ {
			w.writeStringf("%#v, ", v.Index(i))
		}
		w.writeStringf("...}")
		if !isArray {
			w.writeStringf("(len=%d)", ecount)
		}
		return
	}

	// 写入类型
	w.writeStringf("%v{", v.Type()) // 写入类型
	w.writeEndline()                // 换行

	// 写入第一个元素
	w.incLayer() // 增加嵌套数
	w.writeIdents()
	w.writeValue(e, false)

	// 写入其他元素
	for i := 1; i < fcount; i++ {
		w.writeByte(',')
		w.writeEndline()
		w.writeIdents()
		w.writeValue(v.Index(i), false)
	}

	// 如果没显示完，则打印省略号
	if !fmtAll {
		w.writeByte(',')
		w.writeEndline()
		w.writeIdents()
		w.writeStringf("...")
	}

	w.writeEndline() // 换行
	w.decLayer()     // 减少嵌套
	w.writeIdents()  // 缩进
	w.writeByte('}') //数组结束

	// 如果是 slice 显示元素总数
	if !isArray && !fmtAll {
		w.writeStringf("(len=%d)", ecount)
	}
}

func _printPArray(w *s_Writer, v reflect.Value) {
	w.writeByte('&')
	_printArray(w, v)
}

// -----------------------------------------------
// 映射
func _printMap(w *s_Writer, v reflect.Value) {
	ecount := v.Len()
	fcount := w.getFmtCount()
	iter := v.MapRange()

	// 没有元素
	if !iter.Next() {
		w.writeStringf("%v{}", v.Type())
		return
	}

	// fcount 指示不显示任何元素
	if fcount == 0 {
		w.writeStringf("%v{...}(len=%d)", v.Type(), v.Len())
		return
	}

	// 显示全部元素
	fmtAll := fcount < 0 || fcount >= ecount
	if fmtAll {
		fcount = ecount
	}

	// value 的值为基础类型，则不对元素进行换行处理
	if _isBaseType(iter.Value().Type()) {
		if fmtAll {
			w.writeStringf("%#v", v)
			return
		}
		// 只显示部分元素
		w.writeStringf("%v{", v.Type())
		for fcount > 0 {
			w.writeStringf("%#v, ", iter.Value())
			fcount -= 1
			iter.Next()
		}
		w.writeStringf("...}(len=%d)", ecount)
		return
	}

	// 写入类型
	w.writeStringf("%v{", v.Type())
	w.writeEndline()

	// 写入第一个元素
	w.incLayer()    // 增加嵌套数
	w.writeIdents() // 第一个元素的缩进
	w.writeStringf("%#v: ", iter.Key())
	w.writeValue(iter.Value(), false)
	iter.Next()

	// 写入其他元素
	for fcount > 1 {
		w.writeByte(',')
		w.writeEndline()
		w.writeIdents()
		w.writeStringf("%#v: ", iter.Key())
		w.writeValue(iter.Value(), false)
		fcount -= 1
		iter.Next()
	}

	// 如果没显示完，则打印省略号
	if !fmtAll {
		w.writeByte(',')
		w.writeEndline()
		w.writeIdents()
		w.writeStringf("...")
	}

	w.writeEndline()
	w.decLayer()
	w.writeIdents()
	w.writeByte('}')

	// 如果没显示完所有元素，则显示总数
	if !fmtAll {
		w.writeStringf("(len=%d)", ecount)
	}
}

func _printPMap(w *s_Writer, v reflect.Value) {
	w.writeByte('&')
	_printMap(w, v)
}

// -----------------------------------------------
// 结构体
func _printStruct(w *s_Writer, v reflect.Value) {
	if v.NumField() == 0 {
		w.writeStringf("%#v", v)
		return
	}
	t := v.Type()
	w.writeStringf("%v{", t)
	w.incLayer()

	// 写成员
	var name string
	for i := 0; i < v.NumField(); i++ {
		name = t.Field(i).Name
		w.writeEndline()
		w.writeIdents()
		w.writeStringf("%s: ", name)
		w.enterPath(name)
		w.writeValue(v.Field(i), false)
		w.leavePath()
		w.writeByte(',')
	}

	w.decLayer()
	w.writeEndline()
	w.writeIdents()
	w.writeByte('}')
}

// 结构体指针
func _printPStruct(w *s_Writer, v reflect.Value) {
	w.writeByte('&')
	_printStruct(w, v)
}

// ------------------------------------------------
// 有符号整型
func _printNumber(w *s_Writer, v reflect.Value) {
	w.writeStringf("%v(%#v)", v.Type(), v)
}

// 有符号指针类型
func _printPNumber(w *s_Writer, v reflect.Value) {
	w.writeStringf("&%v(%#v)", v.Type(), v)
}

// ------------------------------------------------
// 无符号整型
func _printUNumber(w *s_Writer, v reflect.Value) {
	w.writeStringf("%v(%v=%#v)", v.Type(), v, v)
}

// 无符号指针类型
func _printPUNumber(w *s_Writer, v reflect.Value) {
	w.writeStringf("&%v(%v=%#v)", v.Type(), v, v)
}

// ------------------------------------------------
// 复数
func _printComplex(w *s_Writer, v reflect.Value) {
	w.writeStringf("%v%v", v.Type(), v)
}

// 复数指针
func _printPComplex(w *s_Writer, v reflect.Value) {
	w.writeStringf("&%v%v", v.Type(), v)
}

// ------------------------------------------------
// 其他类型
func _printOther(w *s_Writer, v reflect.Value) {
	w.writeStringf("%#v", v)
}

// 其他类型的指针类型
func _printPOther(w *s_Writer, v reflect.Value) {
	w.writeStringf("&%v(%#v)", v.Type(), v)
}

// ---------------------------------------------------------
func _isBaseType(t reflect.Type) bool {
	k := t.Kind()
	if k == reflect.Ptr {
		k = t.Elem().Kind()
	}
	_, ok := _baseTypes[k]
	return ok
}

var _printers map[reflect.Kind]func(*s_Writer, reflect.Value)  // 类型打印方法
var _pprinters map[reflect.Kind]func(*s_Writer, reflect.Value) // 指针类型方法
var _baseTypes map[reflect.Kind]interface{}                    // 基础类型，这些类型如果为 array、slice、map 的成员，则分列成员时，不换行

func init() {
	_printers = map[reflect.Kind]func(*s_Writer, reflect.Value){
		reflect.Array:      _printArray,
		reflect.Slice:      _printArray,
		reflect.Map:        _printMap,
		reflect.Struct:     _printStruct,
		reflect.Int:        _printNumber,
		reflect.Int8:       _printNumber,
		reflect.Int16:      _printNumber,
		reflect.Int32:      _printNumber,
		reflect.Int64:      _printNumber,
		reflect.Float32:    _printNumber,
		reflect.Float64:    _printNumber,
		reflect.Uint:       _printUNumber,
		reflect.Uint8:      _printUNumber,
		reflect.Uint16:     _printUNumber,
		reflect.Uint32:     _printUNumber,
		reflect.Uint64:     _printUNumber,
		reflect.Complex64:  _printComplex,
		reflect.Complex128: _printComplex,
	}
	_pprinters = map[reflect.Kind]func(*s_Writer, reflect.Value){
		reflect.Array:      _printPArray,
		reflect.Slice:      _printPArray,
		reflect.Map:        _printPMap,
		reflect.Struct:     _printPStruct,
		reflect.Int:        _printPNumber,
		reflect.Int8:       _printPNumber,
		reflect.Int16:      _printPNumber,
		reflect.Int32:      _printPNumber,
		reflect.Int64:      _printPNumber,
		reflect.Float32:    _printPNumber,
		reflect.Float64:    _printPNumber,
		reflect.Uint:       _printPUNumber,
		reflect.Uint8:      _printPUNumber,
		reflect.Uint16:     _printPUNumber,
		reflect.Uint32:     _printPUNumber,
		reflect.Uint64:     _printPUNumber,
		reflect.Complex64:  _printPComplex,
		reflect.Complex128: _printPComplex,
	}

	_baseTypes = map[reflect.Kind]interface{}{
		reflect.Bool:       nil,
		reflect.Int:        nil,
		reflect.Int8:       nil,
		reflect.Int16:      nil,
		reflect.Int32:      nil,
		reflect.Int64:      nil,
		reflect.Uint:       nil,
		reflect.Uint8:      nil,
		reflect.Uint16:     nil,
		reflect.Uint32:     nil,
		reflect.Uint64:     nil,
		reflect.Uintptr:    nil,
		reflect.Float32:    nil,
		reflect.Float64:    nil,
		reflect.Complex64:  nil,
		reflect.Complex128: nil,
		//reflect.String : nil,
		reflect.UnsafePointer: nil,
		reflect.Chan:          nil,
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
type S_FmtOpts struct {
	Prefix    string         // 前缀
	Ident     string         // 缩进字符串
	FmtCounts map[string]int // 显示指定数量的元素，格式为：map["aa.bb.cc"] = count（只对 array、slice、map 有效）
}

// SprintStruct 以初始化结构的格式，将一个结构体格式化为字符串，并写入流中
// 参数：
//  w: 流缓冲
//	st: 要格式化的结构体
//	prefix: 整个输出结构体的每一行的前缀
//	ident: 缩进字符串
func StreamStruct(w io.Writer, obj interface{}, opts *S_FmtOpts) {
	if opts == nil {
		opts = &S_FmtOpts{"    ", "    ", nil}
	}
	tobj := reflect.TypeOf(obj)
	if tobj == nil {
		w.Write([]byte("nil"))
		return
	}

	vobj := reflect.ValueOf(obj)
	if tobj.Kind() == reflect.Ptr {
		if !vobj.IsValid() || vobj.IsNil() {
			w.Write([]byte(fmt.Sprintf("(%v)(nil)", tobj)))
			return
		}
		if tobj.Elem().Kind() != reflect.Struct {
			opts.Prefix = ""
		}
	} else if tobj.Kind() != reflect.Struct {
		opts.Prefix = ""
	}

	writer := _newWriter(w, opts)
	writer.writeStringf(opts.Prefix)
	writer.writeValue(vobj, true)
	writer.flush()
}

// SprintStruct 以初始化结构的格式，将一个结构体格式化为字符串
// 参数：
//	st: 要格式化的结构体
//	prefix: 整个输出结构体的每一行的前缀
//	ident: 缩进字符串
func SprintStruct(obj interface{}, opts *S_FmtOpts) string {
	out := bytes.NewBuffer([]byte{})
	StreamStruct(out, obj, opts)
	return out.String()
}
