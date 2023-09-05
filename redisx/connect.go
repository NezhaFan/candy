package redisx

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
)

// 根据redis配置初始化一个客户端
func Connect() (*redis.Client, error) {

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Liu@1234", // 没有密码，默认值
		DB:       0,          // 默认DB 0
	})

	err := rdb.Ping(context.Background()).Err()

	return rdb, err
}
