/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: plugin client
@author: fanky
@version: 1.0
@date: 2023-05-03
**/

package fsrpcplugin

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	"fsky.pro/fslog"
)

// -----------------------------------------------------------------------------
// client config
// -----------------------------------------------------------------------------
type S_ClientConfig struct {
	SocketFile  string         // unix socket file
	Logger      fslog.I_Logger // logger
	Cmd         *exec.Cmd      // start plugin server command
	DialTimeout time.Duration  // 插件响应超时时间
}

// -----------------------------------------------------------------------------
// client
// -----------------------------------------------------------------------------
type S_Client struct {
	rpcClient *rpc.Client
	conf      *S_ClientConfig
	process   *os.Process
}

func NewClient(cfg *S_ClientConfig) (*S_Client, error) {
	client := &S_Client{conf: cfg}
	if cfg.SocketFile == "" {
		return nil, errors.New("unix socket file must be indicated")
	}
	if cfg.Logger == nil {
		return nil, errors.New("logger must be indicated")
	}
	if cfg.Cmd == nil {
		return nil, errors.New("start plugin server command must't be nil")
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = time.Second * 3
	} else if cfg.DialTimeout < time.Second {
		cfg.DialTimeout = time.Second
	}
	return client, nil
}

// -------------------------------------------------------------------
// private
// -------------------------------------------------------------------
// 将插件进程的输出转发到 client 进程 logger 中
func (this *S_Client) pipout(r io.Reader) {
	const bufferSize = 64 * 1024
	reader := bufio.NewReaderSize(r, bufferSize)
	msg := []byte{}
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			fslog.Errorf("read plugin %s stdout/stderr EOF", this.conf.Cmd.Path)
			return
		} else if err != nil {
			return
		}
		if isPrefix {
			msg = append(msg, line...)
			continue
		}

		level := "info"
		segs := bytes.SplitN(msg, []byte(">"), 7)
		if len(segs) > 1 {
			level = string(segs[0])
			msg = segs[1]
		}
		msg = bytes.ReplaceAll(msg, []byte{1}, []byte{'\n'})
		this.conf.Logger.Direct(level, msg)
		msg = []byte{}
	}
}

// -------------------------------------------------------------------
// public
// -------------------------------------------------------------------
func (this *S_Client) Start() error {
	cmd := this.conf.Cmd
	cmd.Env = append(cmd.Env, fmt.Sprintf("GO_PLUGIN_UNIX_SOCKET_FILE=%s", this.conf.SocketFile))

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("get cmd stdout pip fail, %v", err)
	}
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("get cmd stderr pip fail, %v", err)
	}
	go this.pipout(cmdStdout)
	go this.pipout(cmdStderr)
	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr

	// 启动插件
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start plugin fail, %v", err)
	}
	this.process = cmd.Process

	// 等待插件结束
	go func() {
		cmd.Wait()
		fslog.Infof("plugin %q has exited", cmd.Path)
	}()
	return nil
}

func (this *S_Client) Dial(ctx context.Context) error {
	// 外部结束控制
	go func() {
		<-ctx.Done()
		this.Kill()
	}()

	var err error
	endTime := time.Now().Add(this.conf.DialTimeout)
	for {
		this.rpcClient, err = rpc.Dial("unix", this.conf.SocketFile)
		if err == nil { return nil }
		if time.Now().After(endTime) {
			return fmt.Errorf("dial to unix file %q fail, %v", this.conf.SocketFile, err)
		}
		time.Sleep(time.Microsecond * 200)
	}
}

func (this *S_Client) Call(method string, arg any, reply *string) error {
	if this.rpcClient == nil {
		return fmt.Errorf("plugin client is not ready")
	}
	return this.rpcClient.Call(method, arg, reply)
}

func (this *S_Client) Kill() error {
	if this.process != nil {
		return this.process.Kill()
	}
	return nil
}
