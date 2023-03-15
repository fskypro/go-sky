module fsmysql

go 1.19

//replace fsky.pro => github.com/fskypro/gosky/common latest
replace fsky.pro => ../../common

replace fsky.pro/fsmysql => ./

require (
	fsky.pro v0.0.0-00010101000000-000000000000
	fsky.pro/fsmysql v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.6.0
)
