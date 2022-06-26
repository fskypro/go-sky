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
	"syscall"
	"time"
)

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
	if this.cmd == "" {
		return
	}
	args := []string{log}
	args = append(args, this.args...)
	cmd := exec.Command(this.cmd, args...)
	out, err := cmd.CombinedOutput()
	out = bytes.ReplaceAll(bytes.TrimSpace(out), []byte("\n"), []byte("\n\t"))
	if err != nil {
		logger.Errorf(1, "execute new log file command(%s) fail, error: %v.", this.cmd, err)
	} else {
		logger.Infof(1, "execute new log file command(%s) success! output:\n\t%s", this.cmd, out)
	}
}

// -----------------------------------------------------------------------------
// DayfileLogger
// -----------------------------------------------------------------------------
type S_DayfileLogger struct {
	*S_Logger
	dir     string
	prefix  string
	logPath string
	file    *os.File

	newLogCmd *s_NewLogCmd
	newLogCB  func(string, error)
}

// NewDayfileLogger，新建 DayfileLogger
// root 为 log 的根目录
// filePrefix 为 log 文件名前缀
func NewDayfileLogger(root string, filePrefix string) *S_DayfileLogger {
	if filePrefix == "" {
		filePrefix = "log"
	}
	logger := &S_DayfileLogger{
		dir:       root,
		prefix:    filePrefix,
		newLogCmd: newLogCmd(""),
	}
	logger.S_Logger = newLogger(logger.write)
	return logger
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
	exists := syscall.Access(logPath, syscall.O_RDWR) == nil
	if !exists {
		file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE, 0660)
		if this.newLogCB != nil {
			this.newLogCB(logPath, err)
		}
		if err != nil {
			err = fmt.Errorf("create log file %q fail, %v", logPath, err)
		} else {
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

func (this *S_DayfileLogger) write(t time.Time, msg []byte) {
	if this.file != nil {
		if _, err := this.file.Write(msg); err == nil {
			return
		}
		this.file.Close()
	}
	logPath, file, err := this.newLogFile(t)
	if err != nil {
		os.Stdout.Write(msg)
		os.Stderr.WriteString(err.Error())
	} else {
		file.Write(msg)
		this.logPath = logPath
		this.file = file
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

func (this *S_DayfileLogger) Dispose() {
	this.Lock()
	defer this.Unlock()
	if this.file != nil {
		this.file.Close()
		this.file = nil
		this.logPath = ""
	}
}
