package fslog

/*

fslog.go
	是为方便直接 import 后使用的 log 模块，可以通过 SetLogger 函数来随意设置一个具体的 Logger：
		S_CmdLogger、S_FileLogger、S_NetLogger

baselogger.go
	主要实现 Logger 接口和 Logger 基类
	I_Logger 为 Logger 接口
	BaseLogger 实现了主要的 Logger 基础功能，包括 Logger 格式化；Logger 消息的屏蔽；函数调用链的解释等

cmdlogger.go
	主要实现在控制台打印日志的 Logger

filelogger.go
	主要实现了在指定文件打印日志的 Logger

netlogger.go
	主要实现将日志传输到网络服务器的 Logger
*/
