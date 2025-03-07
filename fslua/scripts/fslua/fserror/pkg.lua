--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: fserrpr package
@author: fanky
@version: 1.0
@date: 2021-05-03
--]]


return {
	init = function(fsky, tofsky)
		require("fserror.error").init(fsky, tofsky)
	end
}
