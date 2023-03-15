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
	ToggleSite(bool) // 是否输出打印文件

	Debug_(int, any, ...any)
	Debugf_(int, string, ...any)

	Info_(int, any, ...any)
	Infof_(int, string, ...any)

	Notic_(int, any, ...any)
	Noticf_(int, string, ...any)

	Warn_(int, any, ...any)
	Warnf_(int, string, ...any)

	Error_(int, any, ...any)
	Errorf_(int, string, ...any)

	Hack_(int, any, ...any)
	Hackf_(int, string, ...any)

	Critical_(int, any, ...any)
	Criticalf_(int, string, ...any)

	Trace_(int, any, ...any)
	Tracef_(int, string, ...any)

	Panic_(int, any, ...any)
	Panicf_(int, string, ...any)

	Fatal_(int, any, ...any)
	Fatalf_(int, string, ...any)

	Shield(...string)
	ShieldAll()
	Unshield(...string)
	UnshieldAll()
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

func UsedLogger() I_Logger {
	return logger
}

func Debug(arg any, args ...any) {
	logger.Debug_(1, arg, args...)
}

func Debugf(msg string, args ...any) {
	logger.Debugf_(1, msg, args...)
}

func Info(arg any, args ...any) {
	logger.Info_(1, arg, args...)
}

func Infof(msg string, args ...any) {
	logger.Infof_(1, msg, args...)
}

func Notic(arg any, args ...any) {
	logger.Notic_(1, arg, args...)
}

func Noticf(msg string, args ...any) {
	logger.Noticf_(1, msg, args...)
}

func Warn(arg any, args ...any) {
	logger.Warn_(1, arg, args...)
}

func Warnf(msg string, args ...any) {
	logger.Warnf_(1, msg, args...)
}

func Error(arg any, args ...any) {
	logger.Error_(1, arg, args...)
}

func Errorf(msg string, args ...any) {
	logger.Errorf_(1, msg, args...)
}

func Hack(arg any, args ...any) {
	logger.Hack_(1, arg, args...)
}

func Hackf(msg string, args ...any) {
	logger.Hackf_(1, msg, args...)
}

func Critical(arg any, args ...any) {
	logger.Critical_(1, arg, args...)
}

func Criticalf(msg string, args ...any) {
	logger.Criticalf_(1, msg, args...)
}

func Trace(arg any, args ...any) {
	logger.Trace_(1, arg, args...)
}

func Tracef(msg string, args ...any) {
	logger.Tracef_(1, msg, args...)
}

func Panic(arg any, args ...any) {
	logger.Panic_(1, arg, args...)
}

func Panicf(msg string, args ...any) {
	logger.Panicf_(1, msg, args...)
}

func Fatal(arg any, args ...any) {
	logger.Fatal_(1, arg, args...)
}

func Fatalf(msg string, args ...any) {
	logger.Fatalf_(1, msg, args...)
}

// ---------------------------------------------------------
func Shield(lvs ...string) {
	logger.Shield(lvs...)
}

func ShieldAll() {
	logger.ShieldAll()
}

func Unshield(lvs ...string) {
	logger.Unshield(lvs...)
}

func UnshieldAll() {
	logger.UnshieldAll()
}
