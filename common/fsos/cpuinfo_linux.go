/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: cpu information
@author: fanky
@version: 1.0
@date: 2022-06-06
**/

package fsos

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// -----------------------------------------------------------------------------
// cpu 核心数、主频等在 /proc/cpuinfo 中获取
// -----------------------------------------------------------------------------
type S_CpuCore struct {
	Order      int    `json:"processor"`
	VendorID   string `json:"vendor_id"`
	CpuFamily  string `json:"cpu family"`
	Model      int    `json:"model"`
	ModelName  string `json:"model name"`
	MicroCode  string `json:"microcode"`
	CpuMHz     int    `json:"cpu MHz"`
	CacheSize  int    `json:"cache size"`
	PhysicalID int    `json:"physical id"`
	CpuCores   string `json:"cpu cores"`
	Flags      string `json:"flags"`
	Bugs       string `json:"bugs"`
	AddrSize   string `json:"address sizes"`
}

func (this *S_CpuCore) parse(block string) error {
	lines := strings.Split(block, "\n")
	members := make(map[string]string)
	for _, line := range lines {
		kv := strings.Split(line, ":")
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		members[key] = value
	}
	cache := 0
	segs := strings.Split(members["cache size"], " ")
	if len(segs) > 0 {
		cache, _ = strconv.Atoi(segs[0])
	}

	this.VendorID = members["vendor_id"]
	this.CpuFamily = members["cpu family"]
	this.ModelName = members["model name"]
	this.MicroCode = members["microcode"]
	this.CacheSize = cache
	this.CpuCores = members["cpu cores"]
	this.Flags = members["flags"]
	this.Bugs = members["bugs"]
	this.AddrSize = members["address sizes"]

	this.Order, _ = strconv.Atoi(members["processor"])
	this.Model, _ = strconv.Atoi(members["model"])
	this.PhysicalID, _ = strconv.Atoi(members["physical id"])
	cpuMHz, _ := strconv.ParseFloat(members["cpu MHz"], 32)
	this.CpuMHz = int(cpuMHz)
	return nil
}

type S_CpuInfo struct {
	Cores []*S_CpuCore
}

// 核心个数
func (this *S_CpuInfo) CoreNum() int {
	return len(this.Cores)
}

// 获取指定序号的核心信息
func (this *S_CpuInfo) GetCore(order int) *S_CpuCore {
	for _, core := range this.Cores {
		if order == core.Order {
			return core
		}
	}
	return nil
}

func GetCpuInfo() (*S_CpuInfo, error) {
	data, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, fmt.Errorf("read cpuinfo file fail, %v", err)
	}

	cpuInfo := &S_CpuInfo{
		Cores: make([]*S_CpuCore, 0),
	}
	blocks := strings.Split(string(data), "\n\n")
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		core := new(S_CpuCore)
		if err = core.parse(block); err != nil {
			return nil, fmt.Errorf("parse cpuinfo file error, %v", err)
		}
		cpuInfo.Cores = append(cpuInfo.Cores, core)
	}
	return cpuInfo, nil
}

// -----------------------------------------------------------------------------
// CPU 使用率统计，数据来自于：/proc/stat
// $ cat /proc/stat
//   cpu  92609 13781 80309 301104595 2755 127740 41005 0 0 0
//   cpu0 12597 2846  11307 75328250  781  11008  2820  0 0 0
//   cpu1 27766 4608  26706 75257664  523  34279  20154 0 0 0
//   cpu2 33055 4758  24332 75281648  971  28668  3276  0 0 0
//   cpu3 19189 1567  17962 75237030  478  53784  14754 0 0 0
//   ...
//
// 数值每列代表的是：
//   user        用户态
//   nice        低优先级用户态（nice 的进程）
//   system      内核态（特权模式）
//   idle        空闲
//   iowait      IO等待
//   irq         硬中断 （Interrupt request）
//   softirq     软中断
//   steal       虚拟化环境下其他系统的运行时间
//   guest       运行虚拟宿主机
//   guest_nice  虚拟宿主机nice
//
// CPU 使用率统计方法：
//   work = user +  nice + system + irq + softirq + steal
//   idle = idle + iowait
//   usage = work / (work + idle)
// -----------------------------------------------------------------------------
type S_CoreStat struct {
	Order   int
	user    int64
	nice    int64
	system  int64
	idle    int64
	iowait  int64
	irq     int64
	softirq int64
	steal   int64
	guest   int64

	UsedPercent float32
	FreePercent float32
}

func (this *S_CoreStat) parse(vars []string) error {
	if len(vars) < 10 {
		return errors.New("invalid cpu stat file")
	}
	sorder := strings.TrimLeft(vars[0], "cpu")
	if sorder == "" {
		this.Order = -1
	} else {
		order, err := strconv.Atoi(sorder)
		if err != nil {
			return errors.New("ivalid cpu stat file")
		}
		this.Order = order
	}
	this.user, _ = strconv.ParseInt(vars[1], 10, 64)
	this.nice, _ = strconv.ParseInt(vars[2], 10, 64)
	this.system, _ = strconv.ParseInt(vars[3], 10, 64)
	this.idle, _ = strconv.ParseInt(vars[4], 10, 64)
	this.iowait, _ = strconv.ParseInt(vars[5], 10, 64)
	this.irq, _ = strconv.ParseInt(vars[6], 10, 64)
	this.softirq, _ = strconv.ParseInt(vars[7], 10, 64)
	this.steal, _ = strconv.ParseInt(vars[8], 10, 64)
	this.guest, _ = strconv.ParseInt(vars[9], 10, 64)

	useds := this.user + this.nice + this.system + this.irq + this.softirq + this.steal
	frees := this.idle + this.iowait
	this.UsedPercent = float32(float64(useds)/float64(useds+frees)) * 100
	this.FreePercent = 100 - this.UsedPercent
	return nil
}

type S_CpuStat struct {
	UsedPercent float32
	FreePercent float32
	CoreStats   []*S_CoreStat
}

func GetCpuStat() (*S_CpuStat, error) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, fmt.Errorf("read cpu stat file fail, %v", err)
	}
	cpuStat := &S_CpuStat{
		CoreStats: make([]*S_CoreStat, 0),
	}

	re := regexp.MustCompile(`\s+`)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "cpu") {
			continue
		}
		vars := re.Split(strings.TrimSpace(line), -1)
		core := new(S_CoreStat)
		if err = core.parse(vars); err != nil {
			return nil, err
		}
		if core.Order >= 0 {
			cpuStat.CoreStats = append(cpuStat.CoreStats, core)
		} else {
			cpuStat.UsedPercent = core.UsedPercent
			cpuStat.FreePercent = core.FreePercent
		}
	}
	return cpuStat, nil
}
