local p = require("fsky").path
local fstable = require("fsky").fstable

print(p.join("abc/", "def", ".txt"))
print(p.join("//abc/", "def", ".txt"))

local info = p.splitext("./root/test.b.a")
print(info.dir, info.fname, info.ext)


print(fstable.listout(p.split("asdf/asdf////sdf/../.././sdf.g")))
print(p.normalize("/root//sdff/sdf/../.././sd.f"))

print(p.filePathExists("./os.lua"))

print(lfs)
