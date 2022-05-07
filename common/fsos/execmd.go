/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: excute system command
@author: fanky
@version: 1.0
@date: 2022-05-06
**/

package fsos

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"sync"
)

type S_Cmd struct {
	sync.Locker
	cmd    *exec.Cmd
	ctx    context.Context
	cancel func()

	OnOut func(string)
	OnErr func(string)
}

func NewCmd(name string, args ...string) *S_Cmd {
	cmd := new(S_Cmd)
	cmd.ctx, cmd.cancel = context.WithCancel(context.Background())
	cmd.cmd = exec.CommandContext(cmd.ctx, name, args...)
	return cmd
}

func (this *S_Cmd) readOutput(reader io.ReadCloser) {
	r := bufio.NewReader(reader)
	for {
		select {
		case <-this.ctx.Done():
			return
		default:
			line, _, err := r.ReadLine()
			if err == io.EOF {
				return
			} else if err != nil {
				if this.OnErr != nil {
					this.OnErr(err.Error())
				}
				return
			} else if this.OnOut != nil {
				this.OnOut(string(line))
			}
		}
	}
}

func (this *S_Cmd) Run() error {
	stdout, err := this.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := this.cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := this.cmd.Start(); err != nil {
		return err
	}
	go this.readOutput(stdout)
	go this.readOutput(stderr)
	return nil
}

func (this *S_Cmd) Wait() error {
	return this.cmd.Wait()
}

func (this *S_Cmd) Stop() {
	this.cancel()
}
