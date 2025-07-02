//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(littlebook-mysql:3318)/littlebook",
	},
	Redis: RedisConfig{
		Addr: "littlebook-redis:6379",
	},
}
