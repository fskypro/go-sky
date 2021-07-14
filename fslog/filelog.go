/**
@copyright: fantasysky 2016
@brief: 将 log 落入文件
@author: fanky
@version: 1.0
@date: 2019-01-05
**/

package fslog

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"fsky.pro/fsos"
	"fsky.pro/fsio"
)

// -----------------------------------------------------------------------------
// inners
// -----------------------------------------------------------------------------
// 获取 log 文件名的日期后缀
func _getLogNamePostfix(utcPostfix bool) string {
	now := time.Now()
	if utcPostfix {
		now = now.UTC()
	}
	return now.Format("2006-01-02")
}

// 获取 log 文件名
func _getLogFilePath(froot, fprefix, fpostfix string) string {
	fname := fmt.Sprintf("%s_%s.log", fprefix, fpostfix)
	return path.Join(froot, fname)
}

// 新建一个 log 文件
func _newLogFile(froot, fprefix, fpostfix string, utcPostfix bool) (string, *os.File, error) {
	logPath := _getLogFilePath(froot, fprefix, fpostfix)
	exists := fsio.IsPathExists(logPath)

	pFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		err = fmt.Errorf("create log file %q fail: %v", logPath, err)
		return "", os.Stdout, err
	}
	if exists {
		now := time.Now()
		if utcPostfix {
			now = now.UTC()
		}
		splitter := fmt.Sprintf("%[3]s%[1]s %[2]s %[1]s%[3]s",
			strings.Repeat("-", 50),
			now.Format("15:04:05"),
			fsos.Endline)
		pFile.WriteString(splitter)
	}
	return logPath, pFile, nil
}

// -----------------------------------------------------------------------------
// new log command
// ---------------------------------------------------------
type s_NewLogCmd struct {
	cmd  string
	args []string
}

func _newLogCmd(cmd string, args ...string) *s_NewLogCmd {
	return &s_NewLogCmd{cmd, args}
}

func (this *s_NewLogCmd) exec(logger *S_FileLogger, log string) {
	if this.cmd == "" {
		return
	}
	args := []string{log}
	args = append(args, this.args...)
	cmd := exec.Command(this.cmd, args...)
	out, err := cmd.CombinedOutput()
	out = bytes.ReplaceAll(bytes.TrimSpace(out), []byte(fsos.Endline), []byte("\n\t"))
	if err != nil {
		logger.Errorf("execute new log file command(%s) fail, error: %v.", this.cmd, err)
	} else {
		logger.Infof("execute new log file command(%s) success! output:%s\t%s", this.cmd, fsos.Endline, out)
	}
}

// -----------------------------------------------------------------------------
// FileLogger
// -----------------------------------------------------------------------------
// 将日志落入文件中，并且每天更换一个文件。
// 日志文件以日期命名
type S_FileLogger struct {
	*S_BaseLogger

	utcPostfix bool                // 是否使用 UTC 时间作为日志文件名的日期后缀
	froot      string              // log 文件所在目录
	fprefix    string              // log 文件名前缀
	fpostfix   string              // log 文件名日期后缀
	logPath    string              // log 路径
	pFile      *os.File            // log 输出
	newLogCmd  *s_NewLogCmd        // 新建 log 文件时，触发该命令
	newLogCB   func(string, error) // 新建 log 文件通知回调
}

// 新建一个 FileLogger。
//	froot：log 文件所在目录
//	fprefix：log 文件名的前缀
//	utcPostfix：是否以 UTC 时间作为 log 文件后缀
//	如果创建 log 文件失败，则返回的第二个错误参数，将会提示失败原因，但 log 仍然可以用，只是 log 默认输出到标准输出
func NewFileLogger(froot string, fprefix string, utcPostfix bool) (*S_FileLogger, error) {
	fpostfix := _getLogNamePostfix(utcPostfix) // 日期后缀
	logPath, w, err := _newLogFile(froot, fprefix, fpostfix, utcPostfix)
	logger := &S_FileLogger{
		S_BaseLogger: NewBaseLogger(w),
		utcPostfix:   utcPostfix,
		froot:        froot,
		fprefix:      fprefix,
		fpostfix:     fpostfix,
		logPath:      logPath,
		pFile:        w,
		newLogCmd:    _newLogCmd(""),
	}
	logger.onBeferPrint = logger._onBeferPrint
	return logger, err
}

// -------------------------------------------------------------------
// prvate
// -------------------------------------------------------------------
// 新建 log 文件
func (this *S_FileLogger) _onBeferPrint(string) {
	postfix := _getLogNamePostfix(this.utcPostfix)
	logPath := _getLogFilePath(this.froot, this.fprefix, this.fpostfix)
	exists := fsio.IsPathExists(logPath)
	if postfix == this.fpostfix && exists {
		return
	}

	logPath, pFile, err := _newLogFile(this.froot, this.fprefix, postfix, this.utcPostfix)
	newLogger := NewBaseLogger(pFile)

	this.Lock()
	if this.logPath != "" {
		this.pFile.Close() // 关闭旧的文件
	}
	this.pFile.Close()
	this.fpostfix = postfix
	this.pFile = pFile
	this.logger = newLogger.logger
	this.logPath = logPath
	this.Unlock()

	if logPath != "" {
		this.newLogCmd.exec(this, logPath)
	}
	if this.newLogCB != nil {
		this.newLogCB(logPath, err)
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_FileLogger) Close() {
	this.Lock()
	defer this.Unlock()
	if this.logPath != "" {
		this.pFile.Close()
	}
	this.logPath = ""
	this.pFile = os.Stdout
	this.logger = NewBaseLogger(this.pFile).logger
}

// 新建 log 文件时将会执行该命令，并把新建的 log 文件作为命令行参数传出
func (this *S_FileLogger) SetNewLogCmd(cmd string, args ...string) {
	this.Lock()
	this.newLogCmd = _newLogCmd(cmd, args...)
	this.Unlock()
	if this.logPath != "" {
		this.newLogCmd.exec(this, this.logPath)
	}
}

// 设置新建 log 文件回调
// cb 第一个参数表示新建文件的路径
// cb 第二个参数表示创建文件失败错误，如果创建成功，则为 nil
func (this *S_FileLogger) SetNewLogCallback(cb func(string, error)) {
	this.newLogCB = cb
}

// --------------------------------------------------------
func (this *S_FileLogger) Debug(vs ...interface{}) {
	this.print(1, "debug", vs)
}

func (this *S_FileLogger) Debugf(format string, vs ...interface{}) {
	this.printf(1, "debug", format, vs)
}

func (this *S_FileLogger) Info(vs ...interface{}) {
	this.print(1, "info", vs)
}

func (this *S_FileLogger) Infof(format string, vs ...interface{}) {
	this.printf(1, "info", format, vs)
}

func (this *S_FileLogger) Warn(vs ...interface{}) {
	this.print(1, "warn", vs)
}

func (this *S_FileLogger) Warnf(format string, vs ...interface{}) {
	this.printf(1, "warn", format, vs)
}

func (this *S_FileLogger) Error(vs ...interface{}) {
	this.print(1, "error", vs)
}

func (this *S_FileLogger) Errorf(format string, vs ...interface{}) {
	this.printf(1, "error", format, vs)
}

func (this *S_FileLogger) Hack(vs ...interface{}) {
	this.print(1, "hack", vs)
}

func (this *S_FileLogger) Hackf(format string, vs ...interface{}) {
	this.printf(1, "hack", format, vs)
}

func (this *S_FileLogger) Panic(vs ...interface{}) {
	this.printChain(1, "panic", vs)
	os.Exit(2)
}

func (this *S_FileLogger) Panicf(format string, vs ...interface{}) {
	this.printChainf(1, "panic", format, vs)
	os.Exit(2)
}

func (this *S_FileLogger) Fatal(vs ...interface{}) {
	this.print(1, "error", vs)
	os.Exit(1)
}

func (this *S_FileLogger) Fatalf(format string, vs ...interface{}) {
	this.printf(1, "fatal", format, vs)
	os.Exit(1)
}

func (this *S_FileLogger) Trace(vs ...interface{}) {
	this.printChain(1, "trace", vs)
}

func (this *S_FileLogger) Tracef(format string, vs ...interface{}) {
	this.printChainf(1, "trace", format, vs)
}
