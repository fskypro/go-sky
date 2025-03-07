--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: base logger
@author: fanky
@version: 1.0
@date: 2021-05-12
--]]


local class = require("fsoo.oo").class
local util = require("fsutil").util
local Array = require("fstable.array").Array
local HashMap = require("fstable.hashmap").HashMap
local logfmt = require("fslog.logfmt").logfmt

----------------------------------------------------------------------
-- private
----------------------------------------------------------------------
local logtypes = HashMap.new({
	debug = "[DEBUG]|",
	error = "[ERROR]|",
	info  = "[INFO] |",
	warn  = "[WARN] |",
	hack  = "[HACK] |",
	trace = "[TRACE]|",
})

----------------------------------------------------------------------
-- BaseLog
----------------------------------------------------------------------
local BaseLog = class("BaseLog")
do
	function BaseLog.f_ctor(this)
		this._outtypes = logtypes.keys()
		this._outputHandler = print
	end

	-- 所有支持的 log 类型
	function BaseLog.logTypes()
		return logtypes.keys()
	end

	-- 设置为输出所有类型
	function BaseLog.outputAll()
		this._outtypes = logtypes.keys()
	end

	-- 屏蔽输出指定类型 log
	function BaseLog.shieldTypes(this, ...)
		for _, t in pairs({...}) do
			local tt = string.lower(t)
			assert(logtypes.hasKey(tt), string.format("log type '%s' is not supported.", t))
			this._outtypes.removeValue(tt)
		end
	end

	----------------------------------------------
	-- 设置 log 输出处理器
	-- handler 为只有一个参数的函数，参数为要输出的 log 消息
	function BaseLog.setOutputHandler(this, handler)
		if not util.callable(handler) then
			this._outputHandler = print
		else
			this._outputHandler = handler
		end
	end

	----------------------------------------------
	-- 添加 log 类别
	function BaseLog.addLogType(this, key, prefix)
		logtypes.add(key, prefix)
	end

	--------------------------------------------------------
	-- 以空格分隔参数，直白返回 log 消息
	--------------------------------------------------------
	function BaseLog.debug(this, msg, ...)
		if not this._outtypes.hasValue("debug") then return end
		msg = logfmt.fmt(1, logtypes.get("debug"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.error(this, msg, ...)
		if not this._outtypes.hasValue("error") then return end
		msg = logfmt.fmt(1, logtypes.get("error"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.info(this, msg, ...)
		if not this._outtypes.hasValue("info") then return end
		msg = logfmt.fmt(1, logtypes.get("info"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.warn(this, msg, ...)
		if not this._outtypes.hasValue("warn") then return end
		msg = logfmt.fmt(1, logtypes.get("warn"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.hack(this, msg, ...)
		if not this._outtypes.hasValue("hack") then return end
		msg = logfmt.fmt(1, logtypes.get("hack"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.trace(this, msg, ...)
		if not this._outtypes.hasValue("trace") then return end
		msg = logfmt.tracefmt(1, logtypes.get("trace"), msg, ...)
		this._outputHandler(msg)
	end

	--------------------------------------------------------
	-- 式化返回 log 消息
	--------------------------------------------------------
	function BaseLog.debugf(this, msg, ...)
		if not this._outtypes.hasValue("debug") then return end
		msg = logfmt.fmtf(1, logtypes.get("debug"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.errorf(this, msg, ...)
		if not this._outtypes.hasValue("error") then return end
		msg = logfmt.fmtf(1, logtypes.get("error"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.infof(this, msg, ...)
		if not this._outtypes.hasValue("info") then return end
		msg = logfmt.fmtf(1, logtypes.get("info"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.warnf(this, msg, ...)
		if not this._outtypes.hasValue("warn") then return end
		msg = logfmt.fmtf(1, logtypes.get("warn"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.hackf(this, msg, ...)
		if not this._outtypes.hasValue("hack") then return end
		msg = logfmt.fmtf(1, logtypes.get("hack"), msg, ...)
		this._outputHandler(msg)
	end

	function BaseLog.tracef(this, msg, ...)
		if not this._outtypes.hasValue("trace") then return end
		msg = logfmt.tracefmtf(1, logtypes.get("trace"), msg, ...)
		this._outputHandler(msg)
	end
end

--------------------------------------------------------------------------------
-- initialize
--------------------------------------------------------------------------------
return {
	init = function(fsky, tofsky)
		fsky.BaseLog = BaseLog
	end,

	BaseLog = BaseLog,
}
