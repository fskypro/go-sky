--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: os utils
@author: fanky
@version: 1.0
@date: 2021-05-03
--]]

return {
	init = function(fsky)
		require("fsos.os").init(fsky)
		require("fsos.path").init(fsky)
	end
}
