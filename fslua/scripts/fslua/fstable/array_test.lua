local Array = require("fsky").Array

local array = Array.new({1,2,3,4, nil, 5, 6})

array.sort(
	function(x, y)
		if x == nil then return -1 > y end
		if y == nil then return x > -1 end
		return x > y
	end
)

print(array)
array.reverse()
print(array)

