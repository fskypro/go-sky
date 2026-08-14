/**
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: cpu info
@author: fanky
@version: 1.0
@date: 2026-02-04
**/

package fshardware

type S_CpuCore struct {
	Order      int      `json:"processor"`    // 核心序号
	VendorID   string   `json:"vendor_id"`    // CPU 厂商ID：Intel、AMD、Apple、ARM 等
	CpuFamily  string   `json:"cpuFamily"`    // CPU 家族编号：如Intel的"6"代表 Core 系列
	Model      string   `json:"model"`        // 厂商内部的具体型号代码，对于 Windows 永远为空字符串
	ModelName  string   `json:"modelName"`    // CPU 完整型号名称：如"Intel Core i7-9750H"、"Apple M3 Max"
	MicroCode  string   `json:"microcode"`    // CPU 内部固件版本，修复bug和漏洞，对于 Windows/MacOS 永远为 0
	CpuMHz     int      `json:"cpuMHz"`       // CPU 当前运行频率（兆赫兹）：如 3200 表示 3.2GHz
	CacheSize  int      `json:"cacheSize"`    // L3 缓存大小，对于 Windows 永远为 0，对于 MacOS 等于 L3
	PhysicalID int      `json:"physicalID"`   // 物理 CPU 插槽编号（双路服务器会有0和1）
	CpuCores   int      `json:"cpuCores"`     // 物理核心数（单颗 CPU 内的核心数），MacOS通常单路
	Flags      []string `json:"flags"`        // CPU 特性标志，对于 MacOS 为空
	Bugs       string   `json:"bugs"`         // 对于 Windows/MacOS 永远为空字符串
	AddrSize   string   `json:"addressSizes"` // 对于 Windows/MacOS  永远为空字符串
	L1Cache    int      `json:"l1_cache_kb"`  // 对于非 MacOS 永远为 0
	L2Cache    int      `json:"l2_cache_kb"`  // 对于非 MacOS 永远为 0
	L3Cache    int      `json:"l3_cache_kb"`  // 对于非 MacOS 永远为 0
}

// -------------------------------------------------------------------
// CpuInfo
// -------------------------------------------------------------------
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
