--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: 
@author: fanky
@version: 1.0
@date: 2021-05-04
--]]


return {
	init = function(fsky, tofsky)
		require("fslog.logfmt").init(fsky, tofsky)
		require("fslog.baselog").init(fsky, tofsky)
		require("fslog.dayfilelog").init(fsky, tofsky)
	end
}
