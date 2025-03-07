一、引用方式
    export LUA_PATH="./fsky/?.lua;;"
    fsky = require("fsky")

二、包属性
    1、fsdefine 模块
        sky.nul
            type: empty table
            说明：表示本包中的空

    2、fsuitl 模块
        fsky.iscallable
            type：function
            说明：判断一个对象是否可调用

    3、fserror 子包
        fsky.Error
            type：fsky.class
            说明：错误基类

    4、fsoo 子包
        fsky.class
            type：function
            说明：创建类的方法，详细用法参看：./fsoo/oo_test.lua

        fsky.Object
            type：fsky.class
            说明：所有 sky.class 创建类的基类

    5、fsstr 子包
        fsky.str.startswith
            type：function
            说明：判断字符串是否以另一个字符串开头

        fsky.str.endswith
            type：function
            说明：判断字符串是否以另一个字符串结尾

        fsky.str.split
            type：function
            以指定子字符串为分隔符拆分字符串(分隔符可以是多个字符的字符串)

        fsky.str.lfill
            type：function
            说明：将字符串设置为指定宽度，不足部分在左边填充指定字符串

        fsky.str.rfill
            type：function
            说明：将字符串设置为指定宽度，不足部分在右边填充自动字符串

    6、fstable 子包
        fsky.Array
            type：fsky.class
            说明：数组，支持插入 nil 值
            参考：fstable/array_test.lua

        fsky.HashMap
            type：fsky.class
            说明：哈希表，支持存储 nil 值
            参考：fstable/hashmap_test.lua

        -----------------------------------------------
        fsky.fstable.listout
            type：function
            说明：返回字符串形式的数组 table

        fsky.fstable.dictout
            type：function
            说明：返回字符串形式的字典 table

        fsky.fstable.update
            type：function
            说明：用参数中指定的多个 table，更新第一个 table，后面的将会覆盖前面的

        fsky.fstale.union
            type：function
            说明：合并多个 table，key 如果在参数中前面的 table 中已经存在，则后面的将会被忽略

    7、fsos 子包
        fsky.systems
            type：table
            说明：所有支持系统列表

        fsky.os.system
            type：function
            说明：返回当前系统编号

        fsky.os.resetSystem
            type：function
            说明：设置系统编号，通常程序启动时第一时间设置，系统类型会影响路径分隔符和文件内容换行符

        fsky.os.pathSpliter
            type：function
            说明：文件夹路径分隔符

        fsky.os.newline
            type：function
            说明：换行符

        fsky.path.join
            type：function
            说明：合并成路径

        -----------------------------------------------
        fsky.path.split
            type：function
            说明：把一个路径按照路径分隔符拆分成一个文件夹数组

        fsky.path.join
            type：function
            说明：把一个跟路径和多个子路径合并成一个路径字符串

        fsky.path.splitext
            type：function
            说明：把一个文件路径，拆分成文件目录、文件名称、文件扩展名(带.)

        fsky.path.normalize
            type：function
            说明：把一个路径修整成最简化，如：/root/aa/bb/.././cc.txt，修正后变为：/root/aa/cc.txt

        fsky.path.filePathExists
            type：function
            说明：判断指定的文件或路径是否存在，存在则返回 true

    8、fslog 子包
        fsky.logfmt
            type：table
            说明：log 文本格式化工具

            fsky.logfmt.fmt            ：将多个参数以空格分开何必成一个字符串输出
            fsky.logfmt.fmtf        ：用参数作为消息的格式化参数
            fsky.logfmt.traceftm    ：与 fmt 相同，但是会附带调用栈信息
            fsky.logfmt.tracefmtf    ：与 fmtf 相同，但会附带调用栈信息

        fsky.BaseLog
            type: fsky.calss
            说明：Log 基类，其拥有以下 log 输出类型(BaseLog 成员方法)：
                debug
                info
                error
                warn
                hack
                trace
                    type：function
                    说明：返回 log 格式字符串，可传入多个参数，以空格分隔输出各个参数
                          类似于：[DEBUG]|2021-05-04 12:33:30.505 dayfilelog_test.lua:4: <msg>
                    参考：fslog/logfmt_test.lua

                debugf
                infof
                errorf
                warnf
                hackf
                tracef
                    type：function
                    说明：返回 log 格式字符串，可传入多个参数以对第一个参数进行格式化
                          类似于：[DEBUG]|2021-05-04 12:33:30.505 dayfilelog_test.lua:4: <msg>
                    参考：fslog/logfmt_test.lua

            BaseLog 还其他成员方法：
                setOutputHandler()  ：设置 log 输出处理函数
                logTypes()          ：获取所有 log 类型
                outputAll()         ：输出所有类型的 log
                shieldType()        ：屏蔽指定类型的 log 输出

        fsky.DFLog
            type：fsky.class，继承于 BaseLog
            说明：以每天创建一个 log 文件的方式，输出 log。

        fsky.gDFLog
            type：fsky.DFLog 对象
            说明：一个全局的 fsky.DFLog。如果整个程序中只有一个 log 系统，则可以直接使用 fsky.gDFLog.init 
                  函数初始化该全局文件 log，然后直接使用该 log 对象进行 log 输出
            参考：fslog/dayfilelog_test.lua
