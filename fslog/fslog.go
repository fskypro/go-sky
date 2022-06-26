/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: logger
@author: fanky
@version: 1.0
@date: 2022-06-25
**/

package fslog

type I_Logger interface {
	Debug(int, any, ...any)
	Debugf(int, string, ...any)

	Info(int, any, ...any)
	Infof(int, string, ...any)

	Warn(int, any, ...any)
	Warnf(int, string, ...any)

	Error(int, any, ...any)
	Errorf(int, string, ...any)

	Hack(int, any, ...any)
	Hackf(int, string, ...any)

	Critical(int, any, ...any)
	Criticalf(int, string, ...any)

	Trace(int, any, ...any)
	Tracef(int, string, ...any)

	Panic(int, any, ...any)
	Panicf(int, string, ...any)

	Fatal(int, any, ...any)
	Fatalf(int, string, ...any)

	Shield(string, ...string)
	Unshield(string, ...string)
}

// -----------------------------------------------------------------------------
// package interfaces
// -----------------------------------------------------------------------------
var logger I_Logger

func init() {
	logger = NewStdoutLogger()
}

func SetLogger(l I_Logger) {
	logger = l
}

func Debug(arg any, args ...any) {
	logger.Debug(1, arg, args...)
}

func Debugf(msg string, args ...any) {
	logger.Debugf(1, msg, args...)
}

func Info(arg any, args ...any) {
	logger.Info(1, arg, args...)
}

func Infof(msg string, args ...any) {
	logger.Infof(1, msg, args...)
}

func Warn(arg any, args ...any) {
	logger.Warn(1, arg, args...)
}

func Warnf(msg string, args ...any) {
	logger.Warnf(1, msg, args...)
}

func Error(arg any, args ...any) {
	logger.Error(1, arg, args...)
}

func Errorf(msg string, args ...any) {
	logger.Errorf(1, msg, args...)
}

func Hack(arg any, args ...any) {
	logger.Hack(1, arg, args...)
}

func Hackf(msg string, args ...any) {
	logger.Hackf(1, msg, args...)
}

func Critical(arg any, args ...any) {
	logger.Critical(1, arg, args...)
}

func Criticalf(msg string, args ...any) {
	logger.Criticalf(1, msg, args...)
}

func Trace(arg any, args ...any) {
	logger.Trace(1, arg, args...)
}

func Tracef(msg string, args ...any) {
	logger.Tracef(1, msg, args...)
}

func Panic(arg any, args ...any) {
	logger.Panic(1, arg, args...)
}

func Panicf(msg string, args ...any) {
	logger.Panicf(1, msg, args...)
}

func Fatal(arg any, args ...any) {
	logger.Fatal(1, arg, args...)
}

func Fatalf(msg string, args ...any) {
	logger.Fatalf(1, msg, args...)
}

// ---------------------------------------------------------
func Shield(lv string, lvs ...string) {
	logger.Shield(lv, lvs...)
}

func Unshield(lv string, lvs ...string) {
	logger.Unshield(lv, lvs...)
}
