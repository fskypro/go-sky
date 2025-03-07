local str = require("fsky").str
local fstable = require("fsky").fstable

print(str.lfill("xxx", 11, ">>|"))
print(str.rfill("xxx", 11, ">>|"))

print(str.startswith("123456", "123"))
print(str.endswith("123456", "56"))

local items = str.split(",123,456,789,,aa", ',,')
print(fstable.listout(items))

items = str.split(",123,456,789,,aa", ',')
print(fstable.listout(items))

print(str.trim(" \t abcdef  "))
print(str.trim("   \t  \n"))
