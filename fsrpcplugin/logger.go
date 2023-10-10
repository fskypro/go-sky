/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: logger for stdout
@author: fanky
@version: 1.0
@date: 2022-06-25
**/

// logs are writen to files and generate one log file every day

package fsrpcplugin

import (
	"os"
	"time"

	"fsky.pro/fslog"
)

type S_Logger struct {
	*fslog.S_Logger
}

func NewLogger() *S_Logger {
	logger := &S_Logger{}
	logger.S_Logger = fslog.NewLogger(logger.write)
	return logger
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func (this *S_Logger) write(t time.Time, lv string, msg []byte) {
	prefix := []byte(lv + ">")
	//msg = bytes.ReplaceAll(msg, []byte{'\n'}, []byte{1})
	//msg = append(prefix, msg...)
	os.Stdout.Write(append(prefix, msg...))
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_Logger) Debug(arg any, args ...any) {
	this.S_Logger.Debug_(1, arg, args...)
}

func (this *S_Logger) Debugf(msg string, args ...any) {
	this.S_Logger.Debugf_(1, msg, args...)
}

func (this *S_Logger) Info(arg any, args ...any) {
	this.S_Logger.Info_(1, arg, args...)
}

func (this *S_Logger) Infof(msg string, args ...any) {
	this.S_Logger.Infof_(1, msg, args...)
}

func (this *S_Logger) Notic(arg any, args ...any) {
	this.S_Logger.Notic_(1, arg, args...)
}

func (this *S_Logger) Noticf(msg string, args ...any) {
	this.S_Logger.Noticf_(1, msg, args...)
}

func (this *S_Logger) Warn(arg any, args ...any) {
	this.S_Logger.Warn_(1, arg, args...)
}

func (this *S_Logger) Warnf(msg string, args ...any) {
	this.S_Logger.Warnf_(1, msg, args...)
}

func (this *S_Logger) Error(arg any, args ...any) {
	this.S_Logger.Error_(1, arg, args...)
}

func (this *S_Logger) Errorf(msg string, args ...any) {
	this.S_Logger.Errorf_(1, msg, args...)
}

func (this *S_Logger) Hack(arg any, args ...any) {
	this.S_Logger.Hack_(1, arg, args...)
}

func (this *S_Logger) Hackf(msg string, args ...any) {
	this.S_Logger.Hackf_(1, msg, args...)
}

func (this *S_Logger) Critical(arg any, args ...any) {
	this.S_Logger.Critical_(1, arg, args...)
}

func (this *S_Logger) Criticalf(msg string, args ...any) {
	this.S_Logger.Criticalf_(1, msg, args...)
}

func (this *S_Logger) Trace(arg any, args ...any) {
	this.S_Logger.Trace_(1, arg, args...)
}

func (this *S_Logger) Tracef(msg string, args ...any) {
	this.S_Logger.Tracef_(1, msg, args...)
}

func (this *S_Logger) Panic(arg any, args ...any) {
	this.S_Logger.Panic_(1, arg, args...)
}

func (this *S_Logger) Panicf(msg string, args ...any) {
	this.S_Logger.Panicf_(1, msg, args...)
}

func (this *S_Logger) Fatal(arg any, args ...any) {
	this.S_Logger.Fatal_(1, arg, args...)
}

func (this *S_Logger) Fatalf(msg string, args ...any) {
	this.S_Logger.Fatalf_(1, msg, args...)
}
