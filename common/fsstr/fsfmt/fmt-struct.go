/**
@copyright: fantasysky 2016
@brief: 实现格式化一个结构体
@author: fanky
@version: 1.0
@date: 2019-01-08
**/

package fsfmt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unsafe"

	"fsky.pro/fsdef"
)

// 获取自定义显示字符串
func (*s_Writer) getViewString(v reflect.Value) (ok bool, str string) {
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

// 指定结构成员是否有指定 tag 标记
func (*s_Writer) hasFieldTag(tags string, tag string) bool {
	for _, t := range strings.Split(tags, "|") {
		if t == tag {
			return true
		}
	}
	return false
}

// -------------------------------------------------------------------
// temp writer
// -------------------------------------------------------------------
type s_Writer struct {
	w         *bufio.Writer
	prefix    string
	ident     string
	fmtCounts map[string]int

	layer  int
	idents string
	path   []string

	// 每进入一层结构体，如果引用的是该结构体的指针，则将指针入栈，离开结构体时出栈
	// 记录下结构体的嵌套结构，用于查找出循环引用的结构体，忽略打印父结构体，以免产生死循环
	structPtrs []reflect.Value
}

func _newWriter(w io.Writer, opts *S_FmtOpts) *s_Writer {
	if opts.Idents < 2 {
		opts.Idents = 2
	}
	return &s_Writer{
		w:          bufio.NewWriter(w),
		prefix:     opts.Prefix,
		ident:      strings.Repeat(" ", opts.Idents),
		fmtCounts:  opts.FmtCounts,
		layer:      0,
		idents:     opts.Prefix,
		path:       make([]string, 0),
		structPtrs: make([]reflect.Value, 0),
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
	count, ok := this.fmtCounts[path]
	if ok {
		return count
	}
	return -1
}

// 进入结构体指针
func (this *s_Writer) enterStructPtr(ptr reflect.Value) {
	this.structPtrs = append(this.structPtrs, ptr)
}

// 离开结构体指针
func (this *s_Writer) leaveStructPtr() {
	this.structPtrs = this.structPtrs[0 : len(this.structPtrs)-1]
}

func (this *s_Writer) isInStruct(ptr reflect.Value) bool {
	for _, p := range this.structPtrs {
		if p.UnsafePointer() == ptr.UnsafePointer() {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------
func (this *s_Writer) writeStringf(s string, args ...interface{}) {
	this.w.WriteString(fmt.Sprintf(s, args...))
}

func (this *s_Writer) writeByte(b byte) {
	this.w.WriteByte(b)
}

func (this *s_Writer) writeIdents(fold bool) {
	if fold {
		this.w.WriteString(this.idents[len(this.ident)-2:])
	} else {
		this.w.WriteString(this.idents)
	}
}

func (this *s_Writer) writeEndline() {
	this.w.WriteString(fsdef.Endline)
}

// ---------------------------------------------------------
// 结构体的类型
func (this *s_Writer) writeStructType(v reflect.Value) {
	t := v.Type()
	if t.PkgPath() == "" { // 匿名结构体
		this.w.WriteString("struct")
	} else {
		this.writeStringf("%v", t)
	}
}

// 结构体指针的类型
func (this *s_Writer) writePStructType(v reflect.Value) {
	t := v.Type().Elem()
	if t.PkgPath() == "" { // 匿名结构体
		if v.IsNil() {
			this.w.WriteString("(*struct)")
		} else {
			this.writeStringf("&struct<%#x>", v.Pointer())
		}
	} else {
		if v.IsNil() {
			this.writeStringf("(*%v)", t)
		} else {
			this.writeStringf("&%v<%#x>", t, v.Pointer())
		}
	}
}

// -----------------------------------------------
// isTop 是否是最顶层结构体
func (this *s_Writer) writeValue(v reflect.Value, tag string, isTop bool) {
	if v.IsValid() && v.IsZero() {
		v = reflect.Zero(v.Type())
	}

	if !isTop {
		if ok, str := this.getViewString(v); ok {
			this.writeStringf(str)
			return
		}
	}

	if v.Type().Kind() != reflect.Ptr {
		printer, ok := _printers[v.Type().Kind()]
		if ok {
			printer(this, v, tag)
		} else {
			_printOther(this, v, tag)
		}
		return
	}

	// 空指针
	if !v.IsValid() {
		this.writeStringf("%v", v)
		return
	}
	if v.IsNil() {
		if v.Type().Elem().Kind() == reflect.Struct {
			this.writePStructType(v)
			this.w.WriteString("(nil)")
		} else {
			this.writeStringf("%v", v)
		}
		return
	}

	elem := v.Elem()
	pprinter, ok := _pprinters[elem.Type().Kind()]
	if ok {
		pprinter(this, v, tag)
	} else {
		_printPOther(this, v, tag)
	}
}

// -------------------------------------------------------------------
// module private
// -------------------------------------------------------------------
// 数组/切片
func _printArray(w *s_Writer, v reflect.Value, tag string) {
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
			w.writeStringf("%v[...]", v.Type())
		} else {
			w.writeStringf("%v[...](len=%d)", v.Type(), v.Len())
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
	//if _isBaseType(e.Type()) {
	// tag 中有不换行标记
	if w.hasFieldTag(tag, "nowrap") {
		// 显示全部
		if fmtAll {
			w.writeStringf("%#v", v)
			return
		}
		// 只显示部分
		w.writeStringf("%v[", v.Type())
		for i := 0; i < fcount; i++ {
			w.writeStringf("%#v, ", v.Index(i))
		}
		w.writeStringf("...]")
		if !isArray {
			w.writeStringf("(len=%d)", ecount)
		}
		return
	}

	// 写入类型
	w.writeStringf("%v[", v.Type()) // 写入类型
	w.writeEndline()                // 换行

	// 写入第一个元素
	w.incLayer() // 增加嵌套数
	w.writeIdents(false)
	w.writeValue(e, tag, false)

	// 写入其他元素
	for i := 1; i < fcount; i++ {
		w.writeByte(',')
		w.writeEndline()
		w.writeIdents(false)
		w.writeValue(v.Index(i), tag, false)
	}

	// 如果没显示完，则打印省略号
	if !fmtAll {
		w.writeByte(',')
		w.writeEndline()
		w.writeIdents(false)
		w.writeStringf("...")
	}

	w.writeEndline()     // 换行
	w.decLayer()         // 减少嵌套
	w.writeIdents(false) // 缩进
	w.writeByte(']')     //数组结束

	// 如果是 slice 显示元素总数
	if !isArray && !fmtAll {
		w.writeStringf("(len=%d)", ecount)
	}
}

func _printPArray(w *s_Writer, v reflect.Value, tag string) {
	w.writeByte('&')
	_printArray(w, v, tag)
}

// -----------------------------------------------
// 映射
func _printMap(w *s_Writer, v reflect.Value, tag string) {
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
	//if _isBaseType(iter.Value().Type()) {
	// tag 中有不换行标记
	if w.hasFieldTag(tag, "nowrap") {
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
	w.incLayer()         // 增加嵌套数
	w.writeIdents(false) // 第一个元素的缩进
	w.writeStringf("%#v: ", iter.Key())
	w.writeValue(iter.Value(), tag, false)
	iter.Next()

	// 写入其他元素
	for fcount > 1 {
		w.writeByte(',')
		w.writeEndline()
		w.writeIdents(false)
		w.writeStringf("%#v: ", iter.Key())
		w.writeValue(iter.Value(), tag, false)
		fcount -= 1
		iter.Next()
	}

	// 如果没显示完，则打印省略号
	if !fmtAll {
		w.writeByte(',')
		w.writeEndline()
		w.writeIdents(false)
		w.writeStringf("...")
	}

	w.writeEndline()
	w.decLayer()
	w.writeIdents(false)
	w.writeByte('}')

	// 如果没显示完所有元素，则显示总数
	if !fmtAll {
		w.writeStringf("(len=%d)", ecount)
	}
}

func _printPMap(w *s_Writer, v reflect.Value, tag string) {
	w.writeByte('&')
	_printMap(w, v, tag)
}

// -----------------------------------------------
func _printStructField(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("{")
	w.incLayer()

	// 写成员
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		tf := t.Field(i)
		tags := tf.Tag.Get("fsfmt")
		if w.hasFieldTag(tags, "hide") { // 隐藏
			continue
		}
		fold := w.hasFieldTag(tags, "fold") // 折叠

		vf := v.Field(i)
		name := t.Field(i).Name
		w.writeEndline()
		w.writeIdents(fold)
		if fold {
			w.writeStringf("+ %s: {...}", name)
			continue
		}

		w.writeStringf("%s: ", name)
		w.enterPath(name)
		w.writeValue(vf, tags, false)
		w.leavePath()
		w.writeByte(',')
	}

	w.decLayer()
	w.writeEndline()
	w.writeIdents(false)
	w.writeByte('}')
}

// 结构体
func _printStruct(w *s_Writer, v reflect.Value, tag string) {
	if v.NumField() == 0 {
		w.writeStringf("%#v", v)
		return
	}

	w.writeStructType(v)
	_printStructField(w, v, tag)
}

// 结构体指针
func _printPStruct(w *s_Writer, v reflect.Value, tag string) {
	// 时父结构体的指针，则不再往下打印，否则死循环
	if w.isInStruct(v) {
		w.writeStringf("&%v<%#x>{parent...}", v.Type().Elem(), v.Pointer())
		return
	}
	w.enterStructPtr(v)
	defer w.leaveStructPtr()

	// nil 在调用该函数之前已经被排除，所以这里可以直接去 Elem
	elem := v.Elem()
	if elem.NumField() == 0 {
		w.writeStringf("&%#v", elem)
		return
	}

	w.writePStructType(v)
	_printStructField(w, elem, tag)
}

// ------------------------------------------------
// 有符号整型
func _printNumber(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("%v(%#v)", v.Type(), v)
}

// 有符号指针类型
func _printPNumber(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("&%v(%#v)", v.Type().Elem(), v)
}

// ------------------------------------------------
// 无符号整型
func _printUNumber(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("%v(%v=%#v)", v.Type(), v, v)
}

// 无符号指针类型
func _printPUNumber(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("&%v(%v=%#v)", v.Type().Elem(), v, v)
}

// ------------------------------------------------
// 复数
func _printComplex(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("%v%v", v.Type(), v)
}

// 复数指针
func _printPComplex(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("&%v%v", v.Type().Elem(), v)
}

// ------------------------------------------------
// 其他类型
func _printOther(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("%#v", v)
}

// 其他类型的指针类型
func _printPOther(w *s_Writer, v reflect.Value, tag string) {
	w.writeStringf("&%v(%#v)", v.Type().Elem(), v)
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

var _printers map[reflect.Kind]func(*s_Writer, reflect.Value, string)  // 类型打印方法
var _pprinters map[reflect.Kind]func(*s_Writer, reflect.Value, string) // 指针类型方法
var _baseTypes map[reflect.Kind]interface{}                            // 基础类型，这些类型如果为 array、slice、map 的成员，则分列成员时，不换行

func init() {
	_printers = map[reflect.Kind]func(*s_Writer, reflect.Value, string){
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
	_pprinters = map[reflect.Kind]func(*s_Writer, reflect.Value, string){
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
	Idents    int            // 缩进空格数
	FmtCounts map[string]int // 显示指定数量的元素，格式为：map["aa.bb.cc"] = count（只对 array、slice、map 有效）
}

// SprintStruct 以初始化结构的格式，将一个结构体格式化为字符串，并写入流中
// 参数：
//
//	 w: 流缓冲
//		st: 要格式化的结构体
//		prefix: 整个输出结构体的每一行的前缀
//		idents: 缩进的空格数
//		如果要格式化输出的结构体成员具有以下 tag，则该成员不会被格式化输出：
//		struct{
//			m int     `fsfmt:"hide"`	// 不会打印
//			mm string `fsfmt:"fold"`    // 折叠输出
//		}
func FPrintStruct(w io.Writer, obj interface{}, opts *S_FmtOpts) {
	if opts == nil {
		opts = &S_FmtOpts{"    ", 4, nil}
	}
	tobj := reflect.TypeOf(obj)
	if tobj == nil {
		w.Write([]byte("nil"))
		return
	}

	writer := _newWriter(w, opts)
	defer writer.flush()

	vobj := reflect.ValueOf(obj)
	if tobj.Kind() == reflect.Ptr {
		if !vobj.IsValid() {
			writer.w.Write([]byte(fmt.Sprintf("(%v)(nil)", tobj)))
			return
		}
		if vobj.IsNil() {
			if tobj.Elem().Kind() == reflect.Struct {
				writer.writePStructType(vobj)
				writer.w.WriteString("(nil)")
			} else {
				writer.w.Write([]byte(fmt.Sprintf("(%v)(nil)", tobj)))
			}
			return
		}
		if tobj.Elem().Kind() != reflect.Struct {
			opts.Prefix = ""
		}
	} else if tobj.Kind() != reflect.Struct {
		opts.Prefix = ""
	}

	writer.writeStringf(opts.Prefix)
	writer.writeValue(vobj, "", true)
}

// SprintStruct 以初始化结构的格式，将一个结构体格式化为字符串
// 参数：
//
//	st: 要格式化的结构体
//	prefix: 整个输出结构体的每一行的前缀
//	idents: 缩进字符串
//	如果要格式化输出的结构体成员具有以下 tag，则该成员不会被格式化输出：
//	struct{
//		m int     `fsfmt:"hide"`    // 打印不输出
//		mm string `fsfmt:"fold"`    // 折叠
//	}
func SprintStruct(obj interface{}, opts *S_FmtOpts) string {
	out := bytes.NewBuffer([]byte{})
	FPrintStruct(out, obj, opts)
	return out.String()
}
