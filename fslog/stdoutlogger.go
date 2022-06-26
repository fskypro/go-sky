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
