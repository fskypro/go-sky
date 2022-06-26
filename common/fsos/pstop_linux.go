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
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// 进程资源消耗信息
type S_PSInfo struct {
	Name  string `json:"name"`
	Tick  int    `json:"tick"`
	Valid bool   `json:"running"`

	Pid  int     `top:"pid" json:"pid"`
	User string  `top:"user" json:"user"`
	Cpu  float32 `top:"cpu" json:"cpu"`
	Mem  float32 `top:"mem" json:"mem"`
	Cmd  string  `top:"cmd" json:"cmd"`
}

func NewPSInfo(name string, pid int) *S_PSInfo {
	return &S_PSInfo{
		Name: name,
		Pid:  pid,
	}
}

func (this *S_PSInfo) update(user, cpu, mem, cmd string) {
	fcpu, _ := strconv.ParseFloat(cpu, 32)
	fmem, _ := strconv.ParseFloat(mem, 32)
	this.User = user
	this.Cpu = float32(fcpu)
	this.Mem = float32(fmem)
	this.Cmd = cmd
	this.Valid = true
}

func (this *S_PSInfo) reset(tick int) {
	this.Tick = tick
	this.User = ""
	this.Cpu = 0
	this.Mem = 0
	this.Cmd = ""
	this.Valid = false
}

// -------------------------------------------------------------------
// PSTop
// -------------------------------------------------------------------
type S_PSTop struct {
	keyIndexs map[string]int
	resp      *regexp.Regexp
	tick      int
}

// keysHeader 指出 PSInfo 中每个键的对应 top 命令头中的哪个列
// 常见：
//	"pid":  "PID",     // 进程号
//	"user": "USER",    // 进程所属用户
//	"cpu":  "%CPUj,    // 进程占用的 Cpu 百分比
//	"mem":  "%MEM",    // 进程占用的内存百分比
//	"cmd":  "COMMAND", // 进程启动时的命令
func NewPSTop(keysHeader map[string]string) (*S_PSTop, error) {
	psTop := &S_PSTop{
		keyIndexs: make(map[string]int),
		resp:      regexp.MustCompile("\\s+"),
	}

	headerKeys := make(map[string]string)
	for key, header := range keysHeader {
		headerKeys[header] = key
		psTop.keyIndexs[key] = -1
	}

	err := psTop.execTop(func(cols []string) bool {
		count := 0
		for index, col := range cols {
			if key, ok := headerKeys[col]; ok {
				psTop.keyIndexs[key] = index
				count++
			}
		}
		return count < len(keysHeader)
	})
	if err != nil {
		return nil, fmt.Errorf("new fail, system doses not support top command")
	}
	for key := range keysHeader {
		if psTop.keyIndexs[key] < 0 {
			return nil, fmt.Errorf("invalid top command header %q", keysHeader[key])
		}
	}
	if psTop.keyIndexs["pid"] < 0 {
		return nil, fmt.Errorf("pid key map header must be indicate")
	}
	return psTop, nil
}

// ---------------------------------------------------------
// private
// ---------------------------------------------------------
func (this *S_PSTop) execTop(f func([]string) bool) error {
	cmd := exec.Command("bash", "-c", "top -cbn1")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("execute top command fail, %v", err)
	}
	buff := bytes.NewBuffer(out)
	for {
		row, err := buff.ReadString('\n')
		if err != nil {
			break
		}
		line := string(row)
		cols := this.resp.Split(strings.TrimSpace(line), -1)
		if len(cols) < len(this.keyIndexs) {
			continue
		}
		if !f(cols) {
			break
		}
	}
	return nil
}

// ---------------------------------------------------------
// public
// ---------------------------------------------------------
func (this *S_PSTop) Execute(psInfos []*S_PSInfo) error {
	for _, psInfo := range psInfos {
		psInfo.reset(this.tick)
	}
	err := this.execTop(func(cols []string) bool {
		spid := cols[this.keyIndexs["pid"]]
		pid, _ := strconv.Atoi(spid)
		if pid == 0 {
			return true
		}
		for _, psInfo := range psInfos {
			if psInfo.Pid == pid {
				psInfo.update(
					cols[this.keyIndexs["user"]],
					cols[this.keyIndexs["cpu"]],
					cols[this.keyIndexs["mem"]],
					cols[this.keyIndexs["cmd"]],
				)
				break
			}
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("execute top command fail, %v", err)
	}
	this.tick++
	return nil
}

func (this *S_PSTop) Reset() {
	this.tick = 0
}
