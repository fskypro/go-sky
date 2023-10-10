module fsky.pro

go 1.20

//replace fsky.pro => github.com/fskypro/gosky/common
replace fsky.pro => ./common

require (
	github.com/mbenkmann/goformat v0.0.0-20180512004123-256ef38c4271 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	winterdrache.de/goformat v0.0.0-20180512004123-256ef38c4271 // indirect
)
