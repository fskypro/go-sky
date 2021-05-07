/**
@copyright: fantasysky 2016
@brief: 实现 Logger 接口，和 Logger 基类
@author: fanky
@version: 1.0
@date: 2019-01-07
**/

package fslog

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"strings"
	"sync"

	"fsky.pro/fsenv"
)

// -----------------------------------------------------------------------------
// Logger
// -----------------------------------------------------------------------------
// Logger 接口
type I_Logger interface {
	Close()

	// ---------------------------------------------------------------
	// private
	// ---------------------------------------------------------------
	print(depth int, prefix string, vs []interface{}) string
	printf(depth int, prefix string, format string, vs []interface{}) string

	printChain(depth int, prefix string, vs []interface{})
	printChainf(depth int, prefix string, format string, vs []interface{})

	// ---------------------------------------------------------------
	// public
	// ---------------------------------------------------------------
	// GetUnshields 获取打印输出的频道，频道名称全为小写
	GetUnshields() []string

	// SetUnshields 设置打印的频道。
	// 如果有频道不存在则设置失败，并返回 false。
	//	可选频道有(不区分大小写)：
	//	info、warn、debug、error、trace
	SetUnshields(chs ...string) bool

	// Unshield 开启某个指定频道的打印
	// 如果指定的频道不存在，则设置失败，并返回 false。
	//	可选频道有(不区分大小写)：
	//	info、warn、debug、error、trace
	Unshield(ch string) bool

	// Shield 屏蔽某个频道的打印
	// 如果指定的频道不存在，则设置失败，并返回 false。
	//	可选频道有(不区分大小写)：
	//	info、warn、debug、error、trace
	Shield(ch string) bool
}

// -----------------------------------------------------------------------------
// inners
// -----------------------------------------------------------------------------
// 频道信息
type _logChannelInfo struct {
	prefix     string // 频道前缀
	shieldable bool   // 是否可屏蔽
}

// {log 频道标签: _logChannelInfo}
var _logChannels = map[string]*_logChannelInfo{
	"debug": &_logChannelInfo{"[DEBUG]", true}, // 调试频道
	"info":  &_logChannelInfo{"[INFO]", true},  // 信息提示频道
	"warn":  &_logChannelInfo{"[WARN]", true},  // 警告频道
	"error": &_logChannelInfo{"[ERROR]", true}, // 错误频道
	"hack":  &_logChannelInfo{"[HACK]", true},  // 欺诈警告频道
	"trace": &_logChannelInfo{"[TRACE]", true}, // 调用链输出频道

	"panic": &_logChannelInfo{"[PANIC]", false}, // 因重大错误，终止应用频道
	"fatal": &_logChannelInfo{"[FATAL]", false}, // 因重大错误，终止应用，并打印调用链的频道
}

// -------------------------------------------------------------------
// log 格式化参数
const _fmtFlags = log.LstdFlags | log.Llongfile | log.Lmicroseconds

// 打印文件起始深度
const _startDepth = 2

// -------------------------------------------------------------------
// inner methods
// -------------------------------------------------------------------
func _getLogPrefix(name string) string {
	prefix := _logChannels[name].prefix
	if prefix == "" {
		panic(fmt.Sprintf("fslog: log cahnnel name '%s' is not exists!", name))
	}
	return prefix
}

// -----------------------------------------------------------------------------
// package inners
// -----------------------------------------------------------------------------
// 初始化函数，加载模块时调用
func initBaseLogger() {
	maxlen := 0
	for _, info := range _logChannels {
		if len(info.prefix) > maxlen {
			maxlen = len(info.prefix)
		}
	}
	for _, info := range _logChannels {
		info.prefix = fmt.Sprintf(fmt.Sprintf("%%-%ds|", maxlen), info.prefix)
	}
}

// -----------------------------------------------------------------------------
// S_BaseLogger
// -----------------------------------------------------------------------------
// S_BaseLogger
type S_BaseLogger struct {
	sync.Mutex
	logger *log.Logger

	// 开放输出的频道
	opens map[string]interface{}

	// 打印回调
	onBeferPrint func(string)
}

// NewByWriter 以指定输出新建一个 S_BaseLogger
func NewBaseLogger(w io.Writer) *S_BaseLogger {
	pLogger := &S_BaseLogger{
		logger:       log.New(w, "", _fmtFlags),
		opens:        make(map[string]interface{}),
		onBeferPrint: func(string) {},
	}
	for ch, _ := range _logChannels {
		pLogger.opens[ch] = nil
	}
	return pLogger
}

// -------------------------------------------------------------------
// inner methods of S_BaseLogger
// -------------------------------------------------------------------
// 判断指定频道是否开启
func (this *S_BaseLogger) _isOppened(ch string) bool {
	ch = strings.ToLower(ch)
	if _, ok := this.opens[ch]; ok {
		return true
	}
	info, ok := _logChannels[ch]
	if !ok {
		this._errorf("fslog: can't indicate channel %q is open or not, it is not exists.%s", ch, fsenv.Endline)
		return false
	}
	return !info.shieldable
}

func (this *S_BaseLogger) _printStack(depth int) {
	this.logger.SetFlags(0)
	this.logger.SetPrefix("")

	// 打印调用链
	start := _startDepth + depth + 5
	lines := bytes.Split(debug.Stack(), []byte(fsenv.Endline))
	for ln := start; ln < len(lines); ln = ln + 1 {
		this.logger.Printf("\t%s%s", lines[ln], fsenv.Endline)
	}

	this.logger.SetFlags(_fmtFlags)
}

func (this *S_BaseLogger) _errorf(format string, vs ...interface{}) {
	this.printf(1, "error", format, vs)
}

// -------------------------------------------------------------------
// private methods of package
// -------------------------------------------------------------------
func (this *S_BaseLogger) print(depth int, prefix string, vs []interface{}) string {
	msg := fmt.Sprintln(vs...)
	if this._isOppened(prefix) {
		this.onBeferPrint(msg)
		this.Lock()
		this.logger.SetPrefix(_getLogPrefix(prefix))
		this.logger.Output(_startDepth+depth, msg)
		this.Unlock()
	}
	return msg
}

func (this *S_BaseLogger) printf(depth int, prefix string, format string, vs []interface{}) string {
	msg := fmt.Sprintf(format+fsenv.Endline, vs...)
	if this._isOppened(prefix) {
		this.onBeferPrint(msg)
		this.Lock()
		this.logger.SetPrefix(_getLogPrefix(prefix))
		this.logger.Output(_startDepth+depth, msg)
		this.Unlock()
	}
	return msg
}

func (this *S_BaseLogger) printChain(depth int, prefix string, vs []interface{}) {
	if !this._isOppened(prefix) {
		return
	}
	msg := fmt.Sprintln(vs...)
	this.onBeferPrint(msg)
	this.Lock()
	this.logger.SetPrefix(_getLogPrefix(prefix))
	this.logger.Output(_startDepth+depth, msg)
	this._printStack(depth + 1)
	this.Unlock()
}

func (this *S_BaseLogger) printChainf(depth int, prefix string, format string, vs []interface{}) {
	if !this._isOppened(prefix) {
		return
	}
	msg := fmt.Sprintf(format, vs...)
	this.onBeferPrint(msg)
	this.Lock()
	this.logger.SetPrefix(_getLogPrefix(prefix))
	this.logger.Output(_startDepth+depth, msg)
	this._printStack(depth + 1)
	this.Unlock()
}

// -------------------------------------------------------------------
// public methods of S_BaseLogger
// -------------------------------------------------------------------
func (this *S_BaseLogger) Close() {

}

// --------------------------------------------------------
// GetUnshields 获取打印输出的频道，频道名称全为小写
func (this *S_BaseLogger) GetUnshields() []string {
	chs := []string{}
	for ch, _ := range this.opens {
		chs = append(chs, ch)
	}
	return chs
}

// SetUnshields 设置打印的频道。
// 如果有频道不存在则设置失败，并返回 false。
//	可选频道有(不区分大小写)：
//	info、warn、debug、error、trace
func (this *S_BaseLogger) SetUnshields(chs ...string) bool {
	var tmp []string
	for _, ch := range chs {
		ch = strings.ToLower(ch)
		info, ok := _logChannels[ch]
		if !ok {
			this._errorf("fslog: can't set open log channels, channel %q is not exists.%s", ch, fsenv.Endline)
			return false
		}
		if !info.shieldable {
			this._errorf("fslog: can't set open log channels, channel %q can't be shielded!%s", ch, fsenv.Endline)
			return false
		} else {
			tmp = append(tmp, ch)
		}
	}
	this.Lock()
	this.opens = make(map[string]interface{})
	for _, ch := range tmp {
		this.opens[ch] = nil
	}
	this.Unlock()
	return true
}

// Unshield 开启某个指定频道的打印
// 如果指定的频道不存在，则设置失败，并返回 false。
//	可选频道有(不区分大小写)：
//	info、warn、debug、error、trace
func (this *S_BaseLogger) Unshield(ch string) bool {
	ch = strings.ToLower(ch)
	if _, ok := _logChannels[ch]; !ok {
		this._errorf("fslog: can't open channel %q, it is not exists!%s", ch, fsenv.Endline)
		return false
	}
	this.Lock()
	this.opens[ch] = nil
	this.Unlock()
	return true
}

// Shield 屏蔽某个频道的打印
// 如果指定的频道不存在，则设置失败，并返回 false。
//	可选频道有(不区分大小写)：
//	info、warn、debug、error、trace
func (this *S_BaseLogger) Shield(ch string) bool {
	ch = strings.ToLower(ch)
	info, ok := _logChannels[ch]
	if !ok {
		this._errorf("fslog: can't shied channel %q, it is not exists!%s", ch, fsenv.Endline)
		return false
	}
	if !info.shieldable {
		this._errorf("can't shied channel %q, it can't be shielded!%s", ch, fsenv.Endline)
		return false
	}
	this.Lock()
	delete(this.opens, ch)
	this.Unlock()
	return true
}
