--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: faky utils
@author: fanky
@version: 1.0
@date: 2021-05-12
--]]

local util = {}

-- 判断一个对象是否可被调用
function util.callable(obj)
	if type(obj) == 'function' then
		return true
	end
	local mtable = getmetatable(obj)
	if mtable == nil then
		return false
	end
	return util.callable(mtable.__call)
end

----------------------------------------------------------------------
-- initialize
----------------------------------------------------------------------
return {
	init = function(fsky, tofsky)
		tofsky(util)
	end,

	util = util,
}

