/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: logger for every day
@author: fanky
@version: 1.0
@date: 2022-06-25
**/

// logs are writen to files and generate one log file every day

package fslog

import (
	"fmt"
	"os"
	"time"
)

// -----------------------------------------------------------------------------
// FileLogger
// -----------------------------------------------------------------------------
type S_FileLogger struct {
	*S_Logger
	file     *os.File
	lastTime time.Time
}

// NewDayfileLogger，新建 DayfileLogger
// root 为 log 的根目录
// filePrefix 为 log 文件名前缀
func NewFileLogger(file string) (*S_FileLogger, error) {
	logger := &S_FileLogger{}
	logger.S_Logger = NewLogger(logger.write)
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return nil, fmt.Errorf("create log file fail, %v", err)
	}
	logger.file = f
	return logger, nil
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
// 父类中的 send 函数中已经 lock，因此这里不需要再上锁了
func (this *S_FileLogger) write(t time.Time, lv string, msg []byte) {
	if this.file == nil {
		now := this.nowTime()
		os.Stderr.WriteString(fmt.Sprintf("%s: %v", now.Format("[ERROR]|2006/01/02 15:04:05.999999 "), "file logger has been closed"))
		os.Stdout.Write(msg)
		return
	}
	if _, err := this.file.Write(msg); err != nil {
		now := this.nowTime()
		os.Stderr.WriteString(fmt.Sprintf("%s: %v", now.Format("[ERROR]|2006/01/02 15:04:05.999999 "), err))
		os.Stdout.Write(msg)
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_FileLogger) Close() {
	this.Lock()
	defer this.Unlock()
	if this.file != nil {
		this.file.Close()
		this.file = nil
	}
}

// ---------------------------------------------------------
func (this *S_FileLogger) Debug(arg any, args ...any) {
	this.S_Logger.Debug_(1, arg, args...)
}

func (this *S_FileLogger) Debugf(msg string, args ...any) {
	this.S_Logger.Debugf_(1, msg, args...)
}

func (this *S_FileLogger) Info(arg any, args ...any) {
	this.S_Logger.Info_(1, arg, args...)
}

func (this *S_FileLogger) Infof(msg string, args ...any) {
	this.S_Logger.Infof_(1, msg, args...)
}

func (this *S_FileLogger) Notic(arg any, args ...any) {
	this.S_Logger.Notic_(1, arg, args...)
}

func (this *S_FileLogger) Noticf(msg string, args ...any) {
	this.S_Logger.Noticf_(1, msg, args...)
}

func (this *S_FileLogger) Warn(arg any, args ...any) {
	this.S_Logger.Warn_(1, arg, args...)
}

func (this *S_FileLogger) Warnf(msg string, args ...any) {
	this.S_Logger.Warnf_(1, msg, args...)
}

func (this *S_FileLogger) Error(arg any, args ...any) {
	this.S_Logger.Error_(1, arg, args...)
}

func (this *S_FileLogger) Errorf(msg string, args ...any) {
	this.S_Logger.Errorf_(1, msg, args...)
}

func (this *S_FileLogger) Hack(arg any, args ...any) {
	this.S_Logger.Hack_(1, arg, args...)
}

func (this *S_FileLogger) Hackf(msg string, args ...any) {
	this.S_Logger.Hackf_(1, msg, args...)
}

func (this *S_FileLogger) Critical(arg any, args ...any) {
	this.S_Logger.Critical_(1, arg, args...)
}

func (this *S_FileLogger) Criticalf(msg string, args ...any) {
	this.S_Logger.Criticalf_(1, msg, args...)
}

func (this *S_FileLogger) Trace(arg any, args ...any) {
	this.S_Logger.Trace_(1, arg, args...)
}

func (this *S_FileLogger) Tracef(msg string, args ...any) {
	this.S_Logger.Tracef_(1, msg, args...)
}

func (this *S_FileLogger) Panic(arg any, args ...any) {
	this.S_Logger.Panic_(1, arg, args...)
}

func (this *S_FileLogger) Panicf(msg string, args ...any) {
	this.S_Logger.Panicf_(1, msg, args...)
}

func (this *S_FileLogger) Fatal(arg any, args ...any) {
	this.S_Logger.Fatal_(1, arg, args...)
}

func (this *S_FileLogger) Fatalf(msg string, args ...any) {
	this.S_Logger.Fatalf_(1, msg, args...)
}
