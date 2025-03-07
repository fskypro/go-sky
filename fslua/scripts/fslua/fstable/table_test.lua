local fstable = require("fsky").fstable

local list = {
	'nil',
	'300',
	'nil',
	'nil',
	'nil',
	'500',
	'nil'
}

list = fstable.reverse(list)

print(fstable.listout(list))

local dict = {
	aa = 100,
	bb = 200,
	cc = 300,
	[10] = 50,
	dd = {
		xx = "abc",
		yy = 444,
	}
}

print(fstable.dictout(dict))

local tb1 = fstable.update(dict, {cc = 400})
print(fstable.dictout(tb1))

local tb2 = fstable.union(dict, {cc = 400})
print(fstable.dictout(tb2))

-- 格式化字典 table
print(fs_dfmt({
	aa= 100,
	bb = 200,
	cc = 300,
	dd = {
		xxx = "abc",
		yy = 400,
		zz = {
			kk = "EF",
			vv = "GG"
		}
	}
}, 0, ">>"))

