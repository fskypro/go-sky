/**
@copyright: fantasysky 2016
@brief: 全局 log 模块
@author: fanky
@version: 1.0
@date: 2019-01-06
**/

// 实现日志处理器
package fslog

import "os"

// -------------------------------------------------------------------
// inner global
// -------------------------------------------------------------------
// 默认日志处理器
var _logger I_Logger

// -------------------------------------------------------------------
// 初始化
// -------------------------------------------------------------------
func init() {
	initBaseLogger()
	_logger = NewCmdLogger()
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// SetLogger 设置日志处理器
func SetLogger(logger I_Logger) {
	_logger = logger
}

func Close() {
	_logger.Close()
}

// --------------------------------------------------------
// GetUnshields 获取所有开放打印的频道
func GetUnshields() []string {
	return _logger.GetUnshields()
}

// SetOppens 设置打印的频道。
// 如果有频道不存在则设置失败，并返回 false。
//	可选频道有(不区分大小写)：
//	info、warn、debug、error、trace
func SetUnshields(chs ...string) bool {
	return _logger.SetUnshields(chs...)
}

// Open 开启某个指定频道的打印
// 如果指定的频道不存在，则设置失败，并返回 false。
//	可选频道有(不区分大小写)：
//	info、warn、debug、error、trace
func Unshield(ch string) bool {
	return _logger.Unshield(ch)
}

// Shield 屏蔽某个频道的打印
// 如果指定的频道不存在，则设置失败，并返回 false。
//	可选频道有(不区分大小写)：
//	info、warn、debug、error、trace
func Shield(ch string) bool {
	return _logger.Shield(ch)
}

// --------------------------------------------------------
func Debug(vs ...interface{}) {
	_logger.print(1, "debug", vs)
}

func Debugf(format string, vs ...interface{}) {
	_logger.printf(1, "debug", format, vs)
}

func Info(vs ...interface{}) {
	_logger.print(1, "info", vs)
}

func Infof(format string, vs ...interface{}) {
	_logger.printf(1, "info", format, vs)
}

func Warn(vs ...interface{}) {
	_logger.print(1, "warn", vs)
}

func Warnf(format string, vs ...interface{}) {
	_logger.printf(1, "warn", format, vs)
}

func Error(vs ...interface{}) {
	_logger.print(1, "error", vs)
}

func Errorf(format string, vs ...interface{}) {
	_logger.printf(1, "error", format, vs)
}

func Hack(vs ...interface{}) {
	_logger.print(1, "hack", vs)
}

func Hackf(format string, vs ...interface{}) {
	_logger.printf(1, "hack", format, vs)
}

func Panic(vs ...interface{}) {
	_logger.printChain(1, "panic", vs)
	panic("")
}

func Panicf(format string, vs ...interface{}) {
	_logger.printChainf(1, "panic", format, vs)
	panic("")
}

func Fatal(vs ...interface{}) {
	_logger.print(1, "fatal", vs)
	os.Exit(1)
}

func Fatalf(format string, vs ...interface{}) {
	_logger.printf(1, "fatal", format, vs)
	os.Exit(1)
}

func Trace(vs ...interface{}) {
	_logger.printChain(1, "trace", vs)
}

func Tracef(format string, vs ...interface{}) {
	_logger.printChainf(1, "trace", format, vs)
}
