--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: utils for table
@author: fanky
@version: 1.0
@date: 2021-05-02
--]]

local fstable = {}

-- 判断是否存在指定的值
function fstable.hasvalue(tb, value)
	for _, v in pairs(tb) do
		if v == value then
			return true
		end
	end
	return false
end

------------------------------------------------------------
-- 计算表长度
function fstable.length(tb)
	local len = 0
	for k, _ in pairs(tb) do
		len = len + 1
	end
	return len
end

------------------------------------------------------------
-- 获取 table 的所有 key，以新的 table 返回
function fstable.keysof(tb)
	local keys = {}
	for k, _ in pairs(tb) do
		table.insert(keys, k)
	end
	return keys
end

-- 获取 table 的所有 value，以新的 table 返回
function fstable.valuesof(tb)
	local values = {}
	for _, v in pairs(tb) do
		table.insert(values, v)
	end
	return values
end

------------------------------------------------------------
-- 合并多个 table，并返回一个新的 table
-- 出现在参数后面的 table，其元素将会覆盖前面的
function fstable.update(tb1, ...)
	local newtb = {}
	for k, v in pairs(tb1) do
		newtb[k] = v
	end

	for _, tb in pairs({...}) do
		for k, v in pairs(tb) do
			newtb[k] = v
		end
	end
	return newtb
end

-- 合并多个 table，并返回一个新的 table
-- 如果在参数前面的 table 已经存在某个 key，则后面的 table 的同名 key 将会被丢弃
function fstable.union(tb1, ...)
	local newtb = {}
	local tbs = {...}
	for i = #tbs, 1, -1 do
		for k, v in pairs(tbs[i]) do
			newtb[k] = v
		end
	end
	for k, v in pairs(tb1) do
		newtb[k] = v
	end
	return newtb
end

-- 忽略 key，将所有 table 的 value 合并成一个 table(数组)
function fstable.concat(tb1, ...)
	local newtb = {}
	local index = 1
	for k, v in pairs(tb1) do
		newtb[index] = v
		index = index + 1
	end

	for _, tb in pairs({...}) do
		for k, v in pairs(tb) do
			newtb[index] = v
			index = index + 1
		end
	end
	return newtb
end

------------------------------------------------------------
-- 对数组 table 元素进行翻转
function fstable.reverse(tb)
	local newtb = {}
	local count = #tb
	for i = count, 1, -1 do
		newtb[count-i+1] = tb[i]
	end
	return newtb
end

------------------------------------------------------------
-- 字符串形式列出所有数组表元素
function fstable.listout(tb)
	if type(tb) == 'string' then return '"' .. tb .. '"' end
	if type(tb) ~= 'table' then return tostring(tb) end

	local items = {}
	local idx = 1
	for k, v in pairs(tb) do
		if k > idx then
			for i = idx, k - 1 do
				table.insert(items, "nil")
			end
		end
		if type(v) == 'string' then
			v = '"' .. v .. '"'
		elseif type(v) == 'table' then
			v = fstable.listout(v)
		else
			v = tostring(v)
		end
		table.insert(items, v)
		idx = k + 1
	end
	return '[' .. table.concat(items, ',') .. ']'
end

-- 字符串形式列出所有映射表元素，只展开第一层
function fstable.dictout(tb)
	if type(tb) == 'string' then return '"' .. tb .. '"' end
	if type(tb) ~= 'table' then return tostring(tb) end
	local items = {}
	for k, v in pairs(tb) do
		if type(v) == 'string' then
			v = '"' .. v .. '"'
		elseif type(v) == 'table' then
			v = fstable.dictout(v)
		else
			v = tostring(v)
		end
		table.insert(items, tostring(k) .. '=' .. v)
	end
	return '{' .. table.concat(items, ',') .. '}'
end

-- 格式化列出所有映射表元素
--   deep 为展开深度，默认为 1，只展开第一层 table，如果小等于 0，则全部展开
--   prefix 为所有行前缀，默认为空字符串
--   ident 为嵌套缩进，默认为四个空格
function fstable.dictfmt(tb, deep, prefix, ident)
	if type(tb) == 'string' then return '"' .. tb .. '"' end
	if type(tb) ~= 'table' then return tostring(tb) end

	if type(deep) ~= 'number' then deep = 1 end
	if type(prefix) ~= 'string' then prefix = "" end
	if type(ident) ~= 'string' then ident = "	" end

	local strs = {prefix}
	local function extend(obj, layer)
		if type(obj) == 'table' and (layer < deep or deep <= 0) then
			table.insert(strs, '{\n')
			local left = prefix .. string.rep(ident, layer+1)
			for k, v in pairs(obj) do
				table.insert(strs, left)
				table.insert(strs, tostring(k) .. ' = ')
				extend(v, layer+1)
				table.insert(strs, ',\n')
			end
			left = prefix .. string.rep(ident, layer)
			table.insert(strs, left .. '}')
		elseif type(obj) == 'string' then
			table.insert(strs, '"'..tostring(obj)..'"')
		else
			table.insert(strs, tostring(obj))
		end
	end
	extend(tb, 0)
	return table.concat(strs, '')
end

----------------------------------------------------------------------
-- initialize
----------------------------------------------------------------------
return {
	init = function(fsky, tofsky)
		fsky.fstable = fstable
	end,

	table = fstable,
}
