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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fsky.pro/fstime"
)

// 文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !os.IsNotExist(err)
}

// -----------------------------------------------------------------------------
// new log command
// -----------------------------------------------------------------------------
type s_NewLogCmd struct {
	cmd  string
	args []string
}

func newLogCmd(cmd string, args ...string) *s_NewLogCmd {
	return &s_NewLogCmd{cmd, args}
}

func (this *s_NewLogCmd) exec(logger *S_DayfileLogger, log string) {
	if this.cmd == "" { return }
	args := []string{log}
	args = append(args, this.args...)
	cmd := exec.Command(this.cmd, args...)
	out, err := cmd.CombinedOutput()
	out = bytes.ReplaceAll(bytes.TrimSpace(out), []byte("\n"), []byte("\n\t"))
	if err != nil {
		go logger.Errorf("execute new log file command(%s) fail, error: %v.", this.cmd, err)
	} else {
		go logger.Infof("execute new log file command(%s) success! output:\n\t%s", this.cmd, out)
	}
}

// -----------------------------------------------------------------------------
// DayfileLogger
// -----------------------------------------------------------------------------
type S_DayfileLogger struct {
	*S_Logger
	dir         string
	prefix      string
	logPath     string
	file        *os.File
	nextDayTime time.Time

	newLogCmd   *s_NewLogCmd
	newLogCB    func(string, error)
	newLinkFile string
}

// NewDayfileLogger，新建 DayfileLogger
// root 为 log 的根目录
// filePrefix 为 log 文件名前缀
func NewDayfileLogger(root string, filePrefix string) *S_DayfileLogger {
	if filePrefix == "" {
		filePrefix = "log"
	}
	logger := &S_DayfileLogger{
		dir:         root,
		prefix:      filePrefix,
		nextDayTime: fstime.Dawn(time.Now()).AddDate(0, 0, 1),
		newLogCmd:   newLogCmd(""),
	}
	logger.S_Logger = NewLogger(logger.write)
	logger.logPath, logger.file, _ = logger.newLogFile(time.Now())
	return logger
}

func (this *S_DayfileLogger) GetRoot() string {
	return this.dir
}

func (this *S_DayfileLogger) GetPrefix() string {
	return this.prefix
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
func (this *S_DayfileLogger) getLogFilePath(t time.Time) string {
	fname := fmt.Sprintf("%s_%s.log", this.prefix, t.Format("2006-01-02"))
	return filepath.Join(this.dir, fname)
}

// 新建一个 log 文件
func (this *S_DayfileLogger) newLogFile(t time.Time) (string, *os.File, error) {
	logPath := this.getLogFilePath(t)
	exists := fileExists(logPath)
	if !exists {
		file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE, 0666)
		if this.newLogCB != nil {
			this.newLogCB(logPath, err)
		}
		if err != nil {
			err = fmt.Errorf("create log file %q fail, %v", logPath, err)
		} else {
			go this.linkTo(logPath)
			this.newLogCmd.exec(this, logPath)
		}
		return logPath, file, err
	}

	// 文件已经存在
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		err = fmt.Errorf("open log file %q fail, %v", logPath, err)
		return logPath, file, err
	}
	file.WriteString("\n")
	splitter := fmt.Sprintf("%[1]s %[2]s %[1]s", strings.Repeat("-", 50), t.Format("15:04:05"))
	file.WriteString(splitter)
	file.WriteString("\n")
	return logPath, file, nil
}

// 父类中的 send 函数中已经 lock，因此这里不需要再上锁了
func (this *S_DayfileLogger) write(t time.Time, lv T_Level, msg []byte) {
	now := this.nowTime()
	if this.file != nil && now.Before(this.nextDayTime) {
		if _, err := this.file.Write(msg); err == nil {
			return
		}
	}
	if this.file != nil {
		this.file.Close()
	}
	logPath, file, err := this.newLogFile(t)
	if err != nil {
		os.Stdout.Write(msg)
		os.Stderr.WriteString(now.Format("[ERROR]|2006/01/02 15:04:05.999999 "))
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
	} else {
		file.Write(msg)
		this.logPath = logPath
		this.file = file
	}
	this.nextDayTime = fstime.Dawn(time.Now().AddDate(0, 0, 1).Add(time.Hour))
}

func (this *S_DayfileLogger) linkTo(logPath string) {
	if this.newLinkFile == "" { return }
	_, err := os.Lstat(this.newLinkFile)
	if err == nil {
		os.Remove(this.newLinkFile)
	}

	err = os.Symlink(logPath, this.newLinkFile)
	if err != nil {
		this.Errorf("link new log file %q to %q fail, %v", logPath, this.newLinkFile, err)
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
// 新建 log 文件时将会执行该命令，并把新建的 log 文件作为命令行参数传出
func (this *S_DayfileLogger) SetNewLogCmd(cmd string, args ...string) {
	this.newLogCmd = newLogCmd(cmd, args...)
	if this.logPath != "" {
		this.newLogCmd.exec(this, this.logPath)
	}
}

// 设置新建 log 文件回调
// cb 第一个参数表示新建文件的路径
// cb 第二个参数表示创建文件失败错误，如果创建成功，则为 nil
func (this *S_DayfileLogger) SetNewLogCallback(cb func(string, error)) {
	this.newLogCB = cb
}

// 对新产生对 log 文件进行软连接到指定地方
func (this *S_DayfileLogger) SetNewLogLinkFile(path string) {
	this.newLinkFile = path
}

func (this *S_DayfileLogger) Close() {
	this.Lock()
	defer this.Unlock()
	if this.file != nil {
		this.file.Close()
		this.file = nil
		this.logPath = ""
	}
}

// ---------------------------------------------------------
func (this *S_DayfileLogger) Debug(arg any, args ...any) {
	this.S_Logger.Debug_(1, arg, args...)
}

func (this *S_DayfileLogger) Debugf(msg string, args ...any) {
	this.S_Logger.Debugf_(1, msg, args...)
}

func (this *S_DayfileLogger) Info(arg any, args ...any) {
	this.S_Logger.Info_(1, arg, args...)
}

func (this *S_DayfileLogger) Infof(msg string, args ...any) {
	this.S_Logger.Infof_(1, msg, args...)
}

func (this *S_DayfileLogger) Notic(arg any, args ...any) {
	this.S_Logger.Notic_(1, arg, args...)
}

func (this *S_DayfileLogger) Noticf(msg string, args ...any) {
	this.S_Logger.Noticf_(1, msg, args...)
}

func (this *S_DayfileLogger) Warn(arg any, args ...any) {
	this.S_Logger.Warn_(1, arg, args...)
}

func (this *S_DayfileLogger) Warnf(msg string, args ...any) {
	this.S_Logger.Warnf_(1, msg, args...)
}

func (this *S_DayfileLogger) Error(arg any, args ...any) {
	this.S_Logger.Error_(1, arg, args...)
}

func (this *S_DayfileLogger) Errorf(msg string, args ...any) {
	this.S_Logger.Errorf_(1, msg, args...)
}

func (this *S_DayfileLogger) Hack(arg any, args ...any) {
	this.S_Logger.Hack_(1, arg, args...)
}

func (this *S_DayfileLogger) Hackf(msg string, args ...any) {
	this.S_Logger.Hackf_(1, msg, args...)
}

func (this *S_DayfileLogger) Illeg(arg any, args ...any) {
	this.S_Logger.Illeg_(1, arg, args...)
}

func (this *S_DayfileLogger) Illegf(msg string, args ...any) {
	this.S_Logger.Illegf_(1, msg, args...)
}

func (this *S_DayfileLogger) Critical(arg any, args ...any) {
	this.S_Logger.Critical_(1, arg, args...)
}

func (this *S_DayfileLogger) Criticalf(msg string, args ...any) {
	this.S_Logger.Criticalf_(1, msg, args...)
}

func (this *S_DayfileLogger) Trace(arg any, args ...any) {
	this.S_Logger.Trace_(1, arg, args...)
}

func (this *S_DayfileLogger) Tracef(msg string, args ...any) {
	this.S_Logger.Tracef_(1, msg, args...)
}

func (this *S_DayfileLogger) Panic(arg any, args ...any) {
	this.S_Logger.Panic_(1, arg, args...)
}

func (this *S_DayfileLogger) Panicf(msg string, args ...any) {
	this.S_Logger.Panicf_(1, msg, args...)
}

func (this *S_DayfileLogger) Fatal(arg any, args ...any) {
	this.S_Logger.Fatal_(1, arg, args...)
}

func (this *S_DayfileLogger) Fatalf(msg string, args ...any) {
	this.S_Logger.Fatalf_(1, msg, args...)
}
