--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: os utils
@author: fanky
@version: 1.0
@date: 2021-05-03
--]]


local fstable = require("fstable.table").table


local os = {}

os.systems = {
	KNOW    = 0,
	LINUX   = 1,
	WINDOWS = 2,
	DARWIN  = 3,
	ANDROID = 4,
	IOS     = 5,
}

local _system = os.systems.LINUX

function os.system()
	return _system
end

-- 设置系统类型
function os.resetSystem(sys)
	if type(sys) ~= 'number' then
		return false
	end
	if os.systems[sys] == nil then
		return false
	end
	_system = sys
end

-- 获取路径分隔符
function os.pathSplitter()
	local splitters = {
		[os.systems.KNOW]    = '/',
		[os.systems.LINUX]   = '/',
		[os.systems.WINDOWS] = '\\',
		[os.systems.DARWIN]  = '/',
		[os.systems.ANDROID] = '/',
		[os.systems.IOS]     = '/',
	}
	return splitters[_system]
end

-- 获取换行符
function os.newline()
	local newlines = {
		[os.systems.KNOW]    = '\n',
		[os.systems.LINUX]   = '\n',
		[os.systems.WINDOWS] = '\r\n',
		[os.systems.DARWIN]  = '\n',
		[os.systems.ANDROID] = '\n',
		[os.systems.IOS]     = '\n',
	}
	return newlines[_system]
end

----------------------------------------------------------------------
-- initialize
----------------------------------------------------------------------
return {
	init = function(fsky)
		fsky.os = os
	end,

	os = os
}
