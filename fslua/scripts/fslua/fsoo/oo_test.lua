local fsky = require("fsky")


--------------------------------------------------
-- 基类1
--------------------------------------------------
local B1 = fsky.class("B1")

function B1.f_ctor(this, value)
	this.value = value
	print("constructor of B1.")
end

function B1.test1(this)
	print("B1.test1")
end

--------------------------------------------------
-- 基类2
--------------------------------------------------
local B2 = fsky.class("B2")   -- 基类2

function B2.f_ctor(this, value)
	this.value = value
	print("constructor of B2.")
end

function B2.test2(this)
	print("B2.test2")
end

--------------------------------------------------
-- 子类1（继承于B1）
--------------------------------------------------
local C1 = fsky.class("C1", B1).cf_members({  -- 定义类变量
	cvalue = 1000,
})

function C1.f_ctor(this, value)
	B1.f_ctor(this, value)           -- 调用基类构造函数
	print("construstor of C1")
end

function C1.func(this)
	print("C1.func()", this.value)
end

--------------------------------------------------
-- 子类2（继承于B1和B2）
--------------------------------------------------
local C2 = fsky.class("C2", B1, B2)

function C2.f_ctor(this, value)
	B1.f_ctor(this, value)
	B2.f_ctor(this, value + 100)
	this.newValue = 10000
	print("construstor of C2.")
end

function C2.func(this)
	print("C2.func()", this.value, this.newValue)
	this.mfunc = function()end
end

-- 调用基类函数
function C2.callbase(this, text)
	print(text)
	B1.test1(this)
	B2.test2(this)
end


--------------------------------------------------
-- 测试
--------------------------------------------------
local c1 = C1.new(100)
local c2 = C2.new(200)

print('------------------------------')
print("打印类C1: ", C1)
print("打印类C2: ", C2)

print("打印对象c1: ", c1)
print("打印对象c2: ", c2)

print('------------------------------')
print("C1是否继承于B1: ", C1:cf_isfrom(B1))
print("C1是否继承于B2: ", C1:cf_isfrom(B2))
print("C2是否继承于B2: ", C2:cf_isfrom(B2))

print('------------------------------')
print("c1是否是C1或其子类的对象：", c1.f_isa(C1))
print("c1是否是B2或其子类对象：", c1.f_isa(B2))
print("c2是否是B2或其子类对象：", c2.f_isa(B2))

print('------------------------------')
c1.func()			-- func 在 C1 中定义
c2.func()			-- func 在 C2 中定义
c1.test1()			-- test1 在 B1 中定义
c2.test2()			-- test2 在 B2 中定义

print('------------------------------')
print("获取c1属性值：", c1.value)
print("获取c2属性值：", c2.value, c2.newValue)
print("获取c2函数类型属性值：", c2.mfunc)

print('------------------------------')
c2.callbase("下面调用基类函数：")

print('------------------------------')
print("类变量 C1.cvalue = ", C1.cvalue)
C1.cvalue = 2000
print("类变量 C1.cvalue = ", C1.cvalue)
print("如果C1的对象c1没有定义同名变量，则也可以通过对象引用类变量：c1.cvalue = ", c1.cvalue)
c1.cvalue = 5000
print("对象 c1 定义了同名成员变量 cvalue 后：")
print("C1.cvalue = ", c1.m_class.cvalue)
print("c1.cvalue = ", c1.cvalue)
