--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: path utils
@author: fanky
@version: 1.0
@date: 2021-05-03
--]]

local os = require("fsos.os").os
local fsstr = require("fsstr.str").str
local fstable = require("fstable.table").table

----------------------------------------------------------------------
-- public
----------------------------------------------------------------------
-- 剪去路径前后的路径分隔符
local function _cutSplitter(pstr)
	local sp = os.pathSplitter()
	while(string.sub(pstr, 1, 1) == sp) do
		pstr = string.sub(pstr, 2)
	end
	while(string.sub(pstr, -1) == sp) do
		pstr = string.sub(pstr, 1, -2)
	end
	return pstr
end

-- 剪去路径右边的分隔符
local function _cutRightSplitter(pstr)
	local sp = os.pathSplitter()
	while(string.sub(pstr, -1) == sp) do
		pstr = string.sub(pstr, 1, -2)
	end
	return pstr
end

----------------------------------------------------------------------
-- public
----------------------------------------------------------------------
local path = {}

-- 以路径分隔符拆分路径为一组文件夹
function path.split(path)
	local sp = os.pathSplitter()
	path = string.gsub(path, sp .. "+", sp)
	return fsstr.split(path, sp)
end

-- 拼接路径
function path.join(root, ...)
	root = _cutRightSplitter(root)
	local dirs = {root}
	local sp = os.pathSplitter()
	for _, seg in pairs({...}) do
		table.insert(dirs, _cutSplitter(seg))
	end
	return table.concat(dirs, sp)
end

-- 拆分路径中的目录、文件名、扩展名(带.)
-- 返回：{dir=目录, fname=文件名, ext=扩展名}
function path.splitext(path)
	local sp = os.pathSplitter()
	local info = {dir="", fname="", ext=""}
	local index = #path
	local str = ""
	local chr
	while(index > 0) do
		chr = string.sub(path, index, index)
		if chr == '.' and info.ext == '' then
			info.ext = '.' .. str
			str = ""
		elseif chr == sp then
			info.fname = str
			goto L
		else
			str = chr .. str
		end
		index = index - 1
	end
	::L::
	info.dir = string.sub(path, 1, index)
	return info
end

-- 最简化路径
function path.normalize(p)
	local dirs = path.split(p)
	local newDirs = {}
	local count = #dirs
	for i = 1, count do
		local dir = dirs[i]
		if dir == '..' then
			local last = newDirs[#newDirs]
			if last == nil or last == '..' then
				table.insert(newDirs)
			else
				newDirs[#newDirs] = nil
			end
		elseif dir ~= '.' then
			table.insert(newDirs, dir)
		end
	end
	return table.concat(newDirs, os.pathSplitter())
end

------------------------------------------------------------
-- 判断文件或路径是否存在
function path.filePathExists(p)
	local f, e = io.open(p, 'rb')
	if e == nil then
		f:close()
		return true
	end
	return false
end

----------------------------------------------------------------------
-- initialize
----------------------------------------------------------------------
return {
	init = function(fsky)
		fsky.path = path
	end,

	path = path
}

