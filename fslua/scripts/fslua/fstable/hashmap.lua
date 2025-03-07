--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: hash table
@author: fanky
@version: 1.0
@date: 2021-05-02
--]]

local oo = require("fsoo.oo")
local null = require("fsdefine").null
local fsutil = require("fsutil").util
local Array = require("fstable.array").Array

local HashMap = oo.class("HashMap")
do
	local function isHashMap(obj)
		if type(obj) ~= 'table' then
			return false
		end
		if not fsutil.callable(obj.f_isa) then
			return false
		end
		return obj.f_isa(HashMap)
	end

	--------------------------------------------------------
	-- public
	--------------------------------------------------------
	function HashMap.f_ctor(this, other)
		this._members = {}
		this.update(other)
	end

	-- 返回元素个数
	function HashMap.count(this)
		local count = 0
		for _ in pairs(this._members) do
			count = count + 1
		end
		return count
	end

	-- 用另一个 hashmap 或者 table 更新内部元素
	function HashMap.update(this, other)
		if other == nil or other == this
		or type(other) ~= 'table' then
			return
		end
		if isHashMap(other) then
			for k, v in pairs(other._members) do
				this._members[k] = v
			end
			return
		end
		for k, v in pairs(other) do
			if v == nil then
				v = null
			end
			this._members[k] = v
		end
	end

	-- 获取指定 key 的元素
	-- 如果 key 不存在，则返回 fsky.null
	function HashMap.get(this, key)
		local value = this._members[key]
		if value == nil then
			return nil
		end
		if value == null then
			return nil
		end
		return value
	end

	-- 设置指定 key 的值，设置成功返回 true；如果 key 不存在，则返回 false
	function HashMap.set(this, key, value)
		if this._members[key] == nil then
			return false
		end
		if value == nil then
			value = null
		end
		this._members[key] = value
		return true
	end

	-- 添加一个键值对，如果键已经存在，则更新值
	function HashMap.add(this, key, value)
		if value == nil then
			value = null
		end
		this._members[key] = value
	end

	--------------------------------------------------------
	-- 判断指定 key 是否在 hashmap 中
	function HashMap.hasKey(this, key)
		return this._members[key] ~= nil
	end

	-- 查找是否存在指定 value 的 item
	function HashMap.hasItem(this, value)
		if value == nil then value = null end
		for k, v in pairs(this._members) do
			if value == v then
				return true
			end
		end
		return false
	end

	--------------------------------------------------------
	-- 获取所有 key
	-- 返回 fsky.Array 对象
	function HashMap.keys(this)
		local keys = Array.new()
		for k, _ in pairs(this._members) do
			keys.pushBack(k)
		end
		return keys
	end

	-- 获取所有 value
	-- 返回 fsky.Array 对象
	function HashMap.values(this)
		local values = Array.new()
		for _, v in pairs(this._members) do
			if v == null then v = nil end
			values.pushBack(v)
		end
		return values
	end

	-- 返回一个新的元素 table
	function HashMap.toTable(this)
		local tb = {}
		for k, v in pairs(this._members) do
			if v == null then v = nil end
			tb[k] = v
		end
		return tb
	end

	--------------------------------------------------------
	-- 返回所有元素的迭代器
	function HashMap.iterall(this)
		k, v = next, this._members, nil
		if v == null then value = nil end
		return k, v
	end

	-- 操作函数遍历所有元素
	-- func 包含两个参数：key，value
	-- 如果 func 返回 true，则结束循环
	function HashMap.travel(this, func)
		for k, v in pairs(this._members) do
			if v == null then v = nil end
			if func(k, v) then
				break
			end
		end
	end

	--------------------------------------------------------
	-- 删除指定 key 的元素
	-- 如果 key 存在，则返回 true，否则返回 false
	function HashMap.delete(this, key)
		if key == nil then return false end
		if this._members[key] == nil then
			return false
		end
		this._members[key] = nil
		return true
	end

	-- 清除所有元素
	function HashMap.clear(this)
		this._members = {}
	end

	--------------------------------------------------------
	-- 格式化输出
	--   deep 为展开深度，默认为 1，只展开第一层 table，如果小等于 0，则全部展开
	--   prefix 为所有行前缀，默认为空字符串
	--   ident 为嵌套缩进，默认为四个空格
	-- 注意：只对 HashMap 展开，不对 table 展开
	function HashMap.fmt(this, deep, prefix, ident)
		if type(deep) ~= 'number' then deep = 1 end
		if type(prefix) ~= 'string' then prefix = "" end
		if type(ident) ~= 'string' then ident = "    " end

		local strs = {prefix}
		local function extend(obj, layer)
			if isHashMap(obj) and (layer < deep or deep <= 0) then
				table.insert(strs, this.m_class.cm_name .. '{\n')
				local left = prefix .. string.rep(ident, layer+1)
				for k, v in pairs(obj._members) do
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
		print("aaaa", this.count())
		extend(this, 0)
		return table.concat(strs, '')
	end

	--------------------------------------------------------
	-- 元方法
	--------------------------------------------------------
	-- 返回一个新的 hashmap 对象，该对象由 other 更新 this 所得
	-- other 可以是一个 table，也可以是另一个 hashmap
	function HashMap.__add(this, other)
		local hm = HashMap.new()
		for k, v in pairs(this._members) do
			hm._members[k] = v
		end

		if type(other) ~= 'table' then 
			return hm
		end
		if isHashMap(other) then
			for k, v in pairs(other._members) do
				hm._members[k] = v
			end
			return hm
		end

		for k, v in pairs(other) do
			if v == nil then v = null end
			hm._members[k] = v
		end
		return hm
	end
	
	function HashMap.__concat(this, other)
		return this.__add(other)
	end

	function HashMap.__tostring(this)
		local items = {}
		for k, v in pairs(this._members) do
			if v == null then v = nil end
			table.insert(items, tostring(k) .. "=" .. tostring(v))
		end
		return this.m_class.cm_name .. "{" .. table.concat(items, ",") .. "}"
	end
end


return {
	init = function(fsky, tofsky)
		fsky.HashMap = HashMap
	end,

	HashMap = HashMap,
}
