module fspgsql

go 1.21

replace fsky.pro => ../../common
replace fsky.pro/fssearch => ../fssearch

require fsky.pro v0.0.0-00010101000000-000000000000
require fsky.pro/fssearch v0.0.0-00010101000000-000000000000

require github.com/lib/pq v1.10.7
