/**
@copyright: fantasysky 2016
@brief: 实现一个将日志输出到控制台的 CmdLogger
@author: fanky
@version: 1.0
@date: 2018-08-31
**/

package fslog

import "os"

// CmdLogger
type S_CmdLogger struct {
	*S_BaseLogger
}

// New 新建一个默认 CmdLogger
func NewCmdLogger() *S_CmdLogger {
	return &S_CmdLogger{
		S_BaseLogger: NewBaseLogger(os.Stdout),
	}
}

// ----------------------------------------------
func (this *S_CmdLogger) Debug(vs ...interface{}) {
	this.print(1, "debug", vs)
}

func (this *S_CmdLogger) Debugf(format string, vs ...interface{}) {
	this.printf(1, "debug", format, vs)
}

func (this *S_CmdLogger) Info(vs ...interface{}) {
	this.print(1, "info", vs)
}

func (this *S_CmdLogger) Infof(format string, vs ...interface{}) {
	this.printf(1, "info", format, vs)
}

func (this *S_CmdLogger) Warn(vs ...interface{}) {
	this.print(1, "warn", vs)
}

func (this *S_CmdLogger) Warnf(format string, vs ...interface{}) {
	this.printf(1, "warn", format, vs)
}

func (this *S_CmdLogger) Error(vs ...interface{}) {
	this.print(1, "error", vs)
}

func (this *S_CmdLogger) Errorf(format string, vs ...interface{}) {
	this.printf(1, "error", format, vs)
}

func (this *S_CmdLogger) Hack(vs ...interface{}) {
	this.print(1, "hack", vs)
}

func (this *S_CmdLogger) Hackf(format string, vs ...interface{}) {
	this.printf(1, "hack", format, vs)
}

func (this *S_CmdLogger) Panic(vs ...interface{}) {
	this.printChain(1, "panic", vs)
	os.Exit(2)
}

func (this *S_CmdLogger) Panicf(format string, vs ...interface{}) {
	this.printChainf(1, "panic", format, vs)
	os.Exit(2)
}

func (this *S_CmdLogger) Fatal(vs ...interface{}) {
	this.print(1, "fatal", vs)
	os.Exit(1)
}

func (this *S_CmdLogger) Fatalf(format string, vs ...interface{}) {
	this.printf(1, "fatal", format, vs)
	os.Exit(1)
}

func (this *S_CmdLogger) Trace(vs ...interface{}) {
	this.printChain(1, "trace", vs)
}

func (this *S_CmdLogger) Tracef(format string, vs ...interface{}) {
	this.printChainf(1, "trace", format, vs)
}
