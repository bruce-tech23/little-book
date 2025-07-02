//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13306)/littlebook",
	},
	Redis: RedisConfig{
		Addr: "localhost:16379",
	},
}
