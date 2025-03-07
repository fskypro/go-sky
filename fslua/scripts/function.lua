local function get_source_and_line()
    local info = debug.getinfo(3, "Sl")  -- 2 表示调用此函数的代码位置
    return info.source, info.currentline
end


function debugf(fmt, ...)
    local source, line = get_source_and_line()
    local prefix = string.format("[%s:%v]:", source, line)
    fs_debugf(prefix .. fmt, ...)
end

function func(name)
    debugf("hello %s!", name)
    return 100, "call lua", {
        strValue = "xxxx",
        floatValue = 200.2,
    }
end
