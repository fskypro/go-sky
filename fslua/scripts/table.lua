-- lua table

function get()
    return 4000
end

config  = {
    BaseName = "base",
    Array = {1, 2, 3},

    Base2Name = "base2",
    Slice = {1, 2, "xxx", 2.3},

    Name = "LuaTable",
    IntValue = 1000000,
    FloatValue = 2.2,
    BoolValue = true,

    subConfig = {
        name = "SubTable",
        intValue = 200,
        boolValue = true,

        NestStructs = {
            xx = { a = "100", b = 200 },
            yy = { a = "300", b = 400 },
        }
    },

    inner = {
        map = {[12] = "xx", [13]="yy"},
        myTime = "2024-12-08 12:13:14",
    },

    nestSlice = {
        {
            xx = 1000,
            yy = 2000,
        },
        {
            xx = 3000,
            yy = get(),
        }
    }
}
