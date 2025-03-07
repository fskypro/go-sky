--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: luasky defination
@author: fanky
@version: 1.0
@date: 2021-04-29
--]]


----------------------------------------------------------------------
-- private
----------------------------------------------------------------------
local null = {}
setmetatable(null, null)
function null:__tostring()
	return "fsky.null"
end


----------------------------------------------------------------------
-- initialize
----------------------------------------------------------------------
return {
	init = function(fsky, tofsky)
		fsky.null = null
	end,

	null = null,
}

