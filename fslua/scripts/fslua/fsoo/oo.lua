--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: class maker
@author: fanky
@version: 1.0
@date: 2021-05-03
--]]


----------------------------------------------------------------------
-- 使用说明
----------------------------------------------------------------------
-- 一、通过 caass 函数创建的类包含以下默认类成员和类方法：
--   cm_name   ：类名称
--   cm_bases  ：基类列表
--   cm_meta   ：meta table
--   cf_members：添加静态成员，通常这样写：C = class("C", Base):cf_members({cvalue = 100})
--   cf_isfrom ：判断参数指定的类是否是自己的基类
--
--   注意：
--       用户避免定义以上名称的类变量
--
-- 二、通过 new 创建的类对象包含以下默认成员和方法：
--   m_class ：对象所属的类
--   m_ptr   ：对象地址（16进制形式字符串）
--   f_isa   ：判断对象所属类是否是参数中指定类的子类，如：inst:f_isa(Object) == true
--   f_isinst：判断对象所属类是否是参数中指定的类
--
--   注意：
--      用户避免定义以上名称的成员变量
--
-- 三、创建的类对象，如果有以下构造函数，则调用所属类的 new 时，该构造函数会被调用
--   f_ctor  ：构造函数，构造函数的参数与类的 new 方法参数一致，注意：
--             对象构造完成时，对象的 f_cotr 会被调用，但是其基类的构造函数并不会被调用。
--             如果用户意图要调用基类构造函数，则需要自己主动显式调用
--
-- 四、注意：
--    1、class 支持多重继承，多个父类成员发生重叠时，class 前面的参数类成员会覆盖后面的
--    2、所有类成员方法都默认带有第一个参数，该参数通常命名为 this，用于表示对象本身
--    3、类可以显式声明构造函数 f_ctor，在 f_ctor 中初始化成员变量，如果用户没有定义 f_ctor
--       则，该类会有一个默认构造函数。
--    4、可以通过类名定义和设置类变量，类变量所有对象共享
--    5、全部方法都用点“.”引用，不要用冒号“:”
--
-- 五、示例：
--    参阅：./oo_test.lua
----------------------------------------------------------------------

----------------------------------------------------------------------
-- private
----------------------------------------------------------------------
-- 获取对象地址
local function _tableaddr(tb)
	if type(tb) ~= 'table' then 
		return "" 
	end
	local str = tostring(tb)
	local site
	_, site = string.find(str, ' ')
	if site == nil then
		return ""
	end
	return string.sub(str, site+1)
end

-- 判断 cls 是否是 base 子类
function _isinherit(cls, base)
	if type(base) ~= 'table' then
		return false
	end
	local bases = cls.cm_bases
	if bases == nil then
		return false
	end
	for _, b in pairs(bases) do
		if b == base then
			return true
		end
	end
	for _, b in pairs(bases) do
		if _isinherit(b, base) then
			return true
		end
	end
	return false
end

----------------------------------------------------------------------
-- method
-- 成员方法封装
----------------------------------------------------------------------
local function _createmethod(obj, key, func)
	local method = {this = obj, name = key, func = func}
	local addr = _tableaddr(method)
	setmetatable(method, {
		__call = function(a, ...)
			return method.func(method.this, ...)
		end,

		__tostring = function()
			return string.format("class(%s) inst method: %s", obj.m_class.cm_name, addr)
		end
	})
	return method
end

----------------------------------------------------------------------
-- Object class
-- 最顶层基类
----------------------------------------------------------------------
local Object = {
	cm_name = "Object",
	cm_bases = {},
}

local _objaddr = _tableaddr(Object)
function Object.__tostring(this)
	return string.format("class(%s): %s", this.cm_name, _objaddr)
end
setmetatable(Object, Object)

-- 是否继承于指定的类
function Object.cf_isfrom(this, base)
	return false
end

function Object.f_ctor(this)
end

----------------------------------------------------------------------
-- implement class simulator
----------------------------------------------------------------------
local function class(name, ...)
	assert(type(name) == 'string', "make class error: class name must be a string.")
	assert(#name > 0, "make class error: class name must be a not empty string.")

	local cls = {cm_name = name, cm_bases = {}, cm_meta = {}}

	--------------------------------------------------------
	-- 验证参数必须是 Object 或 Object 的子类
	--------------------------------------------------------
	for i = 1, select('#', ...) do
		local c = select(i, ...)
		if c == nil then 
			goto continue
		end
		if c == Object then
			goto continue
		end

		local isfrom = c.cf_isfrom
		assert(type(isfrom) == 'function' and c:cf_isfrom(Object),
			"make class error: base class must inherit from class Object.")
		table.insert(cls.cm_bases, c)

		::continue::
	end
	if #cls.cm_bases == 0 then
		cls.cm_bases = {Object}
	end

	--------------------------------------------------------
	-- 添加类方法
	--------------------------------------------------------
	local clsmeta = cls.cm_meta
	-- 添加类方法 cf_isfrom，用于判断是否继承于指定的基类
	function clsmeta:cf_isfrom(base)
		return _isinherit(self, base)
	end

	-- 添加类方法 __tostring
	function clsmeta:__tostring()
		return string.format("class(%s)", self.cm_name)
	end
	setmetatable(cls, clsmeta)
	clsmeta.__index = clsmeta

	-- 添加类成员函数
	function clsmeta.cf_members(members)
		for k, v in pairs(members) do
			cls.cm_meta[k] = v
		end
		return cls
	end

	--------------------------------------------------------
	-- 类对象创建函数实现
	--------------------------------------------------------
	-- 新建对象函数
	function cls.new(...)
		--assert(self ~= nil, "new class Object error: you should use ':' to call new method, but no '.'.")

		local inst = {m_class = cls}
		inst.m_ptr = _tableaddr(inst)

		setmetatable(inst, cls)

		-- 调用对象构造函数
		inst.f_ctor(...)

		return inst
	end

	--------------------------------------------------------
	-- 添加对象方法
	--------------------------------------------------------
	-- 搜索类属性
	getclsprop = function(c, key)
		for k, v in pairs(c) do
			if k == key then
				return v
			end
		end
		if not c.cm_bases then
			return nil
		end
		for _, base in pairs(c.cm_bases) do
			local elem = getclsprop(base, key)
			if elem ~= nil then
				return elem
			end
		end
		return nil
	end

	---------------------------------------------
	-- 判断指定参数是否是对象所属的类或者所属类的基类
	function cls.f_isa(this, c)
		if this.m_class == c then
			return true
		end
		return _isinherit(this.m_class, c)
	end

	-- 判断对象的所属类是否是 c
	function cls.f_isinst(this, c)
		return this.m_class == c
	end

	-- 索引值获取
	function cls.__index(obj, key)
		local prop = getclsprop(cls, key)
		if type(prop) == 'function' then
			return _createmethod(obj, key, prop)
		end
		return prop
	end

	-- 对象展现
	function cls.__tostring(obj)
		return string.format("class(%s) inst: %s", obj.m_class.cm_name, obj.m_ptr)
	end

	return cls
end

----------------------------------------------------------------------
-- initalize
----------------------------------------------------------------------
return {
	init = function(fsky)
		fsky.class = class
		fsky.Object = Object
	end,

	Object = Object,
	class = class,
}

