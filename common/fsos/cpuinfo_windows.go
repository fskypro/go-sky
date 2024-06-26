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
	"github.com/shirou/gopsutil/cpu"
)

// -----------------------------------------------------------------------------
// cpu 核心数、主频等在 /proc/cpuinfo 中获取
// -----------------------------------------------------------------------------
type S_CpuCore struct {
	Order      int    `json:"processor"`
	VendorID   string `json:"vendor_id"`
	CpuFamily  string `json:"cpuFamily"`
	Model      string `json:"model"`		// 永远为空字符串
	ModelName  string `json:"modelName"`
	MicroCode  string `json:"microcode"`    // 永远为 0
	CpuMHz     int    `json:"cpuMHz"`
	CacheSize  int    `json:"cacheSize"`	// 永远为 0
	PhysicalID string  `json:"physicalID"`
	CpuCores   int `json:"cpuCores"`
	Flags      []string `json:"flags"`
	Bugs       string `json:"bugs"`			// 永远为空字符串
	AddrSize   string `json:"addressSizes"` // 永远为空字符串
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
	infoStats, err := cpu.Info()
	if err != nil {	return nil, err	}
	cpuInfo := &S_CpuInfo{ Cores: []*S_CpuCore{} }
	for _, info := range infoStats {
		core := &S_CpuCore{
			Order: int(info.CPU)+1,
			VendorID: info.VendorID,
			CpuFamily: info.Family,
			ModelName: info.ModelName,
			MicroCode: info.Microcode, // 永远为空
			CpuMHz: int(info.Mhz),
			CacheSize: int(info.CacheSize), // 永远为 0
			PhysicalID: info.PhysicalID,
			CpuCores: int(info.Cores),
			Flags: info.Flags,
		}
		cpuInfo.Cores = append(cpuInfo.Cores, core)
	}
	return cpuInfo, nil
}

// -----------------------------------------------------------------------------
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
	user    float64
	nice    float64
	system  float64
	idle    float64
	iowait  float64
	irq     float64
	softirq float64
	steal   float64
	guest   float64

	UsedPercent float32
	FreePercent float32
}

func (this *S_CoreStat) update(stat *cpu.TimesStat) error {
	this.user = stat.User
	this.nice = stat.Nice
	this.system = stat.System
	this.idle = stat.Idle
	this.iowait = stat.Iowait
	this.irq = stat.Irq
	this.softirq = stat.Softirq
	this.steal = stat.Steal
	this.guest = stat.Guest

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
	cpuStat := &S_CpuStat{
		CoreStats: make([]*S_CoreStat, 0),
	}
	timeStats, err := cpu.Times(true)
	if err != nil { return nil, err }

	var usedPercent float32
	for _, stat := range timeStats {
		coreStat := new(S_CoreStat)
		coreStat.update(&stat)
		usedPercent += coreStat.UsedPercent
		cpuStat.CoreStats = append(cpuStat.CoreStats, coreStat)
	}
	cpuStat.UsedPercent = usedPercent / float32(len(timeStats))
	cpuStat.FreePercent = max(0.0, 100.0-cpuStat.UsedPercent)
	return cpuStat, nil
}
