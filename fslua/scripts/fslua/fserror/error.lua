--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: error base
@author: fanky
@version: 1.0
@date: 2021-04-29
--]]


local oo = require("fsoo.oo")

local Error = oo.class("fsky.Error")

function Error.f_ctor(this, msg, ...)
	local count = select('#', ...)
	if count > 0 then
		this._message = string.format(msg, ...)
	else
		this._message = msg
	end
end

function Error.message(this)
	return this._message
end

----------------------------------------------------------------------
-- load
----------------------------------------------------------------------
return {
	init = function(fsky, tofsky)
		fsky.Error = Error
	end,

	Error = Error,
}
