module fsmysql

go 1.16

//replace fsky.pro => github.com/fskypro/gosky/common latest
replace fsky.pro => ../../common

require (
	fsky.pro v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-sql-driver/mysql v1.6.0
)