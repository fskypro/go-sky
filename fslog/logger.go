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

	UseUTCTime bool
}

func newLogger(writer func(time.Time, []byte)) *S_Logger {
	logger := &S_Logger{
		writer: writer,
		levels: map[string]bool{
			"DEBUG": true,
			"INFO":  true,
			"WARN":  true,
			"ERROR": true,
			"HACK":  true,
			"CRIT":  true,
			"TRACE": true,
			"PANIC": true,
			"FATAL": true,
		},
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
func (this *S_Logger) send(t time.Time, msg []byte) {
	this.Lock()
	defer this.Unlock()
	if this.writer != nil {
		this.writer(t, msg)
	}
}

// ---------------------------------------------------------
func (this *S_Logger) writef(bb *bytes.Buffer, msg string, args ...any) {
	bb.Write([]byte(fmt.Sprintf(msg, args...)))
}

func (this *S_Logger) writeGoID(bb *bytes.Buffer) {
	var buf = make([]byte, 23)
	runtime.Stack(buf, false)
	fields := bytes.Fields(buf)
	if len(fields) > 1 {
		this.writef(bb, "[G-%04s]|", fields[1])
	} else {
		bb.WriteString("[G-ERR]|")
	}
}

func (this *S_Logger) writePrefix(bb *bytes.Buffer, level string) {
	this.writef(bb, this.levelFmt, "["+level+"]")
}

func (this *S_Logger) writeDateTime(bb *bytes.Buffer) time.Time {
	now := time.Now()
	if this.UseUTCTime {
		now = now.UTC()
	}
	bb.WriteString(time.Now().Format("2006-01-02 15:04:05.999999"))
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
	bb.WriteByte(' ')
	this.writeCallTopStack(bb, depth+1)
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
	bb.WriteByte(' ')
	this.writeCallTopStack(bb, depth+1)
	bb.WriteString(": ")
	this.writef(bb, msg, args...)
	bb.WriteByte('\n')
	this.send(now, bb.Bytes())
}

// ---------------------------------------------------------
func (this *S_Logger) Debug(depth int, arg any, args ...any) {
	this.Output(depth+1, "DEBUG", arg, args...)
}

func (this *S_Logger) Debugf(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "DEBUG", msg, args...)
}

func (this *S_Logger) Info(depth int, arg any, args ...any) {
	this.Output(depth+1, "INFO", arg, args...)
}

func (this *S_Logger) Infof(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "INFO", msg, args...)
}

func (this *S_Logger) Warn(depth int, arg any, args ...any) {
	this.Output(depth+1, "WARN", arg, args...)
}

func (this *S_Logger) Warnf(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "WARN", msg, args...)
}

func (this *S_Logger) Error(depth int, arg any, args ...any) {
	this.Output(depth+1, "ERROR", arg, args...)
}

func (this *S_Logger) Errorf(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "ERROR", msg, args...)
}

func (this *S_Logger) Hack(depth int, arg any, args ...any) {
	this.Output(depth+1, "HACK", arg, args...)
}

func (this *S_Logger) Hackf(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "HACK", msg, args...)
}

func (this *S_Logger) Critical(depth int, arg any, args ...any) {
	this.Output(depth+1, "CRIT", arg, args...)
}

func (this *S_Logger) Criticalf(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "CRIT", msg, args...)
}

func (this *S_Logger) Fatal(depth int, arg any, args ...any) {
	this.Output(depth+1, "FATAL", arg, args...)
	os.Exit(1)
}

func (this *S_Logger) Fatalf(depth int, msg string, args ...any) {
	this.Outputf(depth+1, "FATAL", msg, args...)
	os.Exit(1)
}

func (this *S_Logger) Panic(depth int, arg any, args ...any) {
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

func (this *S_Logger) Panicf(depth int, msg string, args ...any) {
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

func (this *S_Logger) Trace(depth int, arg any, args ...any) {
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

func (this *S_Logger) Tracef(depth int, msg string, args ...any) {
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
func (this *S_Logger) Shield(lv string, lvs ...string) {
	this.Lock()
	defer this.Unlock()
	for _, lv := range append([]string{lv}, lvs...) {
		lv = strings.ToUpper(lv)
		if _, ok := this.levels[lv]; ok {
			this.levels[lv] = false
		}
	}
}

// 取消过滤日志级别
// debug/info/warn/error/notice/clit/trace
func (this *S_Logger) Unshield(lv string, lvs ...string) {
	this.Lock()
	defer this.Unlock()
	for _, lv := range append([]string{lv}, lvs...) {
		lv = strings.ToUpper(lv)
		if _, ok := this.levels[lv]; ok {
			this.levels[lv] = true
		}
	}
}
