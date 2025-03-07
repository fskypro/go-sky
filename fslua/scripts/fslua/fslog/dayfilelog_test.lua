local dflog = require("fsky").gDFLog
dflog.init("test", "./logs")
dflog.setNewLogCmd("./linklog.sh", "daylog")  -- 每次新建 log 文件时，都会执行该命令，并把新 log 文件路径作为第一个命令参数传给命令

dflog.debug("123", "456")
dflog.hack("456", "789")
dflog.trace("aaaa", "bbbbb")

dflog.hackf("aaaaaa: %s", "hehe")
dflog.infof("xxxxx: %d", 200)
dflog.errorf("vvvvv: %d", 200)
dflog.tracef("yyyy: %s", "嘿嘿")

print(dflog.logTypes())
dflog.shieldTypes("debug")   -- 屏蔽 debug（不区分大小写）
dflog.info("debug 类型 log 已经被屏蔽")
dflog.debug("7890")
dflog.error("dfghdfgh")

