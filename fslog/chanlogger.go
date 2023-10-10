/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: logger write to channel
@author: fanky
@version: 1.0
@date: 2023-05-01
**/

// logs are writen to files and generate one log file every day

package fslog

import (
	"time"
)

type S_ChanLog struct {
	Level string
	Msg   []byte
}

type S_ChanLogger struct {
	*S_Logger
	chMsg chan *S_ChanLog
}

// NewDayfileLogger，新建 DayfileLogger
// root 为 log 的根目录
// filePrefix 为 log 文件名前缀
func NewChanLogger(ch chan *S_ChanLog) *S_ChanLogger {
	logger := &S_ChanLogger{chMsg: ch}
	return logger
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func (this *S_ChanLogger) write(t time.Time, level string, msg []byte) {
	this.chMsg <- &S_ChanLog{level, msg}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_ChanLogger) Debug(arg any, args ...any) {
	this.S_Logger.Debug_(1, arg, args...)
}

func (this *S_ChanLogger) Debugf(msg string, args ...any) {
	this.S_Logger.Debugf_(1, msg, args...)
}

func (this *S_ChanLogger) Info(arg any, args ...any) {
	this.S_Logger.Info_(1, arg, args...)
}

func (this *S_ChanLogger) Infof(msg string, args ...any) {
	this.S_Logger.Infof_(1, msg, args...)
}

func (this *S_ChanLogger) Notic(arg any, args ...any) {
	this.S_Logger.Notic_(1, arg, args...)
}

func (this *S_ChanLogger) Noticf(msg string, args ...any) {
	this.S_Logger.Noticf_(1, msg, args...)
}

func (this *S_ChanLogger) Warn(arg any, args ...any) {
	this.S_Logger.Warn_(1, arg, args...)
}

func (this *S_ChanLogger) Warnf(msg string, args ...any) {
	this.S_Logger.Warnf_(1, msg, args...)
}

func (this *S_ChanLogger) Error(arg any, args ...any) {
	this.S_Logger.Error_(1, arg, args...)
}

func (this *S_ChanLogger) Errorf(msg string, args ...any) {
	this.S_Logger.Errorf_(1, msg, args...)
}

func (this *S_ChanLogger) Hack(arg any, args ...any) {
	this.S_Logger.Hack_(1, arg, args...)
}

func (this *S_ChanLogger) Hackf(msg string, args ...any) {
	this.S_Logger.Hackf_(1, msg, args...)
}

func (this *S_ChanLogger) Critical(arg any, args ...any) {
	this.S_Logger.Critical_(1, arg, args...)
}

func (this *S_ChanLogger) Criticalf(msg string, args ...any) {
	this.S_Logger.Criticalf_(1, msg, args...)
}

func (this *S_ChanLogger) Trace(arg any, args ...any) {
	this.S_Logger.Trace_(1, arg, args...)
}

func (this *S_ChanLogger) Tracef(msg string, args ...any) {
	this.S_Logger.Tracef_(1, msg, args...)
}

func (this *S_ChanLogger) Panic(arg any, args ...any) {
	this.S_Logger.Panic_(1, arg, args...)
}

func (this *S_ChanLogger) Panicf(msg string, args ...any) {
	this.S_Logger.Panicf_(1, msg, args...)
}

func (this *S_ChanLogger) Fatal(arg any, args ...any) {
	this.S_Logger.Fatal_(1, arg, args...)
}

func (this *S_ChanLogger) Fatalf(msg string, args ...any) {
	this.S_Logger.Fatalf_(1, msg, args...)
}
