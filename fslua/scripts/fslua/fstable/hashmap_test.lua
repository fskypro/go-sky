local HashMap = require("fsky").HashMap

local hm = HashMap.new({aa=1, bb=2, dd=nil, cc=3})
hm.set("dd", nil)
local hm2 = HashMap.new({ee=6})
print(hm2..{ff=90} )

print(hm)
hm.delete("dd")
print(hm)

-- 格式化字典 table
print(HashMap.new({
	aa= 100,
	bb = 200,
	cc = 300,
	dd = HashMap.new({
		xxx = "abc",
		yy = 400,
		zz = HashMap.new({
			kk = "EF",
			vv = "GG"
		})
	})
}).fmt(0, ">>"))

