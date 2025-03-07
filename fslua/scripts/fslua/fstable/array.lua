--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: implement an array class
@author: fanky
@version: 1.0
@date: 2021-04-26
--]]


local define = require("fsdefine")
local oo = require("fsoo.oo")
local Error = require("fserror.error")
local fstable = require("fstable.table").fstable

local null = define.null

local module = {}

-------------------------------------------------------------------------------------
-- Array
-------------------------------------------------------------------------------------
local Array = oo.class("Array")
do
	---------------------------------------------------------------------------
	-- 内部方法
	---------------------------------------------------------------------------
	-- 将越界的索引转换为正索引或者 0
	-- 0 表示索引越界
	function Array._pindex(this, index)
		local count = #this._elems
		if index > count then
			return 0
		end

		if index < 0 then
			index = count + 1 + index
		end
		if index < 1 then
			return 0
		end
		return index
	end

	-- 将越界索引修正为合法索引
	function Array._mindex(this, index)
		local count = #this._elems
		if index > count then
			return count + 1
		end

		if index < 0 then
			index = count + 1 + index
		end
		if index < 0 then
			return 0
		end
		return index
	end

	---------------------------------------------------------------------------
	-- 构造函数
	---------------------------------------------------------------------------
	function Array.f_ctor(this, other)
		this._elems = {}
		if other == nil then
			return
		elseif other.m_class == Array then
			for k, v in pairs(other._elems) do
				table.insert(this._elems, v)
			end
		elseif type(other) == 'table' then
			for i = 1, #other do
				local v = other[i]
				if v == nil then
					table.insert(this._elems, null)
				else
					table.insert(this._elems, v)
				end
			end
		end
	end

	---------------------------------------------------------------------------
	-- 单元素操作
	---------------------------------------------------------------------------
	-- 返回数组长度
	function Array.count(this)
		return #this._elems
	end

	--------------------------------------------------------
	-- 获取指定索引处的元素，如果索引超出范围则返回 fsky.null
	-- 支持负索引，-1 表示最后一个元素
	function Array.get(this, index)
		local idx = this._pindex(index)
		if idx == 0 then return null end
		local value = this._elems[idx]
		if value == null then
			return nil
		end
		return value
	end

	-- 设置索引值，如果设置成功则返回 true，否则索引超出范围的话返回 false
	-- 支持负索引，-1 表示最后一个元素的索引
	function Array.set(this, index, value)
		local idx = this._pindex(index)
		if idx == 0 then
			return false
		end
		if value == nil then
			value = null
		end
		this._elems[idx] = value
		return true
	end

	--------------------------------------------------------
	-- 查看指定元素的索引，如果元素不存在，则返回 -1
	function Array.index(this, value)
		if value == nil then
			value = null
		end
		for i, v in pairs(this._elems) do
			if v == value then
				return i
			end
		end
		return -1
	end

	-- 查看数值中是否存在指定元素
	function Array.hasValue(this, value)
		return this.index(value) > 0
	end

	--------------------------------------------------------
	-- 在数组最前面插入一个元素
	function Array.pushFront(this, value)
		local count = #this._elems
		for i = count, 1, -1 do
			this._elems[i+1] = this._elems[i]
		end
		if value == nil then
			value = null
		end
		this._elems[1] = value
	end

	-- 在数组最后面添加一个元素
	function Array.pushBack(this, value)
		local count = #this._elems
		if value == nil then
			value = null
		end
		this._elems[count+1] = value
	end

	-- 在指定索引前面插入一个元素，支持负索引，-1 表示最后一个元素的索引
	function Array.insert(this, index, value)
		if value == nil then
			value = null
		end

		local count = #this._elems
		if index > count then
			table.insert(this._elems, value)
			return
		end

		local idx = this._mindex(index)
		if idx < 1 then idx = 1 end

		for i = count, idx, -1 do
			this._elems[i+1] = this._elems[i]
		end
		this._elems[idx] = value
	end

	--------------------------------------------------------
	-- 删除指定索引处的元素
	-- 删除成功返回 true，如果索引越界，则返回 false
	function Array.remove(this, index)
		local idx = this._pindex(index)
		if idx == 0 then
			return false
		end
		local count = #this._elems
		for i = idx, count do
			this._elems[i] = this._elems[i+1]
		end
		this._elems[count] = nil
		return true
	end

	-- 删除指定索引区段中的元素
	function Array.removes(this, idx1, idx2)
		idx1 = this._mindex(idx1)
		idx2 = this._mindex(idx2)
		if idx1 > idx2 then
			idx1, idx2 = idx2, idx1
		end
		local count = #this._elems
		if idx2 < 1 or idx1 > count then
			return
		end

		if idx1 < 1 then idx1 = 1 end
		if idx2 > count then idx2 = count end
		for i=idx1, count do
			idx2 = idx2 + 1
			this._elems[i] = this._elems[idx2]
		end
	end

	-- 删除第一个指定的元素
	-- 删除成功返回 true，如果元素不存在，则返回 false
	function Array.removeValue(this, value)
		local index = this.index(value)
		if index < 0 then
			return false
		end
		return this.remove(index)
	end

	-- 删除所有指定的元素
	function Array.removeAllOfValue(this, value)
		if value == nil then
			value = null
		end
		local count = #this._elems
		local idx = 1
		local space = 0
		while(idx <= count) do
			if value == this._elems[idx+space] then
				space = space + 1
			end
			if space > 0 then
				this._elems[idx] = this._elems[idx+space]
			end
			idx = idx + 1
		end
	end

	-- 清空所有元素
	function Array.clear(this)
		this._elems = {}
	end

	--------------------------------------------------------
	-- 追加一组元素
	-- arr 可以是 Array 也可以是数组 table
	function Array.concat(this, arr)
		if type(arr) ~= 'table' then
			return false
		end
		if arr.f_isa and arr.f_isa(Array) then
			for _, v in pairs(arr._elems) do
				table.insert(this._elems, v)
			end
			return true
		end
		local v
		for i = 0, #arr do
			v = arr[i]
			if v == nil then
				v = null
			end
			table.insert(this._elems, v)
		end
		return true
	end

	-- 追加一组参数中的元素
	function Array.adds(this, ...)
		for i = 1, select('#', ...) do
			local value = select(i, ...)
			if value == nil then
				value = null
			end
			table.insert(this._elems, value)
		end
	end

	--------------------------------------------------------
	-- 获取指定区段元素（包含 idx2 的元素）
	function Array.slice(this, idx1, idx2)
		idx1 = this._mindex(idx1)
		idx2 = this._mindex(idx2)
		if idx1 > idx2 then
			idx1, idx2 = idx2, idx1
		end
		local arr = Array.new()
		for i = idx1, idx2 do
			table.insert(arr._elems, this._elems[i])
		end
		return arr
	end

	-- 所有元素迭代器，skip 参数指定迭代步长
	-- 注意：
	--   1、如果 idx1、idx2、skip 都不传入，则顺序迭代全部元素
	--   2、如果 idx1 小于 idx2，则进行正向迭代，并且要求 skip 必须大于 0，否则只迭代第一个元素
	--   3、如果 idx1 大于 idx2，则进行反向迭代，并且要求 skip 必须小于 0，否则只迭代第一个元素
	--   4、支持负索引，-1 表示最后一个元素
	function Array.iter(this, idx1, idx2, skip)
		local count = #this._elems
		if type(idx1) ~= 'number' then idx1 = 1 end
		if type(idx2) ~= 'number' then idx2 = count end
		idx1 = this._mindex(idx1)
		idx2 = this._mindex(idx2)

		-- 返回只迭代一个元素的迭代器
		local iterone = function(idx)
			local itered = false
			return function()
				if itered then return end
				itered = true
				local value = this._elems[idx]
				if value == null then value = nil end
				return idx, value
			end
		end

		-- 迭代一个元素
		if idx1 == idx2 then
			if idx1 < 1 or idx1 > count then
				return function()end
			end
			local itered = false
			return iterone(idx1)
		end

		-- 正向迭代
		if idx1 < idx2 then
			if idx2 < 1 or idx1 > count then
				return function() end
			end
			if type(skip) ~= 'number' then skip = 1 end
			if idx1 < 1 then idx1 = 1 end
			if idx2 > count then idx2 = count end

			-- 正向迭代要求 skip 必须大于 0，否则只迭代第一个元素
			if skip <= 0 then return iterone(idx1) end

			local value
			local index = idx1
			return function()
				if index <= idx2 then
					value = this._elems[index]
					if value == null then value = nil end
					local newIndex = index
					index = index + skip
					return newIndex, value
				end
			end
		end

		-- 反向迭代
		if idx1 < 1 or idx2 > count then
			return function() end
		end
		if type(skip) ~= 'number' then skip = -1 end
		if idx1 > count then idx1 = count end
		if idx2 < 1 then idx2 = 1 end

		-- 反向迭代要求 skip 必须小于 0，否则只迭代第一个元素
		if skip >= 0 then return iterone(idx1) end

		local value
		local index = idx1
		return function()
			if index >= idx2 then
				value = this._elems[index]
				if value == null then value = nil end
				local newIndex = index
				index = index + skip
				return newIndex, value
			end
		end
	end

	--------------------------------------------------------
	-- 翻转顺序
	function Array.reverse(this)
		this._elems = fstable.reverse(this._elems)
	end

	-- 排序，如果忽略参数 func，则默认从小到大排序
	-- func 包含两个参数，表示冒泡排序中相邻的两个元素
	function Array.sort(this, func)
		table.sort(this._elems, function(x, y)
			if x == null then x = nil end
			if y == null then y = nil end
			return func(x, y)
		end)
	end

	--------------------------------------------------------
	-- 转换为 table
	-- 注意：如果数组最后的元素为 nil，则转换为 table 后，所有后面的 nil 元素将会丢失
	function Array.toTable()
		local tb = {}
		for _, v in pairs(this._elems) do
			if v == null then
				v = nil
			end
			table.insert(tb, v)
		end
		return tb
	end

	---------------------------------------------------------------------------
	-- 元方法
	---------------------------------------------------------------------------
	-- 合并两个数组元素，并返回一个新的数组
	-- arr 可以是数组 table
	function Array.__add(this, arr)
		local newarr = Array.new()
		for _, v in pairs(this._elems) do
			table.insert(newarr._elems, v)
		end
		if type(arr) ~= 'table' then
			return newarr
		end
		if arr.f_isa and arr.f_isa(Array) then
			for _, v in pairs(arr._elems) do
				table.insert(newarr._elems, v)
			end
			return newarr
		end
		local value
		for i = 1, #arr do
			value = arr[i]
			if value == nil then
				value = null
			end
			table.insert(newarr._elems, value)
		end
		return newarr
	end

	-- 合并两个数组元素，并返回一个新的数组，等同 +
	function Array.__concat(this, arr)
		return Array.__add(this, arr)
	end

	function Array.__tostring(this)
		local elems = {}
		for i = 1, #(this._elems) do
			local value = this._elems[i]
			if value == null then
				table.insert(elems, 'nil')
			elseif type(value) == 'string' then
				table.insert(elems, '"'..value..'"')
			else
				table.insert(elems, tostring(value))
			end
		end
		return 'Array[' .. table.concat(elems, ',') .. ']'
	end
end

return {
	init = function(fsky, tofsky)
		fsky.Array = Array
	end,

	Array = Array,
}

