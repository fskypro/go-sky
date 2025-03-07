--[[
@copyright: fantasysky 2016
@website: https://www.fsky.pro
@brief: string utils
@author: fanky
@version: 1.0
@date: 2021-04-24
--]]

local fsstr = {}

-- 判断字符串是否以另一个字符串开头
function fsstr.startswith(str, prefix)
	local i = 1
	local chr1, chr2
	while(true) do
		chr1 = string.sub(str, i, i)
		chr2 = string.sub(prefix, i, i)
		if chr1 == '' then
			return chr2 == ''
		end
		if chr2 == '' then
			return true
		end
		if chr1 ~= chr2 then
			return false
		end
		i = i + 1
	end
end

-- 判断字符串是否以另一个字符串结尾
function fsstr.endswith(str, suffix)
	local i = -1
	local chr1, chr2
	while(true) do
		chr1 = string.sub(str, i, i)
		chr2 = string.sub(suffix, i, i)
		if chr1 == '' then
			return chr2 == ''
		end
		if chr2 == '' then
			return true
		end
		if chr1 ~= chr2 then
			return false
		end
		i = i - 1
	end
end

-- 以 sp 为分隔符，拆分字符串 str
function fsstr.split(str, sp)
	local items = {}
	local len = #str
	local i = 1
	local item = ""
	local chr
	while(i <= len) do
		if fsstr.startswith(string.sub(str, i), sp) then
			table.insert(items, item)
			item = ""
			i = i + #sp
		else
			item = item .. string.sub(str, i, i)
			i = i + 1
		end
	end
	if item ~= "" then
		table.insert(items, item)
	end
	return items
end

------------------------------------------------------------
-- 去掉字符串左边的空白字符
function fsstr.ltrim(str)
	str = string.gsub(str, "^[ \t\n\r]+", "")
	return str
end

-- 去掉字符串右边的空白字符
function fsstr.rtrim(str)
	str = string.gsub(str, "[ \t\n\r]+$", "")
	return str
end

-- 去掉字符串两边的空白字符
function fsstr.trim(str)  
	str = string.gsub(str, "^[ \t\n\r]+", "")
	str = string.gsub(str, "[ \t\n\r]+$", "")
	return str
end

------------------------------------------------------------
-- 将返回字符串设置为固定长度 len，并将 str 在返回字符串中右对齐，用 fills 字符串填充左边空缺
-- 如：str.fillleft("123", 10, "->") 将返回：->->->-123
function fsstr.lfill(str, len, fills)
	str = tostring(str)
	local left = len - string.len(str);
	if left <= 0 then return str end
	fills = fills or " "
	local fillcount = math.floor(left / string.len(fills));
	local nfills = string.rep(fills, fillcount);
	local filltail = string.sub(fills, 1, left-#nfills)
	return nfills .. filltail .. str 
end 

-- 将返回字符串设置为固定长度 len，并将 str 在返回字符串中左对齐，用 fills 字符串填充右边空缺
-- 如：str.fillleft("123", 10, "<-") 将返回：123<-<-<-<
function fsstr.rfill(str, len, fills)
	str = tostring(str)
	local right = len - string.len(str)
	if right <= 0 then return str end
	fills = fills or ""
	local fillcount = math.floor(right / string.len(fills))
	local nfills = string.rep(fills, fillcount)
	local filltail = string.sub(fills, 1, right-#nfills)
	return str .. nfills .. filltail
end

----------------------------------------------------------------------
-- initialize
----------------------------------------------------------------------
return {
	init = function(fsky, tofsky)
		fsky.str = fsstr
	end,

	str = fsstr
}
