/**
@copyright: fantasysky 2016
@brief: 插件版本管理封装
@author: fanky
@version: 1.0
@date: 2019-04-11
**/

/*
要求用 json 格式配置罗列插件版本，格式如下：
{
	"name": "插件名称",

	"versions": [
		{
			"ver": "v1.0.0",
			"dsp": "版本描述",
			"symbols": ["symbol1", "symbol2", ...]
		},
		{
			"ver": "v1.0.1",
			"dsp": "版本描述",
			"symbols": ["symbol1", "symbol2", ...]
		},
		...
		]
	}
}
1、“name” 为插件描述名称，可选
2、“soName” 为插件 so 名称
3、“root” 如果给出的是相对路径，则为相对可执行程序的相对路径
4、“symbols” 为公开成员列表，必须大写开头
5、so 插件文件，必须放在与配置文件同路径下
6、so 插件文件的名称为：插件名称-版本，譬如：plugin-v1.0.1.so
*/
package fsplugin

import (
	"errors"
	"fmt"
	"path/filepath"
	"plugin"
)
import "fsky.pro/fsjson/jsonex"

// -------------------------------------------------------------------
// 版本配置
// -------------------------------------------------------------------
type VersionInfo struct {
	Version  string   `json:"ver"`
	Descript string   `json:"dsp"`
	Symbols  []string `json:"symbols"`
}

type VersionConfig struct {
	Name     string         `json:"name"`
	SoName   string         `json:"soname"`
	Versions []*VersionInfo `json:"versions"`
}

// -------------------------------------------------------------------
// 插件对象
// -------------------------------------------------------------------
type Plugin struct {
	config      string         // 配置文件路径
	Name        string         // 插件名称
	SoName      string         // so 名称
	Path        string         // 插件路径
	VersionList []*VersionInfo // 插件列表

	Version string                   // 当前使用版本
	symbols map[string]plugin.Symbol // 插件公开符号列表（必须大写开头）
}

func NewPlugin(config string) (pln *Plugin, err error) {
	vc := new(VersionConfig)
	err = jsonex.Load(config, vc)
	if err != nil {
		err = errors.New(fmt.Sprintf("can't read plugin config file: %s", err.Error()))
		return
	}
	pln = &Plugin{
		config:      config,
		Name:        vc.Name,
		SoName:      vc.SoName,
		VersionList: vc.Versions,
	}
	pln.Path, _ = filepath.Split(config)
	return
}

// 获取指定版本的 so 路径
func (this *Plugin) GetVersionFile(v string) string {
	return filepath.Join(this.Path, fmt.Sprintf("%s-%s.so", this.SoName, v))
}

// 获取版本数量
func (this *Plugin) GetVersionCount() int {
	return len(this.VersionList)
}

// 判断指定版本是否是当前使用的版本
func (this *Plugin) IsUsingVersion(v string) bool {
	return v == this.Version
}

// ---------------------------------------------------------
// 重新打开插件配置初始化插件
func (this *Plugin) Reopen() error {
	pln, err := NewPlugin(this.config)
	if err != nil {
		return err
	}
	this.Name = pln.Name
	this.SoName = pln.SoName
	this.VersionList = pln.VersionList
	return nil
}

// 加载指定版本
func (this *Plugin) LoadVersion(version string) error {
	var vInfo *VersionInfo
	for _, v := range this.VersionList {
		if v.Version == version {
			vInfo = v
			break
		}
	}
	if vInfo == nil {
		return errors.New(fmt.Sprintf("plugin version %q is not exists!", version))
	}

	path := this.GetVersionFile(version)
	pln, err := plugin.Open(path)
	if err != nil {
		return errors.New(fmt.Sprintf("open plugin file(%q) fail: %s", path, err.Error()))
	}

	symbols := make(map[string]plugin.Symbol)
	for _, symb := range vInfo.Symbols {
		symbol, err := pln.Lookup(symb)
		if err != nil {
			return errors.New(fmt.Sprintf("symbol %q is not exist in plugin(%q), load fail: %s", symb, path, err.Error()))
		} else {
			symbols[symb] = symbol
		}
	}
	this.Version = version
	this.symbols = symbols
	return nil
}

// 加载最新版本
// 注意：
//	这里的最新，是指配置列表中排在最后面的版本，函数本身并不会对版本号进行对比，所以配置时，要将新版本排布在旧版本的后面
func (this *Plugin) LoadLatestVersion() error {
	count := len(this.VersionList)
	if count == 0 {
		return errors.New(fmt.Sprintf("no versions in config list to load plugin"))
	}
	return this.LoadVersion(this.VersionList[count-1].Version)
}

// 获取插件中的指定成员值
func (this *Plugin) Get(mem string) interface{} {
	if this.symbols == nil {
		return nil
	}
	symbol, ok := this.symbols[mem]
	if ok {
		return symbol
	}
	return nil
}
