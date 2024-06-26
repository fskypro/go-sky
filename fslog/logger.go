/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: base logger
@author: fanky
@version: 1.0
@date: 2022-06-25
**/

package fslog

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"fsky.pro/fspath"
)

type S_Logger struct {
	sync.Mutex
	writer   func(time.Time, T_Level, []byte)
	levels   map[T_Level]bool
	levelFmt string
	goidSize int

	showSite   bool
	cutSrcRoot string
	UseUTCTime bool
}

func NewLogger(writer func(time.Time, T_Level, []byte)) *S_Logger {
	logger := &S_Logger{
		writer: writer,
		levels: map[T_Level]bool{
			"DEBUG": true,
			"INFO":  true,
			"NOTIC": true,
			"WARN":  true,
			"ERROR": true,
			"HACK":  true,
			"ILLEG": true,
			"CRIT":  true,
			"TRACE": true,
			"PANIC": true,
			"FATAL": true,
		},
		showSite: true,
	}
	lvLen := 1
	for lv := range logger.levels {
		if lvLen < len(lv) {
			lvLen = len(lv)
		}
	}
	logger.levelFmt = fmt.Sprintf("%%-%ds|", lvLen+2)
	return logger
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func (this *S_Logger) nowTime() time.Time {
	if this.UseUTCTime {
		return time.Now().UTC()
	}
	return time.Now()
}

func (this *S_Logger) send(t time.Time, lv T_Level, msg []byte) {
	this.Lock()
	defer this.Unlock()
	if this.writer != nil {
		this.writer(t, lv, msg)
	}
}

// ---------------------------------------------------------
func (this *S_Logger) writef(buff *bytes.Buffer, msg string, args ...any) {
	buff.WriteString(fmt.Sprintf(msg, args...))
}

func (this *S_Logger) writeGoID(buff *bytes.Buffer) {
	var buf = make([]byte, 23)
	runtime.Stack(buf, false)
	fields := bytes.Fields(buf)
	if len(fields) > 1 {
		goid := fields[1]
		if len(goid) > this.goidSize {
			this.goidSize = len(goid)
		}
		this.writef(buff, fmt.Sprintf("[G-%%0%ds]|", this.goidSize), goid)
	} else {
		buff.WriteString("[G-ERR]|")
	}
}

func (this *S_Logger) writePrefix(buff *bytes.Buffer, level T_Level) {
	this.writef(buff, this.levelFmt, "["+level+"]")
}

func (this *S_Logger) writeDateTime(buff *bytes.Buffer) time.Time {
	now := this.nowTime()
	buff.WriteString(now.Format("2006/01/02 15:04:05.999999"))
	return now
}

func (this *S_Logger) writeCallStack(buff *bytes.Buffer, depth int) {
	pc := make([]uintptr, 50)
	n := runtime.Callers(depth+2, pc)
	if n < 1 {
		buff.WriteString("\n???:?")
		return
	}
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		this.writef(buff, "\n\t%s(...):\n", frame.Function)
		file := fspath.CleanPath(frame.File)
		if this.cutSrcRoot != "" {
			file = strings.TrimPrefix(file, this.cutSrcRoot)
			file = strings.TrimSuffix(file, ".go")
		}
		this.writef(buff, "\t\t%s:%d", file, frame.Line)
		if !more {
			break
		}
	}
}

func (this *S_Logger) writeCallTopStack(buff *bytes.Buffer, depth int) {
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		buff.WriteString("???:?")
		return
	}
	file = fspath.CleanPath(file)
	if this.cutSrcRoot != "" {
		file = strings.TrimPrefix(file, this.cutSrcRoot)
		file = strings.TrimSuffix(file, ".go")
	}
	this.writef(buff, "%s:%d", file, line)
}

// ---------------------------------------------------------
func (this *S_Logger) isShield(lv T_Level) bool {
	this.Lock()
	defer this.Unlock()
	return !this.levels[lv]
}

// -------------------------------------------------------------------
// protected
// -------------------------------------------------------------------
func (this *S_Logger) Debug_(depth int, arg any, args ...any) {
	this.Output(depth+1, "DEBUG", arg, args...)
}

func (this *S_Logger) Debugf_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "DEBUG", msg, args...)
}

func (this *S_Logger) Info_(depth int, arg any, args ...any) {
	this.Output(depth+1, "INFO", arg, args...)
}

func (this *S_Logger) Infof_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "INFO", msg, args...)
}

func (this *S_Logger) Notic_(depth int, arg any, args ...any) {
	this.Output(depth+1, "NOTIC", arg, args...)
}

func (this *S_Logger) Noticf_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "NOTIC", msg, args...)
}

func (this *S_Logger) Warn_(depth int, arg any, args ...any) {
	this.Output(depth+1, "WARN", arg, args...)
}

func (this *S_Logger) Warnf_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "WARN", msg, args...)
}

func (this *S_Logger) Error_(depth int, arg any, args ...any) {
	this.Output(depth+1, "ERROR", arg, args...)
}

func (this *S_Logger) Errorf_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "ERROR", msg, args...)
}

func (this *S_Logger) Hack_(depth int, arg any, args ...any) {
	this.Output(depth+1, "HACK", arg, args...)
}

func (this *S_Logger) Hackf_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "HACK", msg, args...)
}

func (this *S_Logger) Illeg_(depth int, arg any, args ...any) {
	this.Output(depth+1, "ILLEG", arg, args...)
}

func (this *S_Logger) Illegf_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "ILLEG", msg, args...)
}

func (this *S_Logger) Critical_(depth int, arg any, args ...any) {
	this.Output(depth+1, "CRIT", arg, args...)
}

func (this *S_Logger) Criticalf_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "CRIT", msg, args...)
}

func (this *S_Logger) Fatal_(depth int, arg any, args ...any) {
	this.Output(depth+1, "FATAL", arg, args...)
	os.Exit(1)
}

func (this *S_Logger) Fatalf_(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "FATAL", msg, args...)
	os.Exit(1)
}

func (this *S_Logger) Panic_(depth int, arg any, args ...any) {
	bb := new(bytes.Buffer)
	this.writeGoID(bb)
	this.writePrefix(bb, "PANIC")
	now := this.writeDateTime(bb)
	bb.WriteByte(' ')
	this.writeCallTopStack(bb, depth+1)
	bb.WriteString(": ")
	this.writef(bb, "%v", arg)
	for _, a := range args {
		this.writef(bb, " %v", a)
	}

	msg := bb.String()
	if !this.isShield("PANIC") {
		this.writeCallStack(bb, depth+1)
		bb.WriteByte('\n')
		this.send(now, "PANIC", bb.Bytes())
	}
	panic(msg)
}

func (this *S_Logger) Panicf_(depth int, msg string, args ...any) {
	bb := new(bytes.Buffer)
	this.writeGoID(bb)
	this.writePrefix(bb, "PANIC")
	now := this.writeDateTime(bb)
	bb.WriteByte(' ')
	this.writeCallTopStack(bb, depth+1)
	bb.WriteString(": ")
	this.writef(bb, msg, args...)

	msg = bb.String()
	if !this.isShield("PANIC") {
		this.writeCallStack(bb, depth+1)
		bb.WriteByte('\n')
		this.send(now, "PANIC", bb.Bytes())
	}
	panic(msg)
}

func (this *S_Logger) Trace_(depth int, arg any, args ...any) {
	if this.isShield("TRACE") {
		return
	}
	bb := new(bytes.Buffer)
	this.writeGoID(bb)
	this.writePrefix(bb, "TRACE")
	now := this.writeDateTime(bb)
	bb.WriteByte(' ')
	this.writeCallTopStack(bb, depth+1)
	bb.WriteString(": ")
	this.writef(bb, "%v", arg)
	for _, a := range args {
		this.writef(bb, " %v", a)
	}
	this.writeCallStack(bb, depth+1)
	bb.WriteByte('\n')
	this.send(now, "TRACE", bb.Bytes())
}

func (this *S_Logger) Tracef_(depth int, msg string, args ...any) {
	if this.isShield("TRACE") {
		return
	}
	bb := new(bytes.Buffer)
	this.writeGoID(bb)
	this.writePrefix(bb, "TRACE")
	now := this.writeDateTime(bb)
	bb.WriteByte(' ')
	this.writeCallTopStack(bb, depth+1)
	bb.WriteString(": ")
	this.writef(bb, msg, args...)
	this.writeCallStack(bb, depth+1)
	bb.WriteByte('\n')
	this.send(now, "TRACE", bb.Bytes())
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_Logger) ToggleSite(showSite bool) {
	this.showSite = showSite
}

func (this *S_Logger) CutSrcRoot(root string) {
	this.cutSrcRoot = fspath.CleanPath(root) + string(os.PathSeparator)
}

func (this *S_Logger) SetOutputWriter(writer func(time.Time, T_Level, []byte)) {
	this.Lock()
	defer this.Unlock()
	this.writer = writer
}

// ---------------------------------------------------------
func (this *S_Logger) Output(depth int, level T_Level, arg any, args ...any) {
	if this.isShield(level) {
		return
	}
	bb := new(bytes.Buffer)
	this.writeGoID(bb)
	this.writePrefix(bb, level)
	now := this.writeDateTime(bb)
	if this.showSite {
		bb.WriteByte(' ')
		this.writeCallTopStack(bb, depth+1)
	}
	bb.WriteString(": ")
	this.writef(bb, "%v", arg)
	for _, a := range args {
		this.writef(bb, " %v", a)
	}
	bb.WriteByte('\n')
	this.send(now, level, bb.Bytes())
}

func (this *S_Logger) Outputf(depth int, level T_Level, msg string, args ...any) {
	if this.isShield(level) {
		return
	}
	bb := new(bytes.Buffer)
	this.writeGoID(bb)
	this.writePrefix(bb, level)
	now := this.writeDateTime(bb)
	if this.showSite {
		bb.WriteByte(' ')
		this.writeCallTopStack(bb, depth+1)
	}
	bb.WriteString(": ")
	this.writef(bb, msg, args...)
	bb.WriteByte('\n')
	this.send(now, level, bb.Bytes())
}

// ---------------------------------------------------------
// 过滤输出级别
// debug/info/warn/error/notice/clit/trace
func (this *S_Logger) Shield(lvs ...string) {
	this.Lock()
	defer this.Unlock()
	for _, l := range lvs {
		lv := T_Level(strings.ToUpper(l))
		if _, ok := this.levels[lv]; ok {
			this.levels[lv] = false
		}
	}
}

// 屏蔽所有输出级别
func (this *S_Logger) ShieldAll() {
	this.Lock()
	defer this.Unlock()
	for lv := range this.levels {
		this.levels[lv] = false
	}
}

// 取消过滤日志级别
// debug/info/warn/error/notice/clit/trace
func (this *S_Logger) Unshield(lvs ...string) {
	this.Lock()
	defer this.Unlock()
	for _, l := range lvs {
		lv := T_Level(strings.ToUpper(l))
		if _, ok := this.levels[lv]; ok {
			this.levels[lv] = true
		}
	}
}

// 解除所有屏蔽输出级别
func (this *S_Logger) UnshieldAll() {
	this.Lock()
	defer this.Unlock()
	for lv := range this.levels {
		this.levels[lv] = true
	}
}

// ---------------------------------------------------------
func (this *S_Logger) Whatever(lv T_Level, arg any, args ...any) {
	this.Output(1, lv, arg, args...)
}

func (this *S_Logger) Whateverf(lv T_Level, msg string, args ...any) {
	this.Output(1, lv, msg, args)
}

// 直接输出字符串，不带任何前缀记录
func (this *S_Logger) Direct(lv T_Level, msg string) {
	if this.isShield(lv) { return }
	bb := new(bytes.Buffer)
	this.writeGoID(bb)
	this.writePrefix(bb, lv)
	bb.Write([]byte(msg))
	bb.WriteByte('\n')
	this.send(time.Now(), lv, bb.Bytes())
}

func (this *S_Logger) Directf(lv T_Level, msg string, args ...any) {
	this.Direct(lv, fmt.Sprintf(msg, args...))
}
