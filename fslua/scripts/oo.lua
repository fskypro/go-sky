-- 模拟类对象

-- MyClass.lua
MyClass = {
    addr = "北京"
}
MyClass.__index = MyClass

-- Constructor
function MyClass:new(name, age)
    local obj = setmetatable({}, MyClass)
    obj.name = name or "Unknown"
    obj.age = age or 0
    obj.addr = "深圳"
    return obj
end

function MyClass:get()
    return "name: " .. self.name .. ", age: " .. self.age .. ", address: " .. MyClass.addr
end

-- Member function
function MyClass:printInfo1(s)
    print("获取 go 对象的成员变量：", s.name, s.Age )
    print("调用 go 对象的方法：", s:Func())
    print("调用 go 对象的设置方法：s:SetAge(60)", s:SetAge(60))

    print("获取 go 对象父结构体的成员变量：", s.value)
    print("调用 go 对象父结构体的方法：", s:Test())

    print("获取 go 对象成员对象的成员变量：", s.hand.name)
    print("调用 go 对象成员对象的方法：", s.hand:Func())

    local name, age = s:Func()
    print("Name:", name)
    print("Age:", age)
    
    return self:get()
end

function MyClass.printInfo2(self, s)
    print("获取 go 对象的成员变量：", s.name, s.Age )
    print("调用 go 对象的方法：", s:Func())
    print("调用 go 对象的设置方法：s:SetAge(60)", s:SetAge(60))

    print("获取 go 对象父结构体的成员变量：", s.value)
    print("调用 go 对象父结构体的方法：", s:Test())

    print("获取 go 对象成员对象的成员变量：", s.hand.name)
    print("调用 go 对象成员对象的方法：", s.hand:Func())

    local name, age = s:Func()
    print("Name:", name)
    print("Age:", age)
    
    return self:get()
end

obj1 = MyClass:new("xxx", 100)
obj2 = MyClass:new("yyy", 200)
