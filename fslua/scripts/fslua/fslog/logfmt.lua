--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: log formater
@author: fanky
@version: 1.0
@date: 2021-04-24
--]]

local fsstr = require("fsstr/str").str

----------------------------------------------------------------------
-- private
----------------------------------------------------------------------
local function _gettime()
	local clock = os.clock() * 1000000
	clock = tostring(math.floor(clock))
	clock = string.sub(clock, -4, -1)
	return os.date("%Y-%m-%d %H:%M:%S.") .. clock
end

local function _output(layer, prefix, msg, ...)
	local info = debug.getinfo(layer+1)
	local segs = {
		prefix,
		_gettime(),
		" ",
		string.sub(info.source, 2, -1),
		":",
		info.currentline,
		": ",
		tostring(msg)
	}
	for i=1, select('#', ...) do
		table.insert(segs, " " .. tostring(select(i, ...)))
	end
	return table.concat(segs)
end

local function _outputf(layer, prefix, msg, ...)
	local info = debug.getinfo(layer+1)
	local segs = {
		prefix,
		_gettime(),
		" ",
		string.sub(info.source, 2, -1),
		":",
		info.currentline,
		": ",
		string.format(msg, ...)
	}
	return table.concat(segs)
end

function _outtrace(layer)
	layer = layer + 1
	local indents = "  "
	local segs = {"\n", indents, "traceback:"}
	while(true)
	do
		layer = layer + 1
		local info = debug.getinfo(layer)
		if info == nil then break end
		indents = indents .. "  "
	
		table.insert(segs, "\n")
		table.insert(segs, indents)
		table.insert(segs, info.short_src)
		table.insert(segs, ":" .. info.currentline)
		local fname = info.name
		if fname ~= nil then
			table.insert(segs, ": in function '" .. fname .. "'")
		end
	end
	return table.concat(segs)
end

----------------------------------------------------------------------
-- public
----------------------------------------------------------------------
local logfmt = {}

function logfmt.fmt(layer, prefix, msg, ...)
	return _output(layer+1, prefix, msg, ...)
end

function logfmt.fmtf(layer, prefix, msg, ...)
	return _outputf(layer+1, prefix, msg, ...)
end

function logfmt.tracefmt(layer, prefix, msg, ...)
	return _output(layer+2, prefix, msg, ...) .. _outtrace(layer)
end

function logfmt.tracefmtf(layer, prefix, msg, ...)
	return _outputf(layer+2, prefix, msg, ...) .. _outtrace(layer)
end

----------------------------------------------------------------------
-- initialize
----------------------------------------------------------------------
return {
	init = function(fsky, tofsky)
		fsky.logfmt = logfmt
	end,

	logfmt = logfmt,
}

