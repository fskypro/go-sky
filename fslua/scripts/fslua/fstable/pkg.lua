--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 
@author: fanky
@version: 1.0
@date: 2021-05-03
--]]


return {
	init = function(fsky, tofsky)
		require("fstable.array").init(fsky, tofsky)
		require("fstable.hashmap").init(fsky, tofsky)
		require("fstable.table").init(fsky, tofsky)
	end
}
