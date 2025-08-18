module fsrpc

go 1.24.1

replace fsky.pro => ../common

require (
	fsky.pro v0.0.0-00010101000000-000000000000
	github.com/natefinch/npipe v0.0.0-20160621034901-c1b8fa8bdcce
)

require gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
