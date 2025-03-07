local fsky = require("fsky")

local A = fsky.class("A").cf_members({
    addr = "北京"
})

function A.f_ctor(this, name, age)
    this.name = name
    this.age = age
end

function A.get(this)
    return "name: " .. this.name .. ", age: " .. this.age .. ", address: " .. A.addr
end

function A.printInfo(this, s)
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
    
    return this.get()
end

obj1 = A.new("xxxx", 100)
obj2 = A.new("yyyy", 200)

