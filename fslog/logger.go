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
)

type S_Logger struct {
	sync.Mutex
	writer   func(time.Time, []byte)
	levels   map[string]bool
	levelFmt string
	goidSize int

	showSite   bool
	UseUTCTime bool
}

func newLogger(writer func(time.Time, []byte)) *S_Logger {
	logger := &S_Logger{
		writer: writer,
		levels: map[string]bool{
			"DEBUG": true,
			"INFO":  true,
			"NOTIC": true,
			"WARN":  true,
			"ERROR": true,
			"HACK":  true,
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

func (this *S_Logger) send(t time.Time, msg []byte) {
	this.Lock()
	defer this.Unlock()
	if this.writer != nil {
		this.writer(t, msg)
	}
}

// ---------------------------------------------------------
func (this *S_Logger) writef(bb *bytes.Buffer, msg string, args ...any) {
	bb.WriteString(fmt.Sprintf(msg, args...))
}

func (this *S_Logger) writeGoID(bb *bytes.Buffer) {
	var buf = make([]byte, 23)
	runtime.Stack(buf, false)
	fields := bytes.Fields(buf)
	if len(fields) > 1 {
		goid := fields[1]
		if len(goid) > this.goidSize {
			this.goidSize = len(goid)
		}
		this.writef(bb, fmt.Sprintf("[G-%%0%ds]|", this.goidSize), goid)
	} else {
		bb.WriteString("[G-ERR]|")
	}
}

func (this *S_Logger) writePrefix(bb *bytes.Buffer, level string) {
	this.writef(bb, this.levelFmt, "["+level+"]")
}

func (this *S_Logger) writeDateTime(bb *bytes.Buffer) time.Time {
	now := this.nowTime()
	bb.WriteString(now.Format("2006/01/02 15:04:05.999999"))
	return now
}

func (this *S_Logger) writeCallStack(bb *bytes.Buffer, depth int) {
	pc := make([]uintptr, 50)
	n := runtime.Callers(depth+2, pc)
	if n < 1 {
		bb.WriteString("\n???:?")
		return
	}
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		this.writef(bb, "\n\t%s(...):\n", frame.Function)
		this.writef(bb, "\t\t%s:%d", frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func (this *S_Logger) writeCallTopStack(bb *bytes.Buffer, depth int) {
	_, file, line, ok := runtime.Caller(depth + 1)
	if !ok {
		bb.WriteString("???:?")
		return
	}
	this.writef(bb, "%s:%d", file, line)
}

// ---------------------------------------------------------
func (this *S_Logger) isShield(lv string) bool {
	this.Lock()
	defer this.Unlock()
	return !this.levels[lv]
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_Logger) ToggleSite(showSite bool) {
	this.showSite = showSite
}

func (this *S_Logger) SetOutputWriter(writer func(time.Time, []byte)) {
	this.Lock()
	defer this.Unlock()
	this.writer = writer
}

// ---------------------------------------------------------
func (this *S_Logger) Output(depth int, level string, arg any, args ...any) {
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
	this.send(now, bb.Bytes())
}

func (this *S_Logger) Outputf(depth int, level string, msg string, args ...any) {
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
	this.send(now, bb.Bytes())
}

// ---------------------------------------------------------
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
		this.send(now, bb.Bytes())
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
		this.send(now, bb.Bytes())
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
	this.send(now, bb.Bytes())
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
	this.send(now, bb.Bytes())
}

// ---------------------------------------------------------
// 过滤输出级别
// debug/info/warn/error/notice/clit/trace
func (this *S_Logger) Shield(lvs ...string) {
	this.Lock()
	defer this.Unlock()
	for _, lv := range lvs {
		lv = strings.ToUpper(lv)
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
	for _, lv := range lvs {
		lv = strings.ToUpper(lv)
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
