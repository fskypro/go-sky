module fslog

go 1.21

toolchain go1.21.1

//replace fsky.pro => github.com/fskypro/gosky/common latest
replace fsky.pro => ../common

require fsky.pro v0.0.0-00010101000000-000000000000
