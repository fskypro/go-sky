local Error = require("fsky").Error

local err = Error:new("error message.")
print(err)
print(err:message())

local err = Error:new("error: %d", 10000)
print(err:message())

