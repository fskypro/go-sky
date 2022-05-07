//go:build linux
// +build linux

/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: fetch ps top command info, only for linux
@author: fanky
@version: 1.0
@date: 2022-05-06
**/

package fsos

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type S_PSTop struct {
	pid          string
	cmd          *S_Cmd
	header       map[string]int
	interval     int
	isInitHeader bool
	re           *regexp.Regexp
	tick         int
	hasPs        bool

	OnTick func(map[string]any)
	OnErr  func(string)
}

func NewPSTop(pid int, interval int) *S_PSTop {
	if interval < 1 {
		interval = 3
	}
	cmd := NewCmd("bash", "-c", fmt.Sprintf("top -b -d%d", interval))
	psTop := &S_PSTop{
		pid:      strconv.Itoa(pid),
		cmd:      cmd,
		interval: interval,
		header:   make(map[string]int),
		re:       regexp.MustCompile("\\s+"),
	}

	psTop.header = map[string]int{
		"PID":     -1,
		"USER":    -1,
		"%CPU":    -1,
		"%MEM":    -1,
		"COMMAND": -1,
	}
	cmd.OnOut = psTop.onOutput
	cmd.OnErr = psTop.onError
	return psTop
}

func (this *S_PSTop) onOutput(line string) {
	cols := this.re.Split(strings.TrimSpace(line), -1)
	if len(cols) < len(this.header) {
		return
	}
	if cols[0] == "top" {
		if this.tick > 0 && !this.hasPs {
			this.onError(fmt.Sprintf("process which pid=%s is not exists", this.pid))
		}
		this.tick++
		this.hasPs = false
		return
	}
	if !this.isInitHeader {
	L:
		for n := range this.header {
			for index, c := range cols {
				if c == n {
					this.header[n] = index
					continue L
				}
			}
			return
		}
		this.isInitHeader = true
		return
	}
	if this.pid != cols[this.header["PID"]] {
		return
	} else {
		this.hasPs = true
	}

	cpu, err := strconv.ParseFloat(cols[this.header["%CPU"]], 32)
	if err != nil {
		cpu = -1
	}
	mem, err := strconv.ParseFloat(cols[this.header["%MEM"]], 32)
	if err != nil {
		mem = -1
	}
	psInfo := map[string]any{
		"TICK":    this.tick,
		"PID":     this.pid,
		"%CPU":    cpu,
		"%MEM":    mem,
		"COMMAND": cols[this.header["COMMAND"]],
	}
	if this.OnTick != nil {
		this.OnTick(psInfo)
	}
}

func (this *S_PSTop) onError(err string) {
	if this.OnErr != nil {
		this.OnErr(err)
	}
}

func (this *S_PSTop) Run() error {
	err := this.cmd.Run()
	if err != nil {
		return err
	}
	return this.cmd.Wait()
}

func (this *S_PSTop) Stop() {
	this.cmd.Stop()
}
