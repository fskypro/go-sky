module fsredis

go 1.24.1

replace fsky.pro/fsredis => ./

replace fsky.pro => ../../common

require (
	fsky.pro v0.0.0-00010101000000-000000000000
	fsky.pro/fsredis v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.9.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
)
