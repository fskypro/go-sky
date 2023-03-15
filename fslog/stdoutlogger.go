/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: logger for stdout
@author: fanky
@version: 1.0
@date: 2022-06-25
**/

// logs are writen to files and generate one log file every day

package fslog

import (
	"context"
	"os"
	"time"
)

type S_StdoutLogger struct {
	*S_Logger
	cancel func()
}

// NewDayfileLogger，新建 DayfileLogger
// root 为 log 的根目录
// filePrefix 为 log 文件名前缀
func NewStdoutLogger() *S_StdoutLogger {
	ctx, cancel := context.WithCancel(context.Background())
	logger := NewStdoutLoggerContex(ctx)
	logger.cancel = cancel
	return logger
}

func NewStdoutLoggerContex(ctx context.Context) *S_StdoutLogger {
	logger := &S_StdoutLogger{}
	logger.S_Logger = newLogger(logger.write)
	return logger
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func (this *S_StdoutLogger) write(t time.Time, msg []byte) {
	os.Stdout.Write(msg)
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_StdoutLogger) Debug(arg any, args ...any) {
	this.S_Logger.Debug_(1, arg, args...)
}

func (this *S_StdoutLogger) Debugf(msg string, args ...any) {
	this.S_Logger.Debugf_(1, msg, args...)
}

func (this *S_StdoutLogger) Info(arg any, args ...any) {
	this.S_Logger.Info_(1, arg, args...)
}

func (this *S_StdoutLogger) Infof(msg string, args ...any) {
	this.S_Logger.Infof_(1, msg, args...)
}

func (this *S_StdoutLogger) Notic(arg any, args ...any) {
	this.S_Logger.Notic_(1, arg, args...)
}

func (this *S_StdoutLogger) Noticf(msg string, args ...any) {
	this.S_Logger.Noticf_(1, msg, args...)
}

func (this *S_StdoutLogger) Warn(arg any, args ...any) {
	this.S_Logger.Warn_(1, arg, args...)
}

func (this *S_StdoutLogger) Warnf(msg string, args ...any) {
	this.S_Logger.Warnf_(1, msg, args...)
}

func (this *S_StdoutLogger) Error(arg any, args ...any) {
	this.S_Logger.Error_(1, arg, args...)
}

func (this *S_StdoutLogger) Errorf(msg string, args ...any) {
	this.S_Logger.Errorf_(1, msg, args...)
}

func (this *S_StdoutLogger) Hack(arg any, args ...any) {
	this.S_Logger.Hack_(1, arg, args...)
}

func (this *S_StdoutLogger) Hackf(msg string, args ...any) {
	this.S_Logger.Hackf_(1, msg, args...)
}

func (this *S_StdoutLogger) Critical(arg any, args ...any) {
	this.S_Logger.Critical_(1, arg, args...)
}

func (this *S_StdoutLogger) Criticalf(msg string, args ...any) {
	this.S_Logger.Criticalf_(1, msg, args...)
}

func (this *S_StdoutLogger) Trace(arg any, args ...any) {
	this.S_Logger.Trace_(1, arg, args...)
}

func (this *S_StdoutLogger) Tracef(msg string, args ...any) {
	this.S_Logger.Tracef_(1, msg, args...)
}

func (this *S_StdoutLogger) Panic(arg any, args ...any) {
	this.S_Logger.Panic_(1, arg, args...)
}

func (this *S_StdoutLogger) Panicf(msg string, args ...any) {
	this.S_Logger.Panicf_(1, msg, args...)
}

func (this *S_StdoutLogger) Fatal(arg any, args ...any) {
	this.S_Logger.Fatal_(1, arg, args...)
}

func (this *S_StdoutLogger) Fatalf(msg string, args ...any) {
	this.S_Logger.Fatalf_(1, msg, args...)
}
