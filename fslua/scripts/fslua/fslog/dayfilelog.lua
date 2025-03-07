--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: log printer for a file every day
@author: fanky
@version: 1.0
@date: 2021-04-30
--]]


local class = require("fsoo.oo").class
local fsos = require("fsos.os").os
local fspath = require("fsos.path").path
local fsutil = require("fsutil").util

local logfmt = require("fslog.logfmt").logfmt
local BaseLog = require("fslog.baselog").BaseLog


--------------------------------------------------------------------------------
-- NewLogCmd
--------------------------------------------------------------------------------
local NewLogCmd = class("NewLogCmd")
do
	function NewLogCmd.f_ctor(this, cmd, ...)
		if type(cmd) ~= 'string' then
			this._cmd = ""
			this._args = {}
		else
			this._cmd = cmd
			this._args = {...}
		end
	end

	function NewLogCmd.exec(this, logger, logPath)
		if this._cmd == "" then
			return
		end
		local cmd = string.format("%s %q", this._cmd, logPath)
		if #this._args > 0 then
			cmd = cmd .. " '" .. table.concat(this._args, "' '") .. "'"
		end

		local ret, err = io.popen(cmd)
		if err ~= nil then
			logger.error("execute new log file command fail: ", err)
			return
		end
		logger.infof(
			"excute new-log-file command(%s). output:%s%s",
			cmd,
			fsos.newline(),
			ret:read("*all")
		)
	end
end

--------------------------------------------------------------------------------
-- DayfileLog
--------------------------------------------------------------------------------
local DayFileLog = class("DayFileLog", BaseLog)
do
	------------------------------------------------------------------
	-- private
	------------------------------------------------------------------
	-- 创建 log 文件
	function DayFileLog._newLogFile(this)
		if this._file then
			this._file:close()
		end

		local prefix = this._prefix
		if prefix ~= "" then
			prefix = prefix .. "_"
		end
		local today = os.date('%Y-%m-%d')
		local fileName = prefix .. today .. this._logext
		local path = fspath.join(this._logdir, fileName)
		local err
		if fspath.filePathExists(path) then
			this._file, err = io.open(path, 'a')
		else
			this._file, err = io.open(path, 'w')
		end

		if this._file ~= nil then
			local sp = string.rep('-', 50)
			local now = os.date('%H:%M:%S')
			this._file:write(string.format("%s %s %s%s", sp, now, sp, fsos.newline()))
			this._lastday = os.date("%Y%m%d")
			this._logPath = path
			this._newLogCmd.exec(this, path)
		else
			print("create logfile fail:\n\t", err)
			this._logPath = ""
		end
		if this._newLogCB then
			this._newLogCB(path, err) 
		end
	end

	function DayFileLog._output(this, msg)
		if not this._inited then
			print("warn: dayfilelog hasn't intialized.")
			print(msg)
			return
		end

		if os.date("%Y%m%d") ~= this._lastday then
			this._newLogFile()
		end
		if this._file then
			this._file:write(msg .. fsos.newline())
			this._file:flush()
		else
			print(msg)
		end
	end

	------------------------------------------------------------------
	-- private
	------------------------------------------------------------------
	function DayFileLog.f_ctor(this)
		BaseLog.f_ctor(this)
		this._logPath = ""					-- log 文件路径
		this._file = nil					-- log 文件流
		this._newLogCmd = NewLogCmd.new()	-- 新建 log 文件时，触发的命令
		this._newLogCB = nil				-- 新建 log 文件回调函数
		this._inited = false

		this.setOutputHandler(this._output)
	end

	-- 设置 log 属性
	-- prefix 为 log 文件前缀，省略为没有前缀
	-- logdir 为 log 输出文件夹，省略为 ./logs
	-- ext 为 log 文件扩展名，可省略，省略后，默认为 .log
	function DayFileLog.init(this, prefix, logdir, ext)
		if type(prefix) ~= 'string' then prefix = "" end
		if type(logdir) ~= 'string' then logdir = "./logs" end
		if type(ext) ~= 'string' then ext = ".log" end
		if string.sub(ext, 1, 1) ~= '.' then ext = '.' .. ext end
		this._prefix = prefix				-- log 文件前缀
		this._logdir = logdir				-- log 文件所在目录
		this._logext = ext					-- log 文件扩展名

		this._newLogFile()
		this._inited = true
	end

	-- 设置一个可执行命令，每当新建一个 log 文件时，该命令将会被执行
	-- 可以通过不定参数设置多个命令行参数，但是第一个参数肯定是新 log 文件路径
	function DayFileLog.setNewLogCmd(this, cmd, ...)
		this._newLogCmd = NewLogCmd.new(cmd, ...)
		if this._logPath ~= "" then
			this._newLogCmd.exec(this, this._logPath)
		end
	end

	-- 设置新建 log 文件回调
	-- cb 必须是可调用对象，或者 nil，如果是可调用对象，则带两个参数：
	--   path：新建 log 文件的路径
	--   err ：wei nil，则创建 log 成功，为非空字符串，则创建新 log 文件失败，并指示失败原因
	function DayFileLog.setNewLogCallback(this, cb)
		if cb == nil then
			this._newLogCB = nil
			return
		end
		assert(fsutil.callable(cb), "new-log-file callback must be callable.")
		this._newLogCB = cb
	end

	------------------------------------------------------------
	function DayFileLog.close(this)
		if this._file then
			this._file:close()
		end
	end
end

--------------------------------------------------------------------------------
-- initialize
--------------------------------------------------------------------------------
return {
	init = function(fsky, tofsky)
		fsky.DFLog = DayFileLog
		fsky.gDFLog = DayFileLog.new()
	end,

	DFLog = DayFileLog,
}
