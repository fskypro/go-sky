module fsky.pro

go 1.23

//replace fsky.pro => github.com/fskypro/gosky/common
replace fsky.pro => ./common

require (
	github.com/elastic/go-sysinfo v1.13.1 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/mbenkmann/goformat v0.0.0-20180512004123-256ef38c4271 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/shirou/gopsutil v2.21.11+incompatible // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
	winterdrache.de/goformat v0.0.0-20180512004123-256ef38c4271 // indirect
)
