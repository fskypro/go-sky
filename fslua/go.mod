module fslua

go 1.26

replace fsky.pro => ../common

replace fsky.pro/fslog => ../fslog

require (
	fsky.pro v0.0.0-00010101000000-000000000000
	fsky.pro/fslog v0.0.0-00010101000000-000000000000
	github.com/yuin/gopher-lua v1.1.1
)

require (
	github.com/djherbis/times v1.6.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)
